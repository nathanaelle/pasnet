package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"os"
	"net"
	"time"
	"strconv"
	"strings"
	"syscall"
)


type Dialer struct {
}



func Dial(network, address string) (net.Conn, error) {
	return	(&Dialer{}).Dial(network, address)
}



func (d *Dialer) Dial(proto, addr string) (c net.Conn, err error) {
	r_addr,err	:= net.ResolveTCPAddr(proto,addr)
	if err != nil {
		return nil, err
	}

	fd, err	:= system_Dialer(proto, r_addr)

	if err != nil {
		return nil, err
	}

	file	:= os.NewFile(uintptr(fd), strings.Join([]string { newfile_prefix, strconv.Itoa(fd),strconv.Itoa(os.Getpid()) }, "_" ) )

	c, err = net.FileConn(file)
	if  err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if err = file.Close(); err != nil {
		syscall.Close(fd)
		c.Close()
		return nil, err
	}

	return
}



func system_Dialer(net string, laddr *net.TCPAddr) (fd int, err error) {
	switch {
		case laddr.IP.To4() != nil:
			var addr	[4]byte
			copy(addr[:], laddr.IP[len(laddr.IP)-4:len(laddr.IP)] )

			return  generic_Dialer(
				func() (int,error){
					return syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
				},
				func(fd int) error{
					return syscall.Connect(fd, &syscall.SockaddrInet4{Port: laddr.Port, Addr: addr })
				})


		case laddr.IP.To16() != nil:
			var addr	[16]byte
			copy(addr[:], laddr.IP )

			return generic_Dialer(
				func() (int,error){
					return syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
				},
				func(fd int) error{
					return syscall.Connect(fd, &syscall.SockaddrInet6{Port: laddr.Port, Addr: addr })
				})

		default:
			return	-1, unknown_proto(net)
	}
}


func generic_Dialer( generic_create func() (int,error), generic_connect func(int) error) (fd int, err error){
	if fd, err = generic_create(); err != nil {
		return -1, err
	}
	err = gatling_run( gatling{
		{ bullet_duration(ka_idle)	, 10*time.Second},
		{ bullet_duration(ka_intvl)	, 5*time.Second	},
		{ bullet_int(ka_count)		, 10		},
		{ bullet_duration(so_linger)	, 3*time.Second	},
		{ bullet_nil(generic_connect)	, nil		},
		{ bullet_int(so_rcvbuf)		, 1<<16		},
		{ bullet_int(so_sndbuf)		, 1<<16		},
		{ bullet_bool(so_nodelay)	, true		},
		{ bullet_bool(so_tcpcork)	, false		},
		{ bullet_bool(so_tcpnopush)	, false		},
		{ bullet_bool(so_nonblock)	, true		},
	}, fd )

	if err != nil {
		syscall.Close(fd)
		return -1, err
	}

	return
}

package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"os"
	"net"
	"sync"
	"strconv"
	"strings"
	"syscall"
//	"crypto/tls"
)


type Dialer struct {
	Proxy		Proxy
	NoIPv6		bool
	NoIPv4		bool
	NoReslv		bool

	end		<-chan struct{}
	wg		*sync.WaitGroup
}


func Dial(network, address string, wg *sync.WaitGroup, end <-chan struct{}) (Conn, error) {
	return	(&Dialer{
		end:	end,
		wg:	wg,
	}).Dial(network, address)
}



func (d *Dialer) Dial(proto, addr string) (Conn, error) {
	if d.Proxy == nil {
		return now_dial(proto, addr, d.wg, d.end)
	}

	if err := d.Proxy.Handshake(); err != nil {
		return nil, err
	}

	return	d.Proxy.Dial(proto,addr)
}


func (d *Dialer) DialTLS(proto, addr string, conf *TLSClientConfig) (Conn, error) {
	c, err	:= d.Dial(proto, addr)
	if err != nil {
		return nil, err
	}

	return c.TLS(conf)
}


func now_dial(proto, addr string, wg *sync.WaitGroup, end <-chan struct{}) (Conn, error) {
	switch proto {
	case	"tcp","tcp4","tcp6":
		sc,err	:= new_socket_common(wg, end)
		if err != nil {
			return nil, err
		}

		r_addr,err	:= net.ResolveTCPAddr(proto,addr)
		if err != nil {
			return nil, err
		}

		fd, err	:= system_dialer(proto, r_addr)
		if err != nil {
			return nil, err
		}

		file	:= os.NewFile(uintptr(fd), strings.Join([]string { newfile_prefix, strconv.Itoa(fd),strconv.Itoa(os.Getpid()) }, "_" ) )

		conn, err := net.FileConn(file)
		if  err != nil {
			syscall.Close(fd)
			return nil, err
		}

		if err = file.Close(); err != nil {
			syscall.Close(fd)
			conn.Close()
			return nil, err
		}

		return &conn_tcp{ conn.(*net.TCPConn), sc }, nil
	}
	return nil, &E_UnknownProto{proto}
}

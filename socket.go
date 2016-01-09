package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"net"
	"time"
	"syscall"
)

const	(
	LINGER_TIMEOUT	time.Duration	= 3 *time.Second
	KA_INTERVAL	time.Duration	= 5 *time.Second
	KA_IDLE		time.Duration	= 10*time.Second
	KA_COUNT	int		= 10

)

func system_listener(net string, laddr *net.TCPAddr) (fd int, err error) {
	switch {
		case laddr.IP.To4() != nil:
			var addr	[4]byte
			copy(addr[:], laddr.IP[len(laddr.IP)-4:len(laddr.IP)] )

			return  listener_tcp_( tcp4_socket,
				func(fd int) error{
					return syscall.Bind(fd, &syscall.SockaddrInet4{Port: laddr.Port, Addr: addr })
				})

		case laddr.IP.To16() != nil:
			var addr	[16]byte
			copy(addr[:], laddr.IP )

			return listener_tcp_( tcp6_socket,
				func(fd int) error{
					return syscall.Bind(fd, &syscall.SockaddrInet6{Port: laddr.Port, Addr: addr })
				})

		default:
			return	-1, &E_UnknownProto{net}
	}
}


func listener_tcp_( generic_create func() (int,error), generic_bind func(int) error) (fd int, err error){
	if fd, err = generic_create(); err != nil {
		return -1, err
	}
	err = gatling(
		bullet(so_reuseaddr	, true		),
		bullet(so_reuseport	, true		),
		bullet(so_fastopen	, 10		),
		bullet(ka_idle		, KA_IDLE	),
		bullet(ka_intvl		, KA_INTERVAL	),
		bullet(ka_count		, KA_COUNT	),
		bullet(so_linger	, LINGER_TIMEOUT),
		bullet(generic_bind	, nil		),
		bullet(so_listen	, -1		),
		bullet(so_nodelay	, true		),
		bullet(so_tcpcork	, false		),
		bullet(so_tcpnopush	, false		),
		bullet(so_nonblock	, true		),
	)( fd )

	if err != nil {
		syscall.Close(fd)
		return -1, err
	}

	return
}





func system_dialer(net string, laddr *net.TCPAddr) (fd int, err error) {
	switch {
		case laddr.IP.To4() != nil:
			var addr	[4]byte
			copy(addr[:], laddr.IP[len(laddr.IP)-4:len(laddr.IP)] )

			return  dialer_tcp_( tcp4_socket,
				func(fd int) error{
					return syscall.Connect(fd, &syscall.SockaddrInet4{Port: laddr.Port, Addr: addr })
				})

		case laddr.IP.To16() != nil:
			var addr	[16]byte
			copy(addr[:], laddr.IP )

			return dialer_tcp_( tcp6_socket,
				func(fd int) error{
					return syscall.Connect(fd, &syscall.SockaddrInet6{Port: laddr.Port, Addr: addr })
				})

	}

	return	-1, &E_UnknownProto{net}
}



func dialer_tcp_( generic_create func() (int,error), generic_connect func(int) error) (fd int, err error){
	if fd, err = generic_create(); err != nil {
		return -1, err
	}
	err = gatling(
		bullet(ka_idle		, KA_IDLE	),
		bullet(ka_intvl		, KA_INTERVAL	),
		bullet(ka_count		, KA_COUNT	),
		bullet(so_linger	, LINGER_TIMEOUT),
		bullet(generic_connect	, nil		),
		bullet(so_rcvbuf	, 1<<16		),
		bullet(so_sndbuf	, 1<<16		),
		bullet(so_nodelay	, true		),
		bullet(so_tcpcork	, false		),
		bullet(so_tcpnopush	, false		),
		bullet(so_nonblock	, true		),
	)( fd )

	if err != nil {
		syscall.Close(fd)
		return -1, err
	}

	return
}

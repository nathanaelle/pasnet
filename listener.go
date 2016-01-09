package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"io"
	"os"
	"net"
	"sync"
	"time"
	"strconv"
	"strings"
	"syscall"
)

type	(

	Listener	interface {
		net.Listener
		Report()			(conns, in_byte, out_byte uint64)

		//here I need uint128
		//Report()	(conns, byte_in, byte_out uint64, sq_in, sq_out uint128)
	}

	tcp_listener	struct {
		*net.TCPListener
		sc	socket_common
	}

	unix_listener	struct {
		*net.UnixListener
		sc	socket_common
	}

	unknown_listener struct {
		net.Listener
		sc	socket_common
	}

)


const	(
	newfile_prefix	string		= "prefix_newfile"
)

var	(
	IO_TIMEOUT	time.Duration	= 200*time.Millisecond
)


func NewListener(proto string, laddr *net.TCPAddr, wg *sync.WaitGroup, end <-chan struct{}) (Listener, error)  {
	fd,err	:= system_listener(proto, laddr)
	if err	!= nil {
		return nil, err
	}

	return ListenerFromFD(fd, wg, end)
}


func ListenerFromFD(fd int, wg *sync.WaitGroup, end <-chan struct{}) (ln Listener, err error) {
	file	:= os.NewFile( uintptr(fd), strings.Join( []string { newfile_prefix, strconv.Itoa(fd), strconv.Itoa(os.Getpid()) }, "_" ) )
	l,err	:= net.FileListener(file)
	if err	!= nil {
		syscall.Close(fd)
		return nil, err
	}
	ln,err	= ListenerFromNet(l, wg, end)
	if err	!= nil {
		syscall.Close(fd)
		return nil, err
	}

	if err	= file.Close(); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	return
}


func ListenerFromNet(iln net.Listener, wg *sync.WaitGroup, end <-chan struct{}) (Listener,error) {
	sc,err	:= new_socket_common(wg, end)
	if err != nil {
		return nil, err
	}

	switch	iln.(type) {
	case	*net.UnixListener:
		return	&unix_listener	{ iln.(*net.UnixListener), sc }, nil

	case	*net.TCPListener:
		return	&tcp_listener	{ iln.(*net.TCPListener), sc }, nil

	default:
		return	&unknown_listener { iln, sc }, nil
	}
}


func (lst *unix_listener)Accept() (net.Conn,error) {
	return	lst.AcceptUnix()
}

func (lst *unix_listener)AcceptUnix() (*conn_unix,error) {
	for {
		select {
		case	<-lst.sc.end:
			return nil,io.EOF

		default:
			lst.UnixListener.SetDeadline(time.Now().Add(IO_TIMEOUT))
			fd,err := lst.UnixListener.AcceptUnix()
			switch	{
			case	err == nil:
				lst.sc.AddConn()
				return &conn_unix { fd, lst.sc }, nil

			default:
				if not_temporary_timeout(err) {
					return nil,err
				}
			}
		}
	}
}

func (lst *unix_listener)Close() (err error) {
	err = lst.UnixListener.Close()
	lst.sc.Done()
	return
}

func (lst *unix_listener) Addr() (net.Addr) {
	return	lst.UnixListener.Addr()
}

func (lst *unix_listener) Report()	(uint64,uint64,uint64) {
	return	lst.sc.Report()
}


func (lst *tcp_listener)Accept() (net.Conn,error) {
	return	lst.AcceptTCP()
}


func (lst *tcp_listener)AcceptTCP() (*conn_tcp,error) {
	for {
		select {
		case	<-lst.sc.end:
			return nil,io.EOF

		default:
			lst.TCPListener.SetDeadline(time.Now().Add(IO_TIMEOUT))
			fd,err := lst.TCPListener.AcceptTCP()
			switch	{
			case	err == nil:
				lst.sc.AddConn()
				return &conn_tcp{ fd, lst.sc }, nil

			default:
				if not_temporary_timeout(err) {
					return nil,err
				}
			}
		}
	}
}

func (lst *tcp_listener)Close() (err error) {
	err = lst.TCPListener.Close()
	lst.sc.Done()
	return
}

func (lst *tcp_listener)Addr() (net.Addr) {
	return	lst.TCPListener.Addr()
}

func (lst *tcp_listener)Report()	(uint64,uint64,uint64) {
	return	lst.sc.Report()
}



func (lst *unknown_listener)Accept() (net.Conn,error) {
	for {
		select {
		case	<-lst.sc.end:
			return nil,io.EOF

		default:
			fd,err := lst.Listener.Accept()
			switch	{
			case	err == nil:
				lst.sc.AddConn()
				return &conn_unknown{ fd, lst.sc }, nil

			default:
				if not_temporary_timeout(err) {
					return nil,err
				}
			}
		}
	}
}

func (lst *unknown_listener)Close() (err error) {
	err = lst.Listener.Close()
	lst.sc.Done()
	return
}

func (lst *unknown_listener)Addr() (net.Addr) {
	return	lst.Listener.Addr()
}

func (lst *unknown_listener)Report()	(uint64,uint64,uint64) {
	return	lst.sc.Report()
}

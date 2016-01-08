// +build linux freebsd darwin

package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"syscall"
	"time"
	"os"
)

func tcp4_socket() (int,error){
	return syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
}

func tcp6_socket() (int,error){
	return syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
}

func so_rcvbuf(fd int, n int) error {
	return os.NewSyscallError("so_rcvbuf", syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, n))
}

func so_sndbuf(fd int, n int) error {
	return os.NewSyscallError("so_sndbuf", syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, n))
}

func so_nodelay(fd int, flag bool) error {
	return os.NewSyscallError("so_nodelay", syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, boolint(flag)))
}

func so_reuseaddr(fd int, flag bool) error {
	return os.NewSyscallError("so_reuseaddr", syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, boolint(flag)) )
}

func so_nonblock(fd int, flag bool) error {
	return os.NewSyscallError("so_nonblock", syscall.SetNonblock(fd, true))
}

func so_linger(fd int, d time.Duration) error {
	if d == 0 {
		return os.NewSyscallError("so_linger", syscall.SetsockoptLinger(fd, syscall.SOL_SOCKET, syscall.SO_LINGER, &syscall.Linger { 0, 0 } ))
	}

	// cargo cult from src/net/tcpsockopt_unix.go
	d	+= (time.Second - time.Nanosecond)
	l	:= syscall.Linger { 1, int32(d.Seconds()) }

	return os.NewSyscallError("so_linger", syscall.SetsockoptLinger(fd, syscall.SOL_SOCKET, syscall.SO_LINGER, &l ))
}

func so_listen(fd int,queue int) error {
	if queue <1 {
		return os.NewSyscallError("so_listen", syscall.Listen(fd, syscall.SOMAXCONN) )
	}

	return os.NewSyscallError("so_listen", syscall.Listen(fd, queue) )

}

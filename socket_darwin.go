// +build darwin

package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"syscall"
	"time"
	"os"
)


const TCP_KEEPINTVL = 0x101

func ka_idle(fd int, d time.Duration) error {
	if d == 0 {
		return os.NewSyscallError("ka_idle", syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPALIVE, 0 ))
	}

	// cargo cult from src/net/tcpsockopt_unix.go
	d += (time.Second - time.Nanosecond)
	return os.NewSyscallError("ka_idle", syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPALIVE, int(d.Seconds()) ))
}

func ka_intvl(fd int, d time.Duration) error {
	if d == 0 {
		return os.NewSyscallError("ka_intv", syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, TCP_KEEPINTVL, 0 ))
	}

	// cargo cult from src/net/tcpsockopt_unix.go
	d += (time.Second - time.Nanosecond)
	return os.NewSyscallError("ka_intv", syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, TCP_KEEPINTVL, int(d.Seconds()) ))
}

func ka_count(fd int, n int) error {
	return nil
}

func so_tcpcork(fd int, flag bool) error {
	return nil
}

func so_tcpnopush(fd int, flag bool) error {
	return os.NewSyscallError("so_tcpnopush", syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_NOPUSH, boolint(flag)))
}

func so_reuseport(fd int, flag bool) error {
	return os.NewSyscallError("so_reuseport", syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, boolint(flag)) )
}

func so_fastopen(fd int, n int) error {
	return nil
}

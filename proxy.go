package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
)

type	Proxy	interface {
	Handshake() error
	Dial(net, addr string)		(Conn, error)
	Listen(net, addr string)	(Listener, error)
}

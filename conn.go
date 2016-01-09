package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"io"
	"net"
	"time"

	"crypto/tls"
)

type	(

	Conn	interface {
		net.Conn
		IsTLS()			bool
		TLSState()		*TLSState
		TLS(*TLSClientConfig)	(Conn, error)
	}


	conn_tcp	struct {
		*net.TCPConn
		sc	socket_common
	}

	conn_unix	struct {
		*net.UnixConn
		sc	socket_common
	}

	conn_unknown	struct {
		net.Conn
		sc	socket_common
	}

	conn_tls	struct {
		*tls.Conn
		sc	socket_common
		conf	*TLSClientConfig
		state	*TLSState
	}

)


func common_read(conn net.Conn, sc socket_common, b []byte) (n int, err error)  {
	defer func(){
		if n > 0 {
			sc.AddInBytes(uint64(n))
		}
	}()

	n1	:= 0

	for {
		select {
		case	<-sc.end:
			return 0,io.EOF

		default:
			conn.SetReadDeadline(time.Now().Add(IO_TIMEOUT))
			n1,err = conn.Read(b[n:])
			n+=n1
			if err == nil || n == len(b) {
				conn.SetReadDeadline(time.Time{})
				return n,nil
			}

			if not_temporary_timeout(err) {
				return
			}
		}
	}
}


func common_write(conn net.Conn, sc socket_common, b []byte) (n int, err error)  {
	defer func(){
		if n > 0 {
			sc.AddOutBytes(uint64(n))
		}
	}()
	n1	:= 0

	for {
		select {
		case	<-sc.end:
			return n,io.EOF

		default:
			conn.SetWriteDeadline(time.Now().Add(IO_TIMEOUT))
			n1,err = conn.Write(b[n:])
			n+=n1
			if err == nil || n == len(b) {
				conn.SetWriteDeadline(time.Time{})
				return n,nil
			}

			if not_temporary_timeout(err) {
				return
			}
		}
	}
}

func common_tls(c net.Conn, config *TLSClientConfig, sc socket_common) (Conn, error) {
	if config == nil {
		return nil, &E_MissingArgument{ "conf", "*pasnet.TLSClientConfig" }
	}

	tlsconn	:= tls.Client(c, config.GetTLSConfig())
	err	:= tlsconn.Handshake()
	if err != nil {
		return nil, err
	}

	st,err	:= config.Verify(tlsconn)
	if err != nil {
		return nil, err
	}

	return	&conn_tls{ tlsconn, sc, config, st }, nil
}


func (c *conn_tcp) Read(b []byte) (int, error) {
	return	common_read(c.TCPConn, c.sc, b)
}

func (c *conn_tcp) Write(b []byte) (int, error) {
	return	common_write(c.TCPConn, c.sc, b)
}

func (c *conn_tcp) TLS(config *TLSClientConfig) (Conn, error) {
	return	common_tls(c.TCPConn, config, c.sc)
}

func (c *conn_tcp) IsTLS() bool {
	return	false
}

func (c *conn_tcp) TLSState() *TLSState {
	return	nil
}



func (c *conn_unix) Read(b []byte) (int, error) {
	return	common_read(c.UnixConn, c.sc, b)
}

func (c *conn_unix) Write(b []byte) (int, error) {
	return	common_write(c.UnixConn, c.sc, b)
}

func (c *conn_unix) TLS(config *TLSClientConfig) (Conn, error) {
	return	common_tls(c.UnixConn, config, c.sc)
}

func (c *conn_unix) IsTLS() bool {
	return	false
}

func (c *conn_unix) TLSState() *TLSState {
	return	nil
}





func (c *conn_unknown) Read(b []byte) (int, error) {
	return	common_read(c.Conn, c.sc, b)
}

func (c *conn_unknown) Write(b []byte) (int, error) {
	return	common_write(c.Conn, c.sc, b)
}

func (c *conn_unknown) TLS(config *TLSClientConfig) (Conn, error) {
	return	common_tls(c.Conn, config, c.sc)
}

func (c *conn_unknown) IsTLS() bool {
	return	false
}

func (c *conn_unknown) TLSState() *TLSState {
	return	nil
}



func (c *conn_tls) Read(b []byte) (int, error) {
	return	c.Conn.Read(b)
}

func (c *conn_tls) Write(b []byte) (int, error) {
	return	c.Conn.Write(b)
}

func (c *conn_tls) IsTLS() bool {
	return	true
}

func (c *conn_tls) TLSState() *TLSState {
	return	c.state
}

func (c *conn_tls) TLS(config *TLSClientConfig) (Conn, error) {
	return c, nil
}

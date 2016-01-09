package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"sync"
	"sync/atomic"
)


type	(

	socket_common struct {
		end		<-chan struct{}
		wg		*sync.WaitGroup
		conns,in, out	*uint64
	}

)


func	new_socket_common(wg *sync.WaitGroup, end <-chan struct{}) (sc socket_common,err error) {
	if wg == nil {
		return sc, &E_MissingArgument{"wg", "*sync.WaitGroup"}
	}
	if end == nil {
		return sc, &E_MissingArgument{"end","<-chan struct{}"}
	}

	return socket_common{ end, wg, new(uint64), new(uint64), new(uint64) }, nil
}

func (sc socket_common)	AddConn() {
	sc.wg.Add(1)
	atomic.AddUint64(sc.conns, 1)
}

func (sc socket_common)	AddInBytes(size uint64) {
	atomic.AddUint64(sc.in, size)
}

func (sc socket_common)	AddOutBytes(size uint64) {
	atomic.AddUint64(sc.out, size)
}

func (sc socket_common)	Done() {
	sc.wg.Done()
}

func (sc socket_common)Report()	(uint64,uint64,uint64) {
	return	atomic.LoadUint64(sc.conns), atomic.LoadUint64(sc.in), atomic.LoadUint64(sc.out)
}

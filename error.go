package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"fmt"
)




type	(

	PasNetError		interface {
		error
	}

	E_UnknownProto		struct {
		Proto	string
	}

	E_UnknownSocksVersion	struct {
		Version	int
	}

	E_MissingArgument	struct {
		Arg, Type	string
	}


)


func (e *E_UnknownProto)	Error()	string	{
	return	fmt.Sprintf("Unknown proto %s", e.Proto )
}

func (e *E_UnknownSocksVersion)	Error()	string	{
	return	fmt.Sprintf("Unknown SOCKS version %d", e.Version )
}

func (e *E_MissingArgument)	Error()	string	{
	return	fmt.Sprintf("Missing Non nil Argument %s %s", e.Arg, e.Type )
}

package pasnet	// import "github.com/nathanaelle/pasnet"


import	(
	"net"
)


func boolint(b bool) int {
	switch b {
		case true:	return 1
		default:	return 0
	}
}



func not_temporary_timeout(err error) bool {
	nerr,ok := err.(net.Error)
	return !ok || !(nerr.Timeout() && nerr.Temporary())
}

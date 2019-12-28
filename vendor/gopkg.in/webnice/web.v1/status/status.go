package status

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
)

// Text returns a text for the HTTP status code as string
// If the code is unknown, return "HTTP status code %d"
func Text(code int) (ret string) {
	var ok bool
	if ret, ok = statusText[code]; !ok {
		ret = fmt.Sprintf("HTTP status code %d", code)
	}
	return
}

// Bytes returns a text for the HTTP status code as []byte
// If the code is unknown, return "HTTP status code %d"
func Bytes(code int) (ret []byte) {
	var ok bool
	var buf string
	if buf, ok = statusText[code]; !ok {
		buf = fmt.Sprintf("HTTP status code %d", code)
	}
	ret = []byte(buf)
	return
}

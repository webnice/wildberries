package header

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
)

// Interface is an interface of package
type Interface interface {
	// Add adds the key, value pair to the header.
	Add(key string, value string)

	// Del deletes the values associated with key.
	Del(key string)

	// Get gets the first value associated with the given key.
	Get(key string) string

	// IsSet Check, is key set
	IsSet(key string) bool

	// Set sets the header entries associated with key to the single element value.
	Set(key string, value string)

	// Names Getting a list of all key names
	Names() (ret []string)

	// Len Return number of defined header names
	Len() int

	// Reset all headers
	Reset()
}

// impl is an implementation of package
type impl struct {
	Header http.Header
}

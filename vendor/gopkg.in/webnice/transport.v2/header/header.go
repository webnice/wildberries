package header

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
)

// New creates a new object and return interface
func New(item ...http.Header) Interface {
	var (
		i   int
		hdr = new(impl)
	)

	for i = range item {
		hdr.Header = item[i]
		return hdr
	}
	hdr.Header = make(http.Header)

	return hdr
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
func (hdr *impl) Add(key string, value string) {
	hdr.Header.Add(key, value)
}

// Del deletes the values associated with key.
func (hdr *impl) Del(key string) {
	hdr.Header.Del(key)
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
// To access multiple values of a key, access the map directly
// with CanonicalHeaderKey.
func (hdr *impl) Get(key string) string {
	return hdr.Header.Get(key)
}

// IsSet Check, is key set
func (hdr *impl) IsSet(key string) (ok bool) {
	_, ok = hdr.Header[key]
	return
}

// Set sets the header entries associated with key to
// the single element value.  It replaces any existing
// values associated with key.
func (hdr *impl) Set(key string, value string) {
	hdr.Header.Set(key, value)
}

// Names Getting a list of all key names
func (hdr *impl) Names() (ret []string) {
	for kn := range hdr.Header {
		ret = append(ret, kn)
	}

	return
}

// Len Return number of defined header names
func (hdr *impl) Len() int {
	return len(hdr.Header)
}

// Reset all headers
func (hdr *impl) Reset() {
	for kn := range hdr.Header {
		delete(hdr.Header, kn)
	}
}

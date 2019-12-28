package data

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
)

// NewReadCloser Создание нового объекта на основе io.Reader
func NewReadCloser(w io.Reader, fn func() error) ReadCloser {
	return &readCloserImplementation{essence: w, closer: fn}
}

// Read Реализация Writer
func (rd *readCloserImplementation) Read(p []byte) (int, error) {
	return rd.essence.Read(p)
}

// Close Реализация Close
func (rd *readCloserImplementation) Close() error {
	if rd.closer != nil {
		return rd.closer()
	}
	return nil
}

package data

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"io"
	"os"
)

// WriteCloser is an interface
type WriteCloser interface {
	Write(p []byte) (int, error)
	Close() error
}

// ReadCloser is an interface
type ReadCloser interface {
	Read(p []byte) (int, error)
	Close() error
}

// ReadAtSeekerWriteToCloser is an interface
type ReadAtSeekerWriteToCloser interface {
	io.Seeker
	io.Closer
	io.Reader
	io.ReaderAt
	io.WriterTo

	// Size is an size of content
	Size() int64

	// Done Completion of work with data
	Done()
}

// writeCloserImplementation is an implementation
type writeCloserImplementation struct {
	essence io.Writer
	closer  func() error
}

// readCloserImplementation is an implementation
type readCloserImplementation struct {
	essence io.Reader
	closer  func() error
}

// readAtSeekerWriteToCloserImplementation is an implementation
type readAtSeekerWriteToCloserImplementation struct {
	contentData *bytes.Reader // Контент в памяти
	contentFh   *os.File      // Контент в файле
}

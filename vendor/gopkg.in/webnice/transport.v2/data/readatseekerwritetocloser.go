package data

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"io"
	"os"
)

// NewReadAtSeekerWriteToCloser Реализация интерфейсов io для *bytes.Byffer и *os.File
func NewReadAtSeekerWriteToCloser(r interface{}) ReadAtSeekerWriteToCloser {
	var rd *readAtSeekerWriteToCloserImplementation
	switch t := r.(type) {
	case *bytes.Buffer:
		rd = &readAtSeekerWriteToCloserImplementation{contentData: bytes.NewReader(t.Bytes())}
	case *os.File:
		rd = &readAtSeekerWriteToCloserImplementation{contentFh: t}
	default:
		panic("NewReadAtSeekerCloser allowed only *bytes.Buffer or *os.File types")
	}
	return rd
}

// Seek is the basic Seek method
func (rd *readAtSeekerWriteToCloserImplementation) Seek(offset int64, whence int) (ret int64, err error) {
	if rd.contentFh != nil {
		ret, err = rd.contentFh.Seek(offset, whence)
	}
	if rd.contentData != nil {
		ret, err = rd.contentData.Seek(offset, whence)
	}
	return
}

// ReadAt is the basic ReadAt method
func (rd *readAtSeekerWriteToCloserImplementation) ReadAt(b []byte, off int64) (n int, err error) {
	if rd.contentFh != nil {
		n, err = rd.contentFh.ReadAt(b, off)
	}
	if rd.contentData != nil {
		n, err = rd.contentData.ReadAt(b, off)
	}
	return
}

// Read is the basic Read method
func (rd *readAtSeekerWriteToCloserImplementation) Read(b []byte) (n int, err error) {
	if rd.contentFh != nil {
		n, err = rd.contentFh.Read(b)
	}
	if rd.contentData != nil {
		n, err = rd.contentData.Read(b)
	}
	return
}

// WriteTo is the basic WriteTo method
func (rd *readAtSeekerWriteToCloserImplementation) WriteTo(w io.Writer) (n int64, err error) {
	if rd.contentFh != nil {
		n, err = io.Copy(w, rd.contentFh)
	}
	if rd.contentData != nil {
		n, err = rd.contentData.WriteTo(w)
	}
	return
}

// Close is the basic Close method
func (rd *readAtSeekerWriteToCloserImplementation) Close() (err error) {
	if rd.contentFh != nil {
		_, err = rd.contentFh.Seek(0, io.SeekEnd)
	}
	if rd.contentData != nil {
		_, err = rd.contentData.Seek(0, io.SeekEnd)
	}
	return
}

// Size is an size of content
func (rd *readAtSeekerWriteToCloserImplementation) Size() (n int64) {
	if rd.contentFh != nil {
		if fi, err := rd.contentFh.Stat(); err == nil {
			n = fi.Size()
		}
	}
	if rd.contentData != nil {
		n = rd.contentData.Size()
	}
	return
}

// Done Completion of work with data
func (rd *readAtSeekerWriteToCloserImplementation) Done() {
	if rd.contentFh != nil {
		_ = rd.contentFh.Close()
	}
}

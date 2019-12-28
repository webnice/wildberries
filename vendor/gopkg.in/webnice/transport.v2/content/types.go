package content

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"

	"gopkg.in/webnice/transport.v2/data"

	"golang.org/x/text/encoding"
)

// Interface is an interface of package
type Interface interface {
	io.WriterTo

	// Transcode is an transcoding content from the specified encoding to UTF-8
	Transcode(e encoding.Encoding) Interface

	// Transform is an transforming content using a custom function
	Transform(fn TransformFunc) Interface

	// String Return content as string
	String() (string, error)

	// Bytes Return content as []byte
	Bytes() ([]byte, error)

	// UnmarshalJSON Decoding content like JSON
	UnmarshalJSON(o interface{}) error

	// UnmarshalXML Decoding content like XML
	UnmarshalXML(o interface{}) error

	// UnTar Разархивация контента методом TAR
	UnTar() Interface

	// UnZip Разархивация контента методом ZIP (извлекается только первый файл)
	UnZip() Interface

	// UnGzip Разархивация контента методом GZIP
	UnGzip() Interface

	// UnFlate Разархивация контента методом FLATE
	UnFlate() Interface

	// BackToBegin Returns the content reading pointer to the beginning
	// This allows you to repeat the work with content
	BackToBegin() error
}

// TransformFunc is an func for streaming content conversion
type TransformFunc func(r io.Reader) (io.Reader, error)

// impl is an implementation of package
type impl struct {
	esence data.ReadAtSeekerWriteToCloser // Данные контента
	rdc    io.ReadCloser                  // Интерфейс

	transcode encoding.Encoding // Если не равно nil, то контент перекодируется на лету из указанной кодировки
	transform TransformFunc     // Функция потокового преобразования контента
	unzip     bool              // =true - контент разархивируется алгоритмом сжатия ZIP, возвращается первый файл в архиве
	untar     bool              // =true - контент разархивируется алгоритмом сжатия TAR
	ungzip    bool              // =true - контент разархивируется алгоритмом сжатия GZIP
	unflate   bool              // =true - контент разархивируется алгоритмом сжатия FLATE
}

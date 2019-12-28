package response

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"gopkg.in/webnice/transport.v2/charmap"
	"gopkg.in/webnice/transport.v2/content"
	"gopkg.in/webnice/transport.v2/data"
	"gopkg.in/webnice/transport.v2/header"
)

const (
	// Максимальный размер данных загружаемый в память 250Mb
	maxDataSizeLoadedInMemory = uint64(250 * 1024 * 1024)
)

// Pool is an interface of package
type Pool interface {
	// ResponseGet Извлечение из pool нового элемента Response
	ResponseGet() Interface

	// ResponsePut Возврат в sync.Pool использованного элемента Response
	ResponsePut(req Interface)
}

// Interface is an interface of package
type Interface interface {
	// DebugFunc Set debug func and enable or disable debug mode
	// If fn=not nil - debug mode is enabled. If fn=nil, debug mode is disbled
	DebugFunc(fn DebugFunc) Interface

	// Do Выполнение запроса и получение Response
	Do(client *http.Client, request *http.Request) error

	// Load all response data
	Load() error

	// Error Return latest error
	Error() error

	// Response Returns the http.Response as is
	Response() *http.Response

	// ContentLength records the length of the associated content
	ContentLength() int64

	// Cookies parses and returns the cookies set in the Set-Cookie headers
	Cookies() []*http.Cookie

	// Latency is an request latency for reading body of response without reading header of response
	Latency() time.Duration

	// StatusCode is an http status code of response
	StatusCode() int

	// Status is an http status string of response, for known HTTP codes
	Status() string

	// Header maps header keys to values. If the response had multiple headers with the same key,
	// they may be concatenated, with comma delimiters
	Header() header.Interface

	// Charmap interface
	Charmap() charmap.Charmap

	// Content() Interface for working with response content
	Content() content.Interface
}

// impl is an implementation of package
type impl struct {
	responsePool *sync.Pool // Пул объектов Response
}

// DebugFunc Is an a function for debug request/response data
type DebugFunc func(data []byte)

// Response is an Response implementation
type Response struct {
	err                   error                          // Latest error
	response              *http.Response                 // http.Response object
	debugFunc             DebugFunc                      // Is an a function for debug request/response data. If not nil - debug mode is enabled. If nil, debug mode is disbled
	timeBegin             time.Time                      // Дата и время начала загрузки результата запроса
	timeLatency           time.Duration                  // Время ушедшее за выполнение загрузки результата запроса
	contentInMemory       bool                           // =true - Результат в памяти, =false - результат во временном файле
	contentData           *bytes.Buffer                  // Результат запроса в памяти
	contentFilename       string                         // Имя времененного файла результата запроса
	contentFh             *os.File                       // Интерфейс файлового дескриптора временного файла
	contentTemporaryFiles []string                       // Имена временных файлов
	contentWriteCloser    io.WriteCloser                 // Интерфейс io.WriteCloser к результату запроса в памяти
	contentLength         int64                          // Размер загруженных данных
	contentReader         data.ReadAtSeekerWriteToCloser // Интерфейс к данным результата запроса
	charmap               charmap.Charmap                // charmap interface

	// Переменные
	tmpOk     bool      // Общая переменная
	tmpTm     time.Time // Общая переменная
	tmpString string    // Общая переменная
	tmpI      int       // Общая переменная
}

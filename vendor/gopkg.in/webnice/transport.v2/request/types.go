package request

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"gopkg.in/webnice/transport.v2/header"
	"gopkg.in/webnice/transport.v2/methods"
	"gopkg.in/webnice/transport.v2/response"
)

// Pool is an interface of package
type Pool interface {
	// RequestGet Извлечение из pool нового элемента Request
	RequestGet() Interface

	// RequestPut Возврат в sync.Pool использованного элемента Request
	RequestPut(req Interface)
}

// Interface is an request interface
type Interface interface {
	// Cancel Aborting request
	Cancel() Interface

	// Done Waiting for the request to finish
	Done() Interface

	// Error Return latest error
	Error() error

	// DebugFunc Set debug func and enable or disable debug mode
	// If fn=not nil - debug mode is enabled. If fn=nil, debug mode is disbled
	DebugFunc(fn DebugFunc) Interface

	// Request Returns the http.Request prepared for the request
	Request() (*http.Request, error)

	// Method Set request method
	Method(methods.Value) Interface

	// URL Set request URL
	URL(url string) Interface

	// Referer Setting the referer header
	Referer(referer string) Interface

	// UserAgent Setting the UserAgent Request Header
	UserAgent(userAgent string) Interface

	// ContentType Setting the Content-Type request header
	ContentType(contentType string) Interface

	// Accept Setting the Accept request header
	Accept(accept string) Interface

	// AcceptEncoding Setting the Accept-Encoding request header
	AcceptEncoding(acceptEncoding string) Interface

	// AcceptLanguage Setting the Accept-Language request header
	AcceptLanguage(acceptLanguage string) Interface

	// Settings the Accept-Charset request header
	AcceptCharset(acceptCharset string) Interface

	// Settings custom request header
	CustomHeader(name string, value string) Interface

	// BasicAuth Set login and password for request basic authorization
	BasicAuth(username string, password string) Interface

	// Cookies Adding cookies to the request
	Cookies(cookies []*http.Cookie) Interface

	// Header is an interface for custom headers manipulation
	Header() header.Interface

	// Latency is an request latency without reading body of response
	Latency() time.Duration

	// Response is an response interface
	Response() response.Interface

	// DATA OF THE REQUEST

	// DataStream Data for the request body in the form of an interface of io.Reader
	DataStream(data io.Reader) Interface

	// DataString Data for the request body in the form of an string
	DataString(data string) Interface

	// DataString Data for the request body in the form of an []byte
	DataBytes(data []byte) Interface

	// DataJSON The data for the query is created by serializing the object in JSON
	DataJSON(data interface{}) Interface

	// DataXML The data for the query is created by serializing the object in XML
	DataXML(data interface{}) Interface

	// EXECUTING

	// Do Executing the query and getting the result
	Do(client *http.Client) error
}

// impl is an implementation of package
type impl struct {
	methods      methods.Interface // Интерфейс методов запроса
	requestPool  *sync.Pool        // Пул объектов Request
	responsePool response.Pool     // Интерфейс пула объектов Response
}

// DebugFunc Is an a function for debug request/response data
type DebugFunc func(data []byte)

// Request is an Request implementation
type Request struct {
	context              context.Context    // Context interface
	contextCancelFunc    context.CancelFunc // Context CancelFunc
	method               methods.Value      // Метод запроса данных
	header               header.Interface   // Заголовки запроса
	err                  error              // Latest error
	debugFunc            DebugFunc          // Is an a function for debug request/response data. If not nil - debug mode is enabled. If nil, debug mode is disbled
	url                  *bytes.Buffer      // Запрашиваемый URL без данных
	request              *http.Request      // Объект net/http.Request
	requestData          *bytes.Reader      // Данные запроса
	requestDataInterface io.Reader          // Интерфейс данных запроса
	username             string             // Имя пользователя авторизации, если указан, то передаются заголовки авторизации
	password             string             // Пароль авторизации
	cookie               []*http.Cookie     // Печеньги запроса
	timeBegin            time.Time          // Дата и время начала запроса
	timeLatency          time.Duration      // Время ушедшее за выполнение запроса
	response             response.Interface // Интерфейс результата запроса

	tmpArr     []string // Variable
	tmpOk      bool     // Variable
	tmpCounter int      // Variable
	tmpBytes   []byte   // Variable
}

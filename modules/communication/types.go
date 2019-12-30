package communication

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"context"
	"time"

	"gopkg.in/webnice/transport.v2"
	"gopkg.in/webnice/transport.v2/methods"
	"gopkg.in/webnice/transport.v2/request"
	"gopkg.in/webnice/transport.v2/response"
)

const (
	defaultMaximumIdleConnections        = uint(1000)       // Максимальное общее число бездействующих keepalive соединений
	defaultMaximumIdleConnectionsPerHost = uint(10)         // Максимальное число бездействующих keepalive соединений для каждого хоста
	defaultDialContextTimeout            = time.Second * 3  // Таймаут установки соединения с хостом
	defaultIdleConnectionTimeout         = time.Minute * 5  // Таймаут keepalive соединения до обрыва связи
	defaultTotalTimeout                  = time.Second * 90 // Общий таймаут на весь процесс связи, включает соединение, отправку данных, получение ответа
	defaultRequestPoolSize               = uint16(20)       // Размер пула воркеров готовых для выполнения запросов к хостам
)

const (
	// EncodingGzip gzip compression
	EncodingGzip = `gzip`

	// EncodingDeflate deflate compression
	EncodingDeflate = `deflate`

	// AcceptEncoding Поддерживаемые транспортом протоколы сжатия данных
	AcceptEncoding = `gzip, defalte`

	// AcceptJSON Стандартный заголовок ожидаемого контента ответа
	AcceptJSON = `application/json`

	// AcceptLanguage Стандартный заголовок браузера Accept-Language
	AcceptLanguage = `ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7`

	// CacheControl Стандартный заголовок Cache-Control
	CacheControl = `no-cache`

	// UserAgent Стандартный заголовок User-Agent
	UserAgent = `WEBNICE wildberries/modules/communication/` + version
)

// Interface is an interface of package
type Interface interface {
	// Transport Готовый к использованию интерфейс коммуникации с сервером
	Transport() transport.Interface

	// NewRequest Базовый метод создания объекта запроса
	NewRequest(uri string, mtd methods.Value) (ret request.Interface)

	// RequestJSON Подготовка запроса для получения JSON ответа
	RequestJSON(req request.Interface) (ret request.Interface)

	// RequestResponse Выполнение запроса, ожидание и получение результата
	RequestResponse(ctx context.Context, req request.Interface) (ret response.Interface, err error)

	// RequestResponseStatusCode Выполнение запроса, ожидание и получение результата в виде HTTP статуса
	RequestResponseStatusCode(ctx context.Context, req request.Interface) (statusCode int, err error)

	// RequestResponsePlainText Выполнение запроса, ожидание и получение результата в виде текста
	RequestResponsePlainText(ctx context.Context, req request.Interface) (ret *bytes.Buffer, statusCode int, err error)

	// RequestResponseJSON Выполнение запроса, ожидание и получение результата в виде JSON
	RequestResponseJSON(ctx context.Context, req request.Interface, data interface{}) (statusCode int, err error)

	// ERRORS

	// Errors Ошибки известного состояни, которые могут вернуть функции пакета
	Errors() *Error
}

// impl is an implementation of package
type impl struct {
	singleton transport.Interface
}

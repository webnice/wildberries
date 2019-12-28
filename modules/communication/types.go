package communication // import "git.webdesk.ru/wd/kit/modules/communication"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	"gopkg.in/webnice/transport.v2"
)

const (
	defaultMaximumIdleConnections        = uint(1000)       // Максимальное общее число бездействующих keepalive соединений
	defaultMaximumIdleConnectionsPerHost = uint(100)        // Максимальное число бездействующих keepalive соединений для каждого хоста
	defaultDialContextTimeout            = time.Second * 3  // Таймаут установки соединения с хостом
	defaultIdleConnectionTimeout         = time.Minute * 5  // Таймаут keepalive соединения до обрыва связи
	defaultTotalTimeout                  = time.Second * 90 // Общий таймаут на весь процесс связи, включает соединение, отправку данных, получение ответа
	defaultRequestPoolSize               = uint16(20)       // Размер пула воркеров готовых для выполнения запросов к хостам
)

const (
	// BearerHeader Префикс заголовка авторизации
	BearerHeader = `Bearer `

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
	UserAgent = `WEBDESK modules/communication/` + version
)

var (
	singleton transport.Interface
)

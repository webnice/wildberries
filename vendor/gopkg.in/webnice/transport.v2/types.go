package transport

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/webnice/transport.v2/methods"
	"gopkg.in/webnice/transport.v2/request"
)

const (
	defaultRequestPoolSize               = uint16(1)        // Specifies a number of workers in the query pool
	defaultDialContextTimeout            = time.Duration(0) // Is the maximum amount of time a dial will wait for a connect to complete. The default is no timeout
	defaultDialContextKeepAlive          = 30 * time.Second // Specifies the keep-alive period for an active network connection
	defaultMaximumIdleConnections        = 100              // Controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit
	defaultMaximumIdleConnectionsPerHost = 10               // if non-zero, controls the maximum idle (keep-alive) connections to keep per-host
	defaultIdleConnectionTimeout         = 90 * time.Second // Controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit
	defaultTLSHandshakeTimeout           = 10 * time.Second // Specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout
	defaultTLSInsecureSkipVerify         = false            // Enables skip verify TLS certificate
	defaultDialContextDualStack          = true             // Enables RFC 6555-compliant "Happy Eyeballs" dialing when the network is "tcp" and the host in the address parameter resolves to both IPv4 and IPv6 addresses
	requestChanBuffer                    = int(1000)        // Task channel buffer size
)

// ProxyFunc Is an a function to return a proxy for a given Request
type ProxyFunc func(*http.Request) (*url.URL, error)

// ErrorFunc Is an a client error retrieval function
type ErrorFunc func(err error)

// DebugFunc Is an a function for debug request/response data
type DebugFunc func(data []byte)

// DialTLSFunc Type of custom dial function for creating TLS connections for non-proxied HTTPS requests
type DialTLSFunc func(network, addr string) (net.Conn, error)

// DialContextFunc Type of custom dial function for creating unencrypted TCP connections
type DialContextFunc func(ctx context.Context, network, addr string) (net.Conn, error)

// is an implementation of transport
type impl struct {
	client                        *http.Client           // Объект http клиента
	transport                     *http.Transport        // Объект http транспорта
	cookieJar                     http.CookieJar         // Интерфейс CookieJar
	requestChan                   chan request.Interface // Канал задач запросов
	requestPoolLock               *sync.Mutex            // Лок от двойного запуска
	requestPoolCancelFunc         []context.CancelFunc   // Массив функций остановки пула воркеров
	requestPoolStarted            *atomic.Value          // =true Пул воркеров запущен
	requestPoolDone               *sync.WaitGroup        // WaitGroup для полной корректной остановки пула
	err                           error                  // Latest error
	errFunc                       ErrorFunc              // Колбэк функция получения ошибок на стороне http клиента
	debugFunc                     DebugFunc              // Is an a function for debug request/response data. If not nil - debug mode is enabled. If nil, debug mode is disbled
	methods                       methods.Interface      // Query Methods Interface
	requestPoolInterface          request.Pool           // Query objects pool interface
	requestPoolSize               uint16                 // Specifies a number of workers in the query pool
	proxy                         ProxyFunc              // Specifies a function to return a proxy for a given Request
	proxyConnectHeader            http.Header            // Optionally specifies headers to send to proxies during CONNECT requests
	dialContextTimeout            time.Duration          // Is the maximum amount of time a dial will wait for a connect to complete. The default is no timeout
	dialContextKeepAlive          time.Duration          // Specifies the keep-alive period for an active network connection. If zero, keep-alives are not enabled
	maximumIdleConnections        uint                   // Controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit
	maximumIdleConnectionsPerHost uint                   // If non-zero, controls the maximum idle (keep-alive) connections to keep per-host
	idleConnectionTimeout         time.Duration          // Is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself. Zero means no limit
	tlsHandshakeTimeout           time.Duration          // Specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout
	tlsInsecureSkipVerify         bool                   // Enables skip verify TLS certificate
	tlsClientConfig               *tls.Config            // Specifies the TLS configuration to use with tls.Client. If nil, the default configuration is used. If non-nil, HTTP/2 support may not be enabled by default
	tlsDialFunc                   DialTLSFunc            // Custom dial function for creating TLS connections for non-proxied HTTPS requests
	dialContextCustomFunc         DialContextFunc        // Custom dial function for creating unencrypted TCP connections
	dialContextDualStack          bool                   // Enables RFC 6555-compliant "Happy Eyeballs" dialing when the network is "tcp" and the host in the address parameter resolves to both IPv4 and IPv6 addresses
	totalTimeout                  time.Duration          // Specifies a time limit for requests made by this Client. The timeout includes connection time, any redirects, and reading the response body. A Timeout of zero means no timeout
}

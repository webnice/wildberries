package transport

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"crypto/tls"
	"net/http"
	"time"

	"gopkg.in/webnice/transport.v2/methods"
	"gopkg.in/webnice/transport.v2/request"
)

// Interface is an interface of object
type Interface interface {
	// Method Return interface of request methods
	Method() methods.Interface

	// RequestPoolSize Specifies a number of workers in the query pool
	RequestPoolSize(n uint16) Interface

	// ProxyFunc Specifies a function to return a proxy for a given Request
	ProxyFunc(f ProxyFunc) Interface

	// ProxyConnectHeader Optionally specifies headers to send to proxies during CONNECT requests
	ProxyConnectHeader(v http.Header) Interface

	// DialContextTimeout Is the maximum amount of time a dial will wait for a connect to complete. The default is no timeout
	DialContextTimeout(t time.Duration) Interface

	// DialContextKeepAlive Specifies the keep-alive period for an active network connection. If zero, keep-alives are not enabled
	DialContextKeepAlive(t time.Duration) Interface

	// MaximumIdleConnections Controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit
	MaximumIdleConnections(v uint) Interface

	// MaximumIdleConnectionsPerHost If non-zero, controls the maximum idle (keep-alive) connections to keep per-host
	MaximumIdleConnectionsPerHost(v uint) Interface

	// IdleConnectionTimeout Is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself. Zero means no limit
	IdleConnectionTimeout(t time.Duration) Interface

	// TLSHandshakeTimeout Specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout
	TLSHandshakeTimeout(t time.Duration) Interface

	// TLSSkipVerify Enables skip verify TLS certificate for all requests
	TLSSkipVerify(v bool) Interface

	// TLSClientConfig Specifies the TLS configuration to use with tls.Client.
	TLSClientConfig(v *tls.Config) Interface

	// DialTLS Specifies an custom dial function for creating TLS connections for non-proxied HTTPS requests
	DialTLS(fn DialTLSFunc) Interface

	// DialContextCustomFunc Specifies an custom dial function for creating unencrypted TCP connections
	DialContextCustomFunc(fn DialContextFunc) Interface

	// DualStack Enables RFC 6555-compliant "Happy Eyeballs" dialing when the network is "tcp" and the host in the address parameter resolves to both IPv4 and IPv6 addresses
	DualStack(v bool) Interface

	// TotalTimeout Specifies a time limit for requests made by this Client. The timeout includes connection time, any redirects, and reading the response body.
	// The timer remains running after Get, Head, Post, or Do return and will interrupt reading of the Response.Body.
	// A Timeout of zero means no timeout.
	TotalTimeout(t time.Duration) Interface

	// Transport Specifies of adjusted transport object
	Transport(tr *http.Transport) Interface

	// CookieJar Specifies of Cookie Jar interface
	CookieJar(v http.CookieJar) Interface

	// RequestGet Загрузка из sync.Pool объекта request и возврат интерфейса к нему
	// Полученный объект необходимо возвращать в sync.Pool методом RequestPut во избежании утечки памяти
	RequestGet() request.Interface

	// RequestPut Возврат в sync.Pool объекта request
	RequestPut(req request.Interface)

	// Client Returns the current http.Client
	Client() *http.Client

	// Do Executing the query in synchronous mode. Blocking function
	Do(req request.Interface) Interface

	// Done Stopping the worker pool, closing connections
	Done()

	// DebugFunc Set debug func and enable or disable debug mode
	// If fn=not nil - debug mode is enabled. If fn=nil, debug mode is disbled
	DebugFunc(fn DebugFunc) Interface

	// ERRORS

	// Error Return latest error
	Error() error

	// ErrorFunc Registering the error function on the client side
	ErrorFunc(fn ErrorFunc) Interface
}

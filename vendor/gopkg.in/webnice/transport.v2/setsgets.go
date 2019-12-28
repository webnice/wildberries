package transport

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"crypto/tls"
	"net/http"
	"time"
)

func setDefaults(trt *impl) {
	trt.requestPoolSize, trt.dialContextTimeout, trt.dialContextKeepAlive =
		defaultRequestPoolSize, defaultDialContextTimeout, defaultDialContextKeepAlive
	trt.maximumIdleConnections, trt.maximumIdleConnectionsPerHost, trt.idleConnectionTimeout, trt.tlsHandshakeTimeout, trt.tlsInsecureSkipVerify, trt.dialContextDualStack =
		defaultMaximumIdleConnections, defaultMaximumIdleConnectionsPerHost, defaultIdleConnectionTimeout, defaultTLSHandshakeTimeout, defaultTLSInsecureSkipVerify, defaultDialContextDualStack
	trt.ProxyFunc(http.ProxyFromEnvironment)
}

// RequestPoolSize Specifies a number of workers in the query pool
func (trt *impl) RequestPoolSize(v uint16) Interface {
	if v == 0 {
		return trt
	}
	trt.requestPoolSize = v
	return trt
}

// ProxyFunc Specifies a function to return a proxy for a given Request
func (trt *impl) ProxyFunc(f ProxyFunc) Interface {
	if f == nil {
		return trt
	}
	trt.proxy = f
	return trt
}

// ProxyConnectHeader Optionally specifies headers to send to proxies during CONNECT requests
func (trt *impl) ProxyConnectHeader(v http.Header) Interface {
	trt.proxyConnectHeader = v
	return trt
}

// DialContextTimeout Is the maximum amount of time a dial will wait for a connect to complete. The default is no timeout
func (trt *impl) DialContextTimeout(t time.Duration) Interface {
	trt.dialContextTimeout = t
	return trt
}

// DialContextKeepAlive Specifies the keep-alive period for an active network connection. If zero, keep-alives are not enabled
func (trt *impl) DialContextKeepAlive(t time.Duration) Interface {
	trt.dialContextKeepAlive = t
	return trt
}

// MaximumIdleConnections Controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit
func (trt *impl) MaximumIdleConnections(v uint) Interface {
	trt.maximumIdleConnections = v
	return trt
}

// MaximumIdleConnectionsPerHost If non-zero, controls the maximum idle (keep-alive) connections to keep per-host
func (trt *impl) MaximumIdleConnectionsPerHost(v uint) Interface {
	trt.maximumIdleConnectionsPerHost = v
	return trt
}

// IdleConnectionTimeout Is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself. Zero means no limit
func (trt *impl) IdleConnectionTimeout(t time.Duration) Interface {
	trt.idleConnectionTimeout = t
	return trt
}

// TLSHandshakeTimeout Specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout
func (trt *impl) TLSHandshakeTimeout(t time.Duration) Interface {
	trt.tlsHandshakeTimeout = t
	return trt
}

// TLSSkipVerify Enables skip verify TLS certificate for all requests
func (trt *impl) TLSSkipVerify(v bool) Interface {
	trt.tlsInsecureSkipVerify = v
	return trt
}

// TLSClientConfig Specifies the TLS configuration to use with tls.Client.
// If nil, the default configuration is used.
// If non-nil, HTTP/2 support may not be enabled by default
func (trt *impl) TLSClientConfig(v *tls.Config) Interface {
	trt.tlsClientConfig = v
	return trt
}

// DialTLS Specifies an custom dial function for creating TLS connections for non-proxied HTTPS requests
func (trt *impl) DialTLS(fn DialTLSFunc) Interface {
	trt.tlsDialFunc = fn
	return trt
}

// DialContextCustomFunc Specifies an custom dial function for creating unencrypted TCP connections
func (trt *impl) DialContextCustomFunc(fn DialContextFunc) Interface {
	trt.dialContextCustomFunc = fn
	return trt
}

// DualStack Enables RFC 6555-compliant "Happy Eyeballs" dialing when the network is "tcp" and the host in the address parameter resolves to both IPv4 and IPv6 addresses
func (trt *impl) DualStack(v bool) Interface {
	trt.dialContextDualStack = v
	return trt
}

// TotalTimeout Specifies a time limit for requests made by this Client. The timeout includes connection time, any redirects, and reading the response body.
// The timer remains running after Get, Head, Post, or Do return and will interrupt reading of the Response.Body.
// A Timeout of zero means no timeout.
func (trt *impl) TotalTimeout(t time.Duration) Interface {
	trt.totalTimeout = t
	return trt
}

// Transport Specifies of adjusted transport object
func (trt *impl) Transport(tr *http.Transport) Interface {
	trt.transport = tr
	return trt
}

// CookieJar Specifies of Cookie Jar interface
func (trt *impl) CookieJar(v http.CookieJar) Interface {
	trt.cookieJar = v
	return trt
}

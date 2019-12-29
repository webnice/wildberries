package request

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"gopkg.in/webnice/transport.v2/methods"
	"gopkg.in/webnice/web.v1/header"

	transportHeader "gopkg.in/webnice/transport.v2/header"
	"gopkg.in/webnice/transport.v2/response"
)

// Cancel Aborting request
func (r *Request) Cancel() Interface { r.contextCancelFunc(); return r }

// Done Waiting for the request to finish
func (r *Request) Done() Interface { <-r.context.Done(); return r }

// DoneWithContext Waiting for a request to complete, with the ability to interrupt the request through the context
func (r *Request) DoneWithContext(ctx context.Context) Interface {
	select {
	case <-r.context.Done():
	case <-ctx.Done():
		r.Cancel()
	}
	return r
}

// Error Return latest error
func (r *Request) Error() error { return r.err }

// DebugFunc Set debug func and enable or disable debug mode
// If fn=not nil - debug mode is enabled. If fn=nil, debug mode is disbled
func (r *Request) DebugFunc(fn DebugFunc) Interface {
	r.debugFunc = fn
	r.response.DebugFunc(response.DebugFunc(fn))
	return r
}

// Method Set request method
func (r *Request) Method(m methods.Value) Interface {
	if m == nil {
		r.err = fmt.Errorf("Warning, request method is nil, method not set, used default 'GET'")
		return r
	}
	r.method = m
	return r
}

// URL Set request URL
func (r *Request) URL(url string) Interface {
	r.url.Reset()
	r.url.WriteString(url)
	return r
}

// Referer Setting the referer header
func (r *Request) Referer(referer string) Interface {
	r.header.Add(header.Referer, referer)
	return r
}

// UserAgent Setting the UserAgent request header
func (r *Request) UserAgent(userAgent string) Interface {
	r.header.Add(header.UserAgent, userAgent)
	return r
}

// ContentType Setting the Content-Type request header
func (r *Request) ContentType(contentType string) Interface {
	r.header.Add(header.ContentType, contentType)
	return r
}

// Accept Setting the Accept request header
func (r *Request) Accept(accept string) Interface {
	r.header.Add(header.Accept, accept)
	return r
}

// AcceptEncoding Setting the Accept-Encoding request header
func (r *Request) AcceptEncoding(acceptEncoding string) Interface {
	r.header.Add(header.AcceptEncoding, acceptEncoding)
	return r
}

// AcceptLanguage Setting the Accept-Language request header
func (r *Request) AcceptLanguage(acceptLanguage string) Interface {
	r.header.Add(header.AcceptLanguage, acceptLanguage)
	return r
}

// Settings the Accept-Charset request header
func (r *Request) AcceptCharset(acceptCharset string) Interface {
	r.header.Add(header.AcceptCharset, acceptCharset)
	return r
}

// Settings custom request header
func (r *Request) CustomHeader(name string, value string) Interface {
	r.header.Add(name, value)
	return r
}

// BasicAuth Set login and password for request basic authorization
func (r *Request) BasicAuth(username string, password string) Interface {
	r.username, r.password = username, password
	return r
}

// Cookies Adding cookies to the request
func (r *Request) Cookies(cookies []*http.Cookie) Interface {
	r.cookie = append(r.cookie, cookies...)
	return r
}

// Header is an interface for custom headers manipulation
func (r *Request) Header() transportHeader.Interface { return r.header }

// Latency is an request latency without reading body of response
func (r *Request) Latency() time.Duration { return r.timeLatency }

// Response is an response interface
func (r *Request) Response() response.Interface { return r.response }

// DataStream Data for the request body in the form of an interface of io.Reader
func (r *Request) DataStream(data io.Reader) Interface {
	r.requestData, r.requestDataInterface = &bytes.Reader{}, data
	return r
}

// DataString Data for the request body in the form of an string
func (r *Request) DataString(data string) Interface {
	r.requestDataInterface, r.requestData = nil, bytes.NewReader([]byte(data))
	return r
}

// DataBytes Data for the request body in the form of an []byte
func (r *Request) DataBytes(data []byte) Interface {
	r.requestDataInterface, r.requestData = nil, bytes.NewReader(data)
	return r
}

// DataJSON The data for the query is created by serializing the object in JSON
func (r *Request) DataJSON(data interface{}) Interface {
	if r.tmpBytes, r.err = json.Marshal(data); r.err != nil {
		r.contextCancelFunc()
	}
	r.DataBytes(r.tmpBytes)
	return r
}

// DataXML The data for the query is created by serializing the object in XML
func (r *Request) DataXML(data interface{}) Interface {
	if r.tmpBytes, r.err = xml.Marshal(data); r.err != nil {
		r.contextCancelFunc()
	}
	r.DataBytes(r.tmpBytes)
	return r
}

// MakeRequest Создание запроса на основе метода запроса
func (r *Request) MakeRequest() (err error) {
	const getMethod = `GET`

	// Данные передаём через интерфейс, но только
	// - если интерфейс =nil
	// - если есть данные
	if r.requestDataInterface == nil && r.requestData.Len() > 0 {
		r.requestDataInterface = r.requestData
	}
	// Для метода GET, перенос данных в параметры URN
	// Если в URL нет `?` И есть данные
	if r.method.EqualFold(getMethod) && bytes.Index(r.url.Bytes(), []byte(`?`)) < 0 && r.requestData.Len() > 0 {
		if _, err = r.url.WriteString(`?`); err != nil {
			return
		}
		if _, err = r.requestData.WriteTo(r.url); err != nil {
			return
		}
		r.requestDataInterface = nil
	}
	r.request, err = http.NewRequestWithContext(r.context, r.method.String(), r.url.String(), r.requestDataInterface)

	return
}

// Request Returns the http.Request prepared for the request
func (r *Request) Request() (ret *http.Request, err error) {
	err = r.MakeRequest()
	ret = r.request
	return
}

// Do Executing the query and getting the result
func (r *Request) Do(client *http.Client) error {
	defer r.contextCancelFunc()

	// Создание запроса
	if r.err = r.MakeRequest(); r.err != nil {
		return r.err
	}
	// Заголовки простой авторизации
	if r.username != "" {
		r.request.SetBasicAuth(r.username, r.password)
	}
	// Печеньки запроса
	if len(r.cookie) > 0 {
		for r.tmpCounter = range r.cookie {
			r.request.AddCookie(r.cookie[r.tmpCounter])
		}
	}
	// Заголовки
	if r.header.Len() > 0 {
		r.tmpArr = r.header.Names()
		for r.tmpCounter = range r.tmpArr {
			if _, r.tmpOk = r.request.Header[r.tmpArr[r.tmpCounter]]; r.tmpOk {
				r.request.Header.Set(r.tmpArr[r.tmpCounter], r.header.Get(r.tmpArr[r.tmpCounter]))
			} else {
				r.request.Header.Add(r.tmpArr[r.tmpCounter], r.header.Get(r.tmpArr[r.tmpCounter]))
			}
		}
	}
	// Засекаем время запроса
	r.timeBegin = time.Now().In(time.Local)
	// Выполнение запроса
	r.err = r.response.Do(client, r.request)
	// Подсчитываем время ушедшее на запрос
	r.timeLatency = time.Since(r.timeBegin)
	if r.err != nil {
		return r.err
	}
	// Если включён дебаг
	if r.debugFunc != nil {
		if r.requestData.Size() > 0 {
			_, _ = r.requestData.Seek(0, io.SeekStart)
		}
		if buf, err := httputil.DumpRequestOut(r.request, true); err == nil {
			buf = bytes.Join([][]byte{[]byte("URI: " + r.url.String() + "\r\n"), buf}, []byte(``))
			r.debugRequest(buf)
		}
	}
	if r.response == nil {
		return fmt.Errorf("Request failed, response object is nil")
	}
	// Загрузка всех входящих данных
	if r.err = r.response.Load(); r.err != nil {
		return r.err
	}

	return r.err
}

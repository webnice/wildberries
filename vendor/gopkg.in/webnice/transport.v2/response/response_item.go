package response

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	runtimeDebug "runtime/debug"
	"time"

	"gopkg.in/webnice/transport.v2/charmap"
	"gopkg.in/webnice/transport.v2/content"
	"gopkg.in/webnice/transport.v2/data"
	"gopkg.in/webnice/transport.v2/header"
)

// DebugFunc Set debug func and enable or disable debug mode
// If fn=not nil - debug mode is enabled. If fn=nil, debug mode is disbled
func (r *Response) DebugFunc(fn DebugFunc) Interface { r.debugFunc = fn; return r }

// Do Выполнение запроса и получение Response
func (r *Response) Do(client *http.Client, request *http.Request) (err error) {
	r.response, err = client.Do(request)
	return
}

// Создание пути к месту хранения файла и полного имени временного файла
func (r *Response) makeTemporaryFileName() {
	r.tmpTm = time.Now().In(time.Local)
	r.contentFilename = path.Join(
		os.TempDir(),
		fmt.Sprintf("%020d.tmp", r.tmpTm.UnixNano()),
	)
	return
}

// Создание WriteCloser для загрузки результата запроса
func (r *Response) makeContentContainer(size int64) (err error) {
	if size < 0 || size > int64(maxDataSizeLoadedInMemory) {
		// Создание временного файла для данных
		r.makeTemporaryFileName()
		r.contentFh, err = os.OpenFile(r.contentFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0600))
		r.contentTemporaryFiles = append(r.contentTemporaryFiles, r.contentFilename)
		r.contentWriteCloser = r.contentFh
	} else {
		// Чтение в память
		r.contentData.Grow(int(r.response.ContentLength)) // Grow может паниковать (дебильный код)
		r.contentInMemory = true
		r.contentWriteCloser = data.NewWriteCloser(r.contentData)
	}

	return
}

// Создаёт единый интерфейс чтения загруженного контента
func (r *Response) makeContentReader() {
	if r.contentInMemory {
		r.contentReader = data.NewReadAtSeekerWriteToCloser(r.contentData)
	} else {
		if r.contentFh, r.err = os.OpenFile(r.contentFilename, os.O_RDONLY, os.FileMode(0600)); r.err != nil {
			return
		}
		r.contentReader = data.NewReadAtSeekerWriteToCloser(r.contentFh)
	}
}

// Load Load all response data
func (r *Response) Load() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Catch panic: %s\nGoroutine stack is:\n%s", e.(error), string(runtimeDebug.Stack()))
			return
		}
	}()
	// Создание интерфейса io.WriteCloser к результату запроса в памяти
	if err = r.makeContentContainer(r.response.ContentLength); err != nil {
		return
	}
	// Засекаем время выполнения загрузки
	r.timeBegin = time.Now().In(time.Local)
	// Загрузка данных
	r.contentLength, err = io.Copy(r.contentWriteCloser, r.response.Body)
	// Подсчитываем время ушедшее на загрузку
	r.timeLatency = time.Since(r.timeBegin)
	_ = r.contentWriteCloser.Close()
	_ = r.response.Body.Close()
	if err != nil {
		return
	}
	r.makeContentReader()
	// Если включён дебаг
	if r.debugFunc != nil {
		r.response.Body = r.contentReader
		if buf, err := httputil.DumpResponse(r.response, true); err == nil {
			r.debugResponse(buf)
		}
		_, _ = r.contentReader.Seek(0, io.SeekStart)
	}

	return
}

// Error Return latest error
func (r *Response) Error() error { return r.err }

// Response Returns the http.Response as is
func (r *Response) Response() *http.Response {
	_, _ = r.contentReader.Seek(0, io.SeekStart)
	r.response.Body = r.contentReader
	return r.response
}

// ContentLength records the length of the associated content
func (r *Response) ContentLength() int64 { return r.response.ContentLength }

// Cookies parses and returns the cookies set in the Set-Cookie headers
func (r *Response) Cookies() []*http.Cookie { return r.response.Cookies() }

// Latency is an request latency for reading body of response without reading header of response
func (r *Response) Latency() time.Duration { return r.timeLatency }

// StatusCode is an http status code of response
func (r *Response) StatusCode() int { return r.response.StatusCode }

// Status is an http status string of response, for known HTTP codes
func (r *Response) Status() string { return r.response.Status }

// Header maps header keys to values. If the response had multiple headers with the same key,
// they may be concatenated, with comma delimiters
func (r *Response) Header() header.Interface { return header.New(r.response.Header) }

// Charmap interface
func (r *Response) Charmap() charmap.Charmap { return r.charmap }

// Content() Interface for working with response content
func (r *Response) Content() content.Interface { return content.New(r.contentReader) }

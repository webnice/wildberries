package response

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
	"os"
	"sync"
	"time"

	"gopkg.in/webnice/transport.v2/charmap"
)

// New creates a new object and return interface
func New() Pool {
	var rsp = new(impl)
	rsp.responsePool = new(sync.Pool)
	rsp.responsePool.New = rsp.NewResponseItem
	return rsp
}

// NewResponseItem Конструктор sync.Pool для Response
func (rsp *impl) NewResponseItem() interface{} {
	var ret = &Response{
		contentData: &bytes.Buffer{},
		charmap:     charmap.NewCharmap(),
	}
	return ret
}

// ResponseGet Извлечение из pool нового элемента Response
func (rsp *impl) ResponseGet() Interface {
	return rsp.responsePool.Get().(*Response)
}

// ResponsePut Возврат в sync.Pool использованного элемента Response
func (rsp *impl) ResponsePut(req Interface) {
	rsp.responseClean(req.(*Response))
	rsp.responsePool.Put(req)
}

// Очистка данных объекта Response, подготовка к переиспользованию
func (rsp *impl) responseClean(r *Response) {
	r.err = nil
	r.response = nil
	r.debugFunc = nil
	r.timeBegin, r.timeLatency = time.Time{}, 0
	r.contentInMemory = false
	r.contentData.Reset()
	r.contentFilename = r.contentFilename[:0]
	if r.contentFh != nil {
		_ = r.contentFh.Close()
	}
	r.contentFh = nil
	r.contentWriteCloser = nil
	// Временные файлы
	for r.tmpI = range r.contentTemporaryFiles {
		_ = os.Remove(r.contentTemporaryFiles[r.tmpI])
	}
	r.contentTemporaryFiles = r.contentTemporaryFiles[:0]
	r.contentLength = 0
	if r.contentReader != nil {
		r.contentReader.Done()
	}
	r.contentReader = nil

	// Переменные для внутренних целей
	r.tmpOk, r.tmpTm, r.tmpString, r.tmpI = false, time.Time{}, r.tmpString[:0], 0
}

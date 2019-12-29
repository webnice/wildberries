package request

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bytes"
)

// Отображение запроса в дебаг режиме
func (r *Request) debugRequest(data []byte) {
	const prefixKey = `> `
	var (
		buf []byte
		tmp [][]byte
		i   int
	)

	defer func() { _ = recover() }()
	tmp, buf = bytes.Split(data, []byte("\n")), buf[:0]
	for i = range tmp {
		tmp[i] = bytes.TrimRight(tmp[i], "\r")
		buf = bytes.Join([][]byte{buf, []byte(prefixKey), tmp[i], []byte("\r\n")}, []byte(``))
	}
	r.debugFunc(bytes.TrimRight(buf, "\r\n"))
}

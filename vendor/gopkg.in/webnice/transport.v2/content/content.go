package content

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"gopkg.in/webnice/transport.v2/charmap"
	"gopkg.in/webnice/transport.v2/data"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// New creates a new object and return interface
func New(d data.ReadAtSeekerWriteToCloser) Interface { return &impl{esence: d} }

// WriteTo is an io.WriterTo interface implementation
func (cnt *impl) WriteTo(w io.Writer) (n int64, err error) {
	if err = cnt.ReaderCloser(); err != nil || cnt.rdc == nil {
		return
	}
	n, err = io.Copy(w, cnt.rdc)

	return
}

// ReaderCloser Получение io.ReadCloser для контента
func (cnt *impl) ReaderCloser() (err error) {
	// Разархивация ZIP
	if cnt.unzip {
		if cnt.rdc, err = cnt.UncompressZip(cnt.esence); err != nil {
			return
		}
	}
	// Разархивация TAR
	if cnt.untar {
		if cnt.rdc, err = cnt.UncompressTar(cnt.esence); err != nil {
			return
		}
	}
	// Разархивация GZIP
	if cnt.ungzip {
		if cnt.rdc, err = cnt.UncompressGzip(cnt.esence); err != nil {
			return
		}
	}
	// Разархивация FLATE
	if cnt.unflate {
		if cnt.rdc, err = cnt.UncompressFlate(cnt.esence); err != nil {
			return
		}
	}
	// Перекодирование контента если установлен транскодер
	if cnt.transcode != nil && cnt.esence != nil {
		// Создание ReadCloser из Reader + func Close
		cnt.rdc = data.NewReadCloser(
			transform.NewReader(cnt.esence, cnt.transcode.NewDecoder()),
			cnt.esence.Close,
		)
	} else if cnt.rdc == nil && cnt.esence != nil {
		cnt.rdc = cnt.esence
	}
	// Преобразование контента если установлен трансформер
	if cnt.transform != nil && cnt.rdc != nil {
		var newReader io.Reader
		newReader, err = cnt.transform(cnt.rdc)
		cnt.rdc = data.NewReadCloser(newReader, cnt.rdc.Close)
	}

	return
}

// UncompressZip Uncompress content as zip
func (cnt *impl) UncompressZip(r data.ReadAtSeekerWriteToCloser) (rdr io.ReadCloser, err error) {
	var zipReader *zip.Reader

	if zipReader, err = zip.NewReader(r, r.Size()); err != nil {
		err = fmt.Errorf("Zip archive error: %s", err.Error())
		return
	}
	if len(zipReader.File) <= 0 {
		err = fmt.Errorf("There are no files in the archive")
		return
	}
	if rdr, err = zipReader.File[0].Open(); err != nil {
		err = fmt.Errorf("Zip archive error, can't open file '%s': %s", zipReader.File[0].Name, err.Error())
		return
	}

	return
}

// UncompressTar Uncompress content as tar
func (cnt *impl) UncompressTar(r data.ReadAtSeekerWriteToCloser) (rdr io.ReadCloser, err error) {
	var tarReader *tar.Reader

	tarReader = tar.NewReader(r)
	_, err = tarReader.Next()
	if err == io.EOF {
		err = fmt.Errorf("There are no files in the archive")
		return
	}
	if err != nil {
		err = fmt.Errorf("Tar archive error: %s", err.Error())
		return
	}
	rdr = data.NewReadCloser(tarReader, r.Close)

	return
}

// UncompressGzip Uncompress content as gzip
func (cnt *impl) UncompressGzip(r data.ReadAtSeekerWriteToCloser) (rdr io.ReadCloser, err error) {
	var gzipReader *gzip.Reader

	if gzipReader, err = gzip.NewReader(r); err != nil && err != io.EOF {
		err = fmt.Errorf("GZIP content error: %s", err.Error())
		return
	} else if err == io.EOF {
		rdr, err = r, nil
		return
	}
	rdr = data.NewReadCloser(gzipReader, func() error { _ = gzipReader.Close(); return r.Close() })

	return
}

// UncompressFlate Uncompress content as flate
func (cnt *impl) UncompressFlate(r data.ReadAtSeekerWriteToCloser) (rdr io.ReadCloser, err error) {
	var flateReader io.ReadCloser

	if flateReader = flate.NewReader(r); flateReader == nil && err != io.EOF {
		err = fmt.Errorf("FLATE reader error")
		return
	} else if err == io.EOF {
		rdr, err = r, nil
		return
	}
	rdr = data.NewReadCloser(flateReader, func() error { _ = flateReader.Close(); return r.Close() })

	return
}

// String Получение контента в виде строки
func (cnt *impl) String() (ret string, err error) {
	var tmp = &bytes.Buffer{}

	if _, err = cnt.WriteTo(tmp); err != nil {
		return
	}
	ret = tmp.String()

	return
}

// Bytes Получение контента в виде среза байт
func (cnt *impl) Bytes() (ret []byte, err error) {
	var tmp = &bytes.Buffer{}

	if _, err = cnt.WriteTo(tmp); err != nil {
		return
	}
	ret = tmp.Bytes()

	return
}

// Transcode Перекодирование контента из указанной кодировки в UTF-8
func (cnt *impl) Transcode(from encoding.Encoding) Interface {
	return &impl{
		esence:    cnt.esence,
		transcode: from,
		transform: cnt.transform,
		untar:     cnt.untar,
		unzip:     cnt.unzip,
		ungzip:    cnt.ungzip,
		unflate:   cnt.unflate,
	}
}

// UnmarshalJSON Декодирование контента в структуру, предполагается что контент является json
func (cnt *impl) UnmarshalJSON(i interface{}) (err error) {
	var decoder *json.Decoder

	if err = cnt.ReaderCloser(); err == io.EOF {
		return nil
	} else if err != nil {
		return
	}
	defer func() { _ = cnt.rdc.Close() }()
	decoder = json.NewDecoder(cnt.rdc)
	err = decoder.Decode(i)

	return
}

// UnmarshalXML Декодирование контента в структуру, предполагается что контент является xml
func (cnt *impl) UnmarshalXML(i interface{}) (err error) {
	var decoder *xml.Decoder

	if err = cnt.ReaderCloser(); err == io.EOF {
		return nil
	} else if err != nil {
		return
	}
	defer func() { _ = cnt.rdc.Close() }()
	decoder = xml.NewDecoder(cnt.rdc)
	decoder.CharsetReader = cnt.MakeCharsetReader()
	err = decoder.Decode(i)

	return
}

// MakeCharsetReader Creating a function for streaming data reading with transcoding
func (cnt *impl) MakeCharsetReader() func(string, io.Reader) (io.Reader, error) {
	return func(cs string, input io.Reader) (rd io.Reader, err error) {
		// Перекодирование контента на уровне вышестоящего ридера
		if cnt.transcode != nil {
			rd = input
			return
		}
		// Поиск кодовой страницы
		var enc = charmap.NewCharmap().FindByName(cs)
		if enc == nil {
			err = fmt.Errorf("Could not find the code page '%s'", cs)
			return
		}
		// Новый ридер с перекодированием
		rd = data.NewReadCloser(
			transform.NewReader(input, enc.NewDecoder()),
			nil, // Поток будет закрыт в родительской функции, Closer не требуется
		)
		return
	}
}

// Transform Трансформирование исходного контента путём пропуска контента через переданный в функции ридер
func (cnt *impl) Transform(fn TransformFunc) Interface {
	return &impl{
		esence:    cnt.esence,
		transcode: cnt.transcode,
		transform: fn,
		untar:     cnt.untar,
		unzip:     cnt.unzip,
		ungzip:    cnt.ungzip,
		unflate:   cnt.unflate,
	}
}

// UnTar Разархивация контента методом TAR
func (cnt *impl) UnTar() Interface {
	return &impl{
		esence:    cnt.esence,
		transcode: cnt.transcode,
		transform: cnt.transform,
		untar:     true,
		unzip:     cnt.unzip,
		ungzip:    cnt.ungzip,
		unflate:   cnt.unflate,
	}
}

// UnZip Разархивация контента методом ZIP (извлекается только первый файл)
func (cnt *impl) UnZip() Interface {
	return &impl{
		esence:    cnt.esence,
		transcode: cnt.transcode,
		transform: cnt.transform,
		untar:     cnt.untar,
		unzip:     true,
		ungzip:    cnt.ungzip,
		unflate:   cnt.unflate,
	}
}

// UnGzip Разархивация контента методом GZIP
func (cnt *impl) UnGzip() Interface {
	return &impl{
		esence:    cnt.esence,
		transcode: cnt.transcode,
		transform: cnt.transform,
		untar:     cnt.untar,
		unzip:     cnt.unzip,
		ungzip:    true,
		unflate:   cnt.unflate,
	}
}

// UnFlate Разархивация контента методом FLATE
func (cnt *impl) UnFlate() Interface {
	return &impl{
		esence:    cnt.esence,
		transcode: cnt.transcode,
		transform: cnt.transform,
		untar:     cnt.untar,
		unzip:     cnt.unzip,
		ungzip:    cnt.ungzip,
		unflate:   true,
	}
}

// BackToBegin Returns the content reading pointer to the beginning
// This allows you to repeat the work with content
func (cnt *impl) BackToBegin() (err error) {
	if cnt.esence == nil {
		err = fmt.Errorf("request failed, response object is nil")
		return
	}
	_, err = cnt.esence.Seek(0, io.SeekStart)

	return
}

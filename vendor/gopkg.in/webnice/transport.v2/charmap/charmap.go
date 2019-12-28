//
// The "Program" is all-in-one
// This mirror golang.org/x/text/encoding/charmap library

package charmap

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// Преобразование различного написания кодировок к одному путём удаления символов разделителей
var rexReplaceCodepageName = regexp.MustCompile(`[-_ ]`)

// Charmap is an interface
type Charmap interface {
	FindByName(string) encoding.Encoding
	CodePage037() encoding.Encoding
	CodePage437() encoding.Encoding
	CodePage850() encoding.Encoding
	CodePage852() encoding.Encoding
	CodePage855() encoding.Encoding
	CodePage858() encoding.Encoding
	CodePage860() encoding.Encoding
	CodePage862() encoding.Encoding
	CodePage863() encoding.Encoding
	CodePage865() encoding.Encoding
	CodePage866() encoding.Encoding
	CodePage1047() encoding.Encoding
	CodePage1140() encoding.Encoding
	Iso8859_1() encoding.Encoding
	Iso8859_2() encoding.Encoding
	Iso8859_3() encoding.Encoding
	Iso8859_4() encoding.Encoding
	Iso8859_5() encoding.Encoding
	Iso8859_6() encoding.Encoding
	Iso8859_6E() encoding.Encoding
	Iso8859_6I() encoding.Encoding
	Iso8859_7() encoding.Encoding
	Iso8859_8() encoding.Encoding
	Iso8859_8E() encoding.Encoding
	Iso8859_8I() encoding.Encoding
	Iso8859_9() encoding.Encoding
	Iso8859_10() encoding.Encoding
	Iso8859_13() encoding.Encoding
	Iso8859_14() encoding.Encoding
	Iso8859_15() encoding.Encoding
	Iso8859_16() encoding.Encoding
	Koi8R() encoding.Encoding
	Koi8U() encoding.Encoding
	Macintosh() encoding.Encoding
	MacintoshCyrillic() encoding.Encoding
	Windows874() encoding.Encoding
	Windows1250() encoding.Encoding
	Windows1251() encoding.Encoding
	Windows1252() encoding.Encoding
	Windows1253() encoding.Encoding
	Windows1254() encoding.Encoding
	Windows1255() encoding.Encoding
	Windows1256() encoding.Encoding
	Windows1257() encoding.Encoding
	Windows1258() encoding.Encoding
	XUserDefined() encoding.Encoding
}

// implementation is an methods implementation
type implementation struct {
}

// NewCharmap Function create new charmap implementation
func NewCharmap() Charmap {
	return &implementation{}
}

// CodePage037 is the IBM Code Page 037 encoding
func (cm *implementation) CodePage037() encoding.Encoding {
	return charmap.CodePage037
}

// CodePage437 is the IBM Code Page 437 encoding
func (cm *implementation) CodePage437() encoding.Encoding {
	return charmap.CodePage437
}

// CodePage850 is the IBM Code Page 850 encoding
func (cm *implementation) CodePage850() encoding.Encoding {
	return charmap.CodePage850
}

// CodePage852 is the IBM Code Page 852 encoding
func (cm *implementation) CodePage852() encoding.Encoding {
	return charmap.CodePage852
}

// CodePage855 is the IBM Code Page 855 encoding
func (cm *implementation) CodePage855() encoding.Encoding {
	return charmap.CodePage855
}

// CodePage858 is the Windows Code Page 858 encoding
func (cm *implementation) CodePage858() encoding.Encoding {
	return charmap.CodePage858
}

// CodePage860 is the IBM Code Page 860 encoding
func (cm *implementation) CodePage860() encoding.Encoding {
	return charmap.CodePage860
}

// CodePage862 is the IBM Code Page 862 encoding
func (cm *implementation) CodePage862() encoding.Encoding {
	return charmap.CodePage862
}

// CodePage863 is the IBM Code Page 863 encoding
func (cm *implementation) CodePage863() encoding.Encoding {
	return charmap.CodePage863
}

// CodePage865 is the IBM Code Page 865 encoding
func (cm *implementation) CodePage865() encoding.Encoding {
	return charmap.CodePage865
}

// CodePage866 is the IBM Code Page 866 encoding
func (cm *implementation) CodePage866() encoding.Encoding {
	return charmap.CodePage866
}

// CodePage1047 is the IBM Code Page 1047 encoding
func (cm *implementation) CodePage1047() encoding.Encoding {
	return charmap.CodePage1047
}

// CodePage1140 is the IBM Code Page 1140 encoding
func (cm *implementation) CodePage1140() encoding.Encoding {
	return charmap.CodePage1140
}

// Iso8859_1 is the ISO 8859-1 encoding
func (cm *implementation) Iso8859_1() encoding.Encoding {
	return charmap.ISO8859_1
}

// Iso8859_2 is the ISO 8859-2 encoding
func (cm *implementation) Iso8859_2() encoding.Encoding {
	return charmap.ISO8859_2
}

// Iso8859_3 is the ISO 8859-3 encoding
func (cm *implementation) Iso8859_3() encoding.Encoding {
	return charmap.ISO8859_3
}

// Iso8859_4 is the ISO 8859-4 encoding
func (cm *implementation) Iso8859_4() encoding.Encoding {
	return charmap.ISO8859_4
}

// Iso8859_5 is the ISO 8859-5 encoding
func (cm *implementation) Iso8859_5() encoding.Encoding {
	return charmap.ISO8859_5
}

// Iso8859_6 is the ISO 8859-6 encoding
func (cm *implementation) Iso8859_6() encoding.Encoding {
	return charmap.ISO8859_6
}

// Iso8859_6E is the ISO 8859-6E encoding
func (cm *implementation) Iso8859_6E() encoding.Encoding {
	return charmap.ISO8859_6E
}

// Iso8859_6I is the ISO 8859-6I encoding
func (cm *implementation) Iso8859_6I() encoding.Encoding {
	return charmap.ISO8859_6I
}

// Iso8859_7 is the ISO 8859-7 encoding
func (cm *implementation) Iso8859_7() encoding.Encoding {
	return charmap.ISO8859_7
}

// Iso8859_8 is the ISO 8859-8 encoding
func (cm *implementation) Iso8859_8() encoding.Encoding {
	return charmap.ISO8859_8
}

// Iso8859_8E is the ISO 8859-8E encoding
func (cm *implementation) Iso8859_8E() encoding.Encoding {
	return charmap.ISO8859_8E
}

// Iso8859_8I is the ISO 8859-8I encoding
func (cm *implementation) Iso8859_8I() encoding.Encoding {
	return charmap.ISO8859_8I
}

// Iso8859_9 is the ISO 8859-9 encoding
func (cm *implementation) Iso8859_9() encoding.Encoding {
	return charmap.ISO8859_9
}

// Iso8859_10 is the ISO 8859-10 encoding
func (cm *implementation) Iso8859_10() encoding.Encoding {
	return charmap.ISO8859_10
}

// Iso8859_13 is the ISO 8859-13 encoding
func (cm *implementation) Iso8859_13() encoding.Encoding {
	return charmap.ISO8859_13
}

// Iso8859_14 is the ISO 8859-14 encoding
func (cm *implementation) Iso8859_14() encoding.Encoding {
	return charmap.ISO8859_14
}

// Iso8859_15 is the ISO 8859-15 encoding
func (cm *implementation) Iso8859_15() encoding.Encoding {
	return charmap.ISO8859_15
}

// Iso8859_16 is the ISO 8859-16 encoding
func (cm *implementation) Iso8859_16() encoding.Encoding {
	return charmap.ISO8859_16
}

// Koi8R is the KOI8-R encoding
func (cm *implementation) Koi8R() encoding.Encoding {
	return charmap.KOI8R
}

// Koi8U is the KOI8-U encoding
func (cm *implementation) Koi8U() encoding.Encoding {
	return charmap.KOI8U
}

// Macintosh is the Macintosh encoding
func (cm *implementation) Macintosh() encoding.Encoding {
	return charmap.Macintosh
}

// MacintoshCyrillic is the Macintosh Cyrillic encoding
func (cm *implementation) MacintoshCyrillic() encoding.Encoding {
	return charmap.MacintoshCyrillic
}

// Windows874 is the Windows 874 encoding
func (cm *implementation) Windows874() encoding.Encoding {
	return charmap.Windows874
}

// Windows1250 is the Windows 1250 encoding
func (cm *implementation) Windows1250() encoding.Encoding {
	return charmap.Windows1250
}

// Windows1251 is the Windows 1251 encoding
func (cm *implementation) Windows1251() encoding.Encoding {
	return charmap.Windows1251
}

// Windows1252 is the Windows 1252 encoding
func (cm *implementation) Windows1252() encoding.Encoding {
	return charmap.Windows1252
}

// Windows1253 is the Windows 1253 encoding
func (cm *implementation) Windows1253() encoding.Encoding {
	return charmap.Windows1253
}

// Windows1254 is the Windows 1254 encoding
func (cm *implementation) Windows1254() encoding.Encoding {
	return charmap.Windows1254
}

// Windows1255 is the Windows 1255 encoding
func (cm *implementation) Windows1255() encoding.Encoding {
	return charmap.Windows1255
}

// Windows1256 is the Windows 1256 encoding
func (cm *implementation) Windows1256() encoding.Encoding {
	return charmap.Windows1256
}

// Windows1257 is the Windows 1257 encoding
func (cm *implementation) Windows1257() encoding.Encoding {
	return charmap.Windows1257
}

// Windows1258 is the Windows 1258 encoding
func (cm *implementation) Windows1258() encoding.Encoding {
	return charmap.Windows1258
}

// XUserDefined is the X-User-Defined encoding
func (cm *implementation) XUserDefined() encoding.Encoding {
	return charmap.XUserDefined
}

// FindByName Поиск кодивой страницы по имени кодировки
func (cm *implementation) FindByName(name string) (ret encoding.Encoding) {
	var (
		sm, nm string
		item   encoding.Encoding
		i      int
	)

	sm = rexReplaceCodepageName.ReplaceAllString(strings.ToLower(name), ``)
	for i = range charmap.All {
		item, nm = charmap.All[i], rexReplaceCodepageName.ReplaceAllString(strings.ToLower(fmt.Sprintf("%s", item)), ``)
		if strings.EqualFold(sm, nm) {
			ret = item
		}
	}

	return
}

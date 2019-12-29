package types

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	wildberriesLocalLocation = `Europe/Moscow`
	fatalErrorInTime         = `Library error, parse location failed with error: `
	yearError                = `RFC 3339 is clear that years are 4 digits exactly. Year outside of range [0,9999]`
	quoteTime                = '"'

	// Форматы даты и времени валберис, стандартные RFC и локальные валберис
	wbLayout00 = time.RFC3339Nano                 // RFC3339 с указанием таймзоны и наносекундами
	wbLayout01 = time.RFC3339                     // RFC3339 с указанием таймзоны
	wbLayout02 = `2006-01-02T15:04:05.999999999`  // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout03 = `2006-01-02T15:04:05.99999999`   // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout04 = `2006-01-02T15:04:05.9999999`    // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout05 = `2006-01-02T15:04:05.999999`     // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout06 = `2006-01-02T15:04:05.99999`      // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout07 = `2006-01-02T15:04:05.9999`       // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout08 = `2006-01-02T15:04:05.999`        // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout09 = `2006-01-02T15:04:05.99`         // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout10 = `2006-01-02T15:04:05.9`          // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout11 = `2006-01-02T15:04:05.999999999Z` // Не RFC, UTC
	wbLayout12 = `2006-01-02T15:04:05.99999999Z`  // Не RFC, UTC
	wbLayout13 = `2006-01-02T15:04:05.9999999Z`   // Не RFC, UTC
	wbLayout14 = `2006-01-02T15:04:05.999999Z`    // Не RFC, UTC
	wbLayout15 = `2006-01-02T15:04:05.99999Z`     // Не RFC, UTC
	wbLayout16 = `2006-01-02T15:04:05.9999Z`      // Не RFC, UTC
	wbLayout17 = `2006-01-02T15:04:05.999Z`       // Не RFC, UTC
	wbLayout18 = `2006-01-02T15:04:05.99Z`        // Не RFC, UTC
	wbLayout19 = `2006-01-02T15:04:05.9Z`         // Не RFC, UTC
	wbLayout20 = `2006-01-02T15:04:05`            // Не RFC, локальная валберис (wildberriesLocalLocation)
	wbLayout21 = `2006-01-02T15:04:05Z`           // Не RFC, UTC
	wbLayout22 = `2006-01-02`                     // Не RFC, локальная валберис (wildberriesLocalLocation)
)

// WildberriesTimezoneLocal В соответствии с ответом технической поддержки валберис,
// дата и время (локальные) без указания таймзоны или смещения от UTC соответствуют Московскому времени
var WildberriesTimezoneLocal *time.Location

func init() {
	var err error

	if WildberriesTimezoneLocal, err = time.LoadLocation(wildberriesLocalLocation); err != nil {
		// В данном случае паника применима, так как ошибка возможна только если time.Time фатально повреждена
		panic(fatalErrorInTime + err.Error())
	}
}

// WildberriesTime Тип данных - конвертор даты валберис в дату полностью соответствующую RFC3339.
// В соответствии с ответом технической поддержки валберис, дата и время (локальные) без указания таймзоны
// соответствуют Московскому времени
type WildberriesTime struct {
	src []byte    // Исходная дата и время в представлении валберис
	fmt string    // Формат, который подошел при парсинге данных
	obj time.Time // Стандартизированная дата и время
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (wbt WildberriesTime) MarshalBinary() (ret []byte, err error) {
	ret, err = wbt.obj.MarshalBinary()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (wbt *WildberriesTime) UnmarshalBinary(data []byte) (err error) {
	err = wbt.obj.UnmarshalBinary(data)
	return
}

// GobEncode implements the gob.GobEncoder interface.
func (wbt WildberriesTime) GobEncode() ([]byte, error) { return wbt.MarshalBinary() }

// GobDecode implements the gob.GobDecoder interface.
func (wbt *WildberriesTime) GobDecode(data []byte) error { return wbt.UnmarshalBinary(data) }

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (wbt WildberriesTime) MarshalJSON() (ret []byte, err error) {
	var tmp []byte

	if tmp, err = wbt.MarshalText(); err != nil {
		return
	}
	ret = make([]byte, 0, len(time.RFC3339Nano)+2)
	ret = append(ret, quoteTime)
	ret = append(ret, tmp...)
	ret = append(ret, quoteTime)

	return
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (wbt *WildberriesTime) UnmarshalJSON(data []byte) (err error) {
	var (
		src     string
		n       int
		formats = []string{
			wbLayout00, wbLayout01, wbLayout02, wbLayout03, wbLayout04, wbLayout05, wbLayout06, wbLayout07, wbLayout08,
			wbLayout09, wbLayout10, wbLayout11, wbLayout12, wbLayout13, wbLayout14, wbLayout15, wbLayout16, wbLayout17,
			wbLayout18, wbLayout19, wbLayout20, wbLayout21, wbLayout22,
		}
	)

	wbt.src = make([]byte, len(data))
	_ = copy(wbt.src, data)
	src = strings.Trim(string(data), string(quoteTime))
	// Попытка распарсить дату
	for n = range formats {
		switch wbt.fmt = formats[n]; n {
		case 0, 1:
			wbt.obj, err = time.Parse(wbt.fmt, src)
		case 11, 12, 13, 14, 15, 16, 17, 18, 19, 21:
			wbt.obj, err = time.ParseInLocation(wbt.fmt, src, time.UTC)
		default:
			wbt.obj, err = time.ParseInLocation(wbt.fmt, src, WildberriesTimezoneLocal)
		}
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("failed to parse date and time %q", string(data))

	return
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in RFC 3339 format, with sub-second precision added if present.
func (wbt WildberriesTime) MarshalText() (ret []byte, err error) {
	var y int

	if wbt.fmt == "" {
		wbt.fmt = time.RFC3339Nano
		if wbt.obj.Hour() == 0 && wbt.obj.Minute() == 0 && wbt.obj.Second() == 0 {
			wbt.fmt = wbLayout01
		}
	}
	if y = wbt.obj.Year(); y < 0 || y >= 10000 {
		err = errors.New(yearError)
		return
	}
	ret = make([]byte, 0, len(time.RFC3339Nano))
	ret = wbt.obj.AppendFormat(ret, wbt.fmt)

	return
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (wbt *WildberriesTime) UnmarshalText(data []byte) error {
	return wbt.obj.UnmarshalJSON(data)
}

// String returns the time formatted using the format string "2006-01-02 15:04:05.999999999 -0700 MST"
func (wbt WildberriesTime) String() string { return wbt.obj.String() }

// Format returns a textual representation of the time value formatted
// according to layout, which defines the format by showing how the reference
// time, defined to be "Mon Jan 2 15:04:05 -0700 MST 2006"
func (wbt WildberriesTime) Format(layout string) string { return wbt.obj.Format(layout) }

// Time return object as a time.Time object
func (wbt WildberriesTime) Time() time.Time { return wbt.obj }

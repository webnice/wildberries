package orders

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/webnice/wildberries/modules/communication"
	wildberriesTypes "github.com/webnice/wildberries/types"

	"gopkg.in/webnice/transport.v2/request"
	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
)

// New creates a new object and return interface
func New(com communication.Interface, serverURI string, apiKey string, fromAt time.Time) Interface {
	var ods = &impl{
		ctx:       context.Background(),
		apiKey:    apiKey,
		serverURI: serverURI,
		com:       com,
		fromAt:    fromAt,
	}
	return ods
}

// WithContext Using context to interrupt requests to service
func (ods *impl) WithContext(ctx context.Context) Interface {
	if ctx == nil {
		return ods
	}
	ods.ctx = ctx
	return ods
}

// From Set of date and time of the beginning of the period for data request
func (ods *impl) From(fromAt time.Time) Interface {
	if fromAt.IsZero() {
		return ods
	}
	ods.fromAt = fromAt

	return ods
}

// Выбор значения fromAt
func (ods *impl) getFrom(fromAt ...time.Time) (ret time.Time) {
	var n int

	// Переопределения даты и времени начала периода для запроса, если fromAt передан
	ret = ods.fromAt
	for n = range fromAt {
		if fromAt[n].IsZero() {
			continue
		}
		ret = fromAt[n]
		break
	}

	return
}

// Report Load report data from the service.
// The onThisDay parameter indicates that data for the selected day is requested.
// If not set the fromAt parameter, then the data will be loaded for the current day
// or starting from the date and time set by the From function
func (ods *impl) Report(onThisDay bool, fromAt ...time.Time) (ret []*wildberriesTypes.Order, err error) {
	const (
		urn          = `%s/orders`
		keyDate      = `dateFrom`
		keyApi       = `key`
		keyOnThisDay = `flag`
		rawQueryFmt  = `%s=%s&%s=%s&%s=%s`
	)
	var (
		req     request.Interface
		from    time.Time
		uri     *url.URL
		flagKey string
	)

	// Подготовка данных
	from = ods.getFrom(fromAt...)
	if uri, err = url.Parse(fmt.Sprintf(urn, ods.serverURI)); err != nil {
		err = fmt.Errorf("can't create request URI, error: %s", err)
		return
	}
	if flagKey = "0"; onThisDay {
		flagKey = "1"
	}
	uri.RawQuery = fmt.Sprintf(
		rawQueryFmt,
		keyDate, from.In(wildberriesTypes.WildberriesTimezoneLocal).Format(wildberriesNonRFC3339TimeFormat),
		keyApi, url.QueryEscape(ods.apiKey),
		keyOnThisDay, flagKey,
	)
	// Создание запроса
	req = ods.com.NewRequestBaseJSON(uri.String(), ods.com.Transport().Method().Get())
	defer ods.com.Transport().RequestPut(req)
	req.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	// Выполнение запроса
	if err = ods.com.RequestResponseJSON(req, &ret); err != nil {
		err = fmt.Errorf("service response error: %s", err)
		return
	}

	return
}

package sales

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
	var sle = &impl{
		ctx:       context.Background(),
		apiKey:    apiKey,
		serverURI: serverURI,
		com:       com,
		fromAt:    fromAt,
	}
	return sle
}

// WithContext Using context to interrupt requests to service
func (sle *impl) WithContext(ctx context.Context) Interface {
	if ctx == nil {
		return sle
	}
	sle.ctx = ctx
	return sle
}

// From Set of date and time of the beginning of the period for data request
func (sle *impl) From(fromAt time.Time) Interface {
	if fromAt.IsZero() {
		return sle
	}
	sle.fromAt = fromAt

	return sle
}

// Выбор значения fromAt
func (sle *impl) getFrom(fromAt ...time.Time) (ret time.Time) {
	var n int

	// Переопределения даты и времени начала периода для запроса, если fromAt передан
	ret = sle.fromAt
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
// or starting from the date and time set by the From function.
// PriceWithDisc calculation formula:
//   Pricewithdisc = totalprice*((100 – discountPercent)/100 ) *((100 – promoCodeDiscount)/100 ) *((100 – spp)/100 )
func (sle *impl) Report(onThisDay bool, fromAt ...time.Time) (ret []*wildberriesTypes.Sale, err error) {
	const (
		urn          = `%s/sales`
		keyDate      = `dateFrom`
		keyApi       = `key`
		keyOnThisDay = `flag`
	)
	var (
		req  request.Interface
		from time.Time
		uri  *url.URL
	)

	// Подготовка данных
	from = sle.getFrom(fromAt...)
	if uri, err = url.Parse(fmt.Sprintf(urn, sle.serverURI)); err != nil {
		err = fmt.Errorf("can't create request URI, error: %s", err)
		return
	}
	uri.Query().Set(keyDate, from.Format(time.RFC3339))
	uri.Query().Set(keyApi, sle.apiKey)
	uri.Query().Set(keyOnThisDay, "0")
	if onThisDay {
		uri.Query().Set(keyOnThisDay, "1")
	}
	// Создание запроса
	req = sle.com.NewRequestBaseJSON(uri.String(), sle.com.Transport().Method().Get())
	defer sle.com.Transport().RequestPut(req)
	req.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	// Выполнение запроса
	if err = sle.com.RequestResponseJSON(req, &ret); err != nil {
		err = fmt.Errorf("service response error: %s", err)
		return
	}

	return
}

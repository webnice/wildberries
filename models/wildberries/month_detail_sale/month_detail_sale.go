package monthsale

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/webnice/wildberries/modules/communication"
	wildberriesTypes "github.com/webnice/wildberries/types"

	"gopkg.in/webnice/transport.v2/request"
	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
)

// New creates a new object and return interface
func New(com communication.Interface, serverURI string, apiKey string, fromAt time.Time) Interface {
	var mds = &impl{
		ctx:       context.Background(),
		apiKey:    apiKey,
		serverURI: serverURI,
		com:       com,
		fromAt:    fromAt,
	}
	return mds
}

// WithContext Using context to interrupt requests to service
func (mds *impl) WithContext(ctx context.Context) Interface {
	if ctx == nil {
		return mds
	}
	mds.ctx = ctx
	return mds
}

// From Set of date and time of the beginning of the period for data request
func (mds *impl) From(fromAt time.Time) Interface {
	if fromAt.IsZero() {
		return mds
	}
	mds.fromAt = fromAt

	return mds
}

// Выбор значения fromAt
func (mds *impl) getFrom(fromAt ...time.Time) (ret time.Time) {
	var n int

	// Переопределения даты и времени начала периода для запроса, если fromAt передан
	ret = mds.fromAt
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
// If not set the fromAt parameter, then the data will be loaded for the current day
// or starting from the date and time set by the From function
func (mds *impl) Report(
	rowID uint64,
	limit uint64,
	fromAt ...time.Time,
) (ret []*wildberriesTypes.MonthDetailSale, err error) {
	const (
		urn            = `%s/reportDetailMart`
		keyDate        = `dateFrom`
		keyApi         = `key`
		keyLimit       = `limit`
		keyRowID       = `rrdid`
		rawQueryFmt    = `%s=%s&%s=%s`
		rawQueryAddFmt = `&%s=%s`
	)
	var (
		req  request.Interface
		from time.Time
		uri  *url.URL
	)

	// Подготовка данных
	from = mds.getFrom(fromAt...)
	if uri, err = url.Parse(fmt.Sprintf(urn, mds.serverURI)); err != nil {
		err = fmt.Errorf("can't create request URI, error: %s", err)
		return
	}
	uri.RawQuery = fmt.Sprintf(
		rawQueryFmt,
		keyDate, from.In(wildberriesTypes.WildberriesTimezoneLocal).Format(wildberriesNonRFC3339TimeFormat),
		keyApi, url.QueryEscape(mds.apiKey),
	)
	if rowID > 0 {
		uri.RawQuery += fmt.Sprintf(rawQueryAddFmt, keyRowID, strconv.FormatUint(rowID, 10))
	}
	if limit > 0 {
		uri.RawQuery += fmt.Sprintf(rawQueryAddFmt, keyLimit, strconv.FormatUint(limit, 10))
	}
	// Создание запроса
	req = mds.com.NewRequestBaseJSON(uri.String(), mds.com.Transport().Method().Get())
	defer mds.com.Transport().RequestPut(req)
	req.Header().Add(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	// Выполнение запроса
	if err = mds.com.RequestResponseJSON(req, &ret); err != nil {
		err = fmt.Errorf("service response error: %s", err)
		return
	}

	return
}

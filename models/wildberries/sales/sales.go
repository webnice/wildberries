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

// UntilDone Configures repeated requests with a progressive timeout until a
// response is successfully received from the server, but not more than retryMax requests
func (sle *impl) UntilDone(retryTimeout time.Duration, retryMax uint) Interface {
	sle.retryTimeout, sle.retryMax = retryTimeout, retryMax
	return sle
}

// Выполнение запроса к серверу, получение и разбор результата
func (sle *impl) request(onThisDay bool, fromAt ...time.Time) (statusCode int, ret []*wildberriesTypes.Sale, err error) {
	const (
		urn          = `%s/sales`
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
	from = sle.getFrom(fromAt...)
	if uri, err = url.Parse(fmt.Sprintf(urn, sle.serverURI)); err != nil {
		err = fmt.Errorf("can't create request URI, error: %s", err)
		return
	}
	if flagKey = "0"; onThisDay {
		flagKey = "1"
	}
	uri.RawQuery = fmt.Sprintf(
		rawQueryFmt,
		keyDate, from.In(wildberriesTypes.WildberriesTimezoneLocal).Format(wildberriesNonRFC3339TimeFormat),
		keyApi, url.QueryEscape(sle.apiKey),
		keyOnThisDay, flagKey,
	)
	// Создание запроса
	req = sle.com.RequestJSON(sle.com.NewRequest(uri.String(), sle.com.Transport().Method().Get()))
	defer sle.com.Transport().RequestPut(req)
	// Выполнение запроса
	if statusCode, err = sle.com.RequestResponseJSON(sle.ctx, req, &ret); err != nil {
		err = fmt.Errorf("service response error: %s", err)
		return
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
	var (
		statusCode int
		n          uint
	)

	for {
		n++
		statusCode, ret, err = sle.request(onThisDay, fromAt...)
		// Успешный ответ
		if err == nil && (statusCode > 199 && statusCode < 300) {
			break
		}
		// Если выключены повторы или попытки кончились
		if sle.retryTimeout == 0 || sle.retryMax <= n {
			break
		}
		// Если было выполнено прерывание через контекст
		if err = sle.ctx.Err(); err != nil {
			break
		}
		// Ожидание прерывания или таймаута между повторами
		select {
		case <-time.After(sle.retryTimeout * time.Duration(n)):
		case <-sle.ctx.Done():
			err = sle.ctx.Err()
			break
		}
	}

	return
}

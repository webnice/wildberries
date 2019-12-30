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

// UntilDone Configures repeated requests with a progressive timeout until a
// response is successfully received from the server, but not more than retryMax requests
func (ods *impl) UntilDone(retryTimeout time.Duration, retryMax uint) Interface {
	ods.retryTimeout, ods.retryMax = retryTimeout, retryMax
	return ods
}

// Выполнение запроса к серверу, получение и разбор результата
func (ods *impl) request(onThisDay bool, fromAt ...time.Time) (statusCode int, ret []*wildberriesTypes.Order, err error) {
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
	req = ods.com.RequestJSON(ods.com.NewRequest(uri.String(), ods.com.Transport().Method().Get()))
	defer ods.com.Transport().RequestPut(req)
	// Выполнение запроса
	if statusCode, err = ods.com.RequestResponseJSON(ods.ctx, req, &ret); err != nil {
		err = fmt.Errorf("service response error: %s", err)
		return
	}

	return
}

// Report Load report data from the service.
// The onThisDay parameter indicates that data for the selected day is requested.
// If not set the fromAt parameter, then the data will be loaded for the current day
// or starting from the date and time set by the From function
func (ods *impl) Report(onThisDay bool, fromAt ...time.Time) (ret []*wildberriesTypes.Order, err error) {
	var (
		statusCode int
		n          uint
	)

	for {
		n++
		statusCode, ret, err = ods.request(onThisDay, fromAt...)
		// Успешный ответ
		if err == nil && (statusCode > 199 && statusCode < 300) {
			break
		}
		// Если выключены повторы или попытки кончились
		if ods.retryTimeout == 0 || ods.retryMax <= n {
			break
		}
		// Если было выполнено прерывание через контекст
		if err = ods.ctx.Err(); err != nil {
			break
		}
		// Ожидание прерывания или таймаута между повторами
		select {
		case <-time.After(ods.retryTimeout * time.Duration(n)):
		case <-ods.ctx.Done():
			err = ods.ctx.Err()
			break
		}
	}

	return
}

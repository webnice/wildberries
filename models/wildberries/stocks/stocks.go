package stocks

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
	var stk = &impl{
		ctx:       context.Background(),
		apiKey:    apiKey,
		serverURI: serverURI,
		com:       com,
		fromAt:    fromAt,
	}
	return stk
}

// WithContext Using context to interrupt requests to service
func (stk *impl) WithContext(ctx context.Context) Interface {
	if ctx == nil {
		return stk
	}
	stk.ctx = ctx
	return stk
}

// From Set of date and time of the beginning of the period for data request
func (stk *impl) From(fromAt time.Time) Interface {
	if fromAt.IsZero() {
		return stk
	}
	stk.fromAt = fromAt

	return stk
}

// Выбор значения fromAt
func (stk *impl) getFrom(fromAt ...time.Time) (ret time.Time) {
	var n int

	// Переопределения даты и времени начала периода для запроса, если fromAt передан
	ret = stk.fromAt
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
func (stk *impl) UntilDone(retryTimeout time.Duration, retryMax uint) Interface {
	stk.retryTimeout, stk.retryMax = retryTimeout, retryMax
	return stk
}

// Выполнение запроса к серверу, получение и разбор результата
func (stk *impl) request(fromAt ...time.Time) (statusCode int, ret []*wildberriesTypes.Stock, err error) {
	const (
		urn         = `%s/stocks`
		keyDate     = `dateFrom`
		keyApi      = `key`
		rawQueryFmt = `%s=%s&%s=%s`
	)
	var (
		req  request.Interface
		from time.Time
		uri  *url.URL
	)

	// Подготовка данных
	from = stk.getFrom(fromAt...)
	if uri, err = url.Parse(fmt.Sprintf(urn, stk.serverURI)); err != nil {
		err = fmt.Errorf("can't create request URI, error: %s", err)
		return
	}
	uri.RawQuery = fmt.Sprintf(
		rawQueryFmt,
		keyDate, from.In(wildberriesTypes.WildberriesTimezoneLocal).Format(wildberriesNonRFC3339TimeFormat),
		keyApi, url.QueryEscape(stk.apiKey),
	)
	// Создание запроса
	req = stk.com.RequestJSON(stk.com.NewRequest(uri.String(), stk.com.Transport().Method().Get()))
	defer stk.com.Transport().RequestPut(req)
	// Выполнение запроса
	if statusCode, err = stk.com.RequestResponseJSON(stk.ctx, req, &ret); err != nil {
		err = fmt.Errorf("service response error: %s", err)
		return
	}

	return
}

// Report Load report data from the service.
// If not set the fromAt parameter, then the data will be loaded for the current day
// or starting from the date and time set by the From function
func (stk *impl) Report(fromAt ...time.Time) (ret []*wildberriesTypes.Stock, err error) {
	var (
		statusCode int
		n          uint
	)

	for {
		n++
		statusCode, ret, err = stk.request(fromAt...)
		// Успешный ответ
		if err == nil && (statusCode > 199 && statusCode < 300) {
			break
		}
		// Если выключены повторы или попытки кончились
		if stk.retryTimeout == 0 || stk.retryMax <= n {
			break
		}
		// Если было выполнено прерывание через контекст
		if err = stk.ctx.Err(); err != nil {
			break
		}
		// Ожидание прерывания или таймаута между повторами
		select {
		case <-time.After(stk.retryTimeout * time.Duration(n)):
		case <-stk.ctx.Done():
			err = stk.ctx.Err()
			break
		}
	}

	return
}

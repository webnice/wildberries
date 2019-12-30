package incomes // import "incomes"

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
	var inc = &impl{
		ctx:       context.Background(),
		apiKey:    apiKey,
		serverURI: serverURI,
		com:       com,
		fromAt:    fromAt,
	}

	return inc
}

// WithContext Using context to interrupt requests to service
func (inc *impl) WithContext(ctx context.Context) Interface {
	if ctx == nil {
		return inc
	}
	inc.ctx = ctx

	return inc
}

// From Set of date and time of the beginning of the period for data request
func (inc *impl) From(fromAt time.Time) Interface {
	if fromAt.IsZero() {
		return inc
	}
	inc.fromAt = fromAt

	return inc
}

// Выбор значения fromAt
func (inc *impl) getFrom(fromAt ...time.Time) (ret time.Time) {
	var n int

	// Переопределения даты и времени начала периода для запроса, если fromAt передан
	ret = inc.fromAt
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
func (inc *impl) UntilDone(retryTimeout time.Duration, retryMax uint) Interface {
	inc.retryTimeout, inc.retryMax = retryTimeout, retryMax
	return inc
}

// Выполнение запроса к серверу, получение и разбор результата
func (inc *impl) request(fromAt ...time.Time) (statusCode int, ret []*wildberriesTypes.Income, err error) {
	const (
		urn         = `%s/incomes`
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
	from = inc.getFrom(fromAt...)
	if uri, err = url.Parse(fmt.Sprintf(urn, inc.serverURI)); err != nil {
		err = fmt.Errorf("can't create request URI, error: %s", err)
		return
	}
	uri.RawQuery = fmt.Sprintf(
		rawQueryFmt,
		keyDate, from.In(wildberriesTypes.WildberriesTimezoneLocal).Format(wildberriesNonRFC3339TimeFormat),
		keyApi, url.QueryEscape(inc.apiKey),
	)
	// Создание запроса
	req = inc.com.RequestJSON(inc.com.NewRequest(uri.String(), inc.com.Transport().Method().Get()))
	defer inc.com.Transport().RequestPut(req)
	// Выполнение запроса
	if statusCode, err = inc.com.RequestResponseJSON(inc.ctx, req, &ret); err != nil {
		err = fmt.Errorf("service response error: %s", err)
		return
	}

	return
}

// Report Load report data from the service.
// If not set the fromAt parameter, then the data will be loaded for the current day
// or starting from the date and time set by the From function
func (inc *impl) Report(fromAt ...time.Time) (ret []*wildberriesTypes.Income, err error) {
	var (
		statusCode int
		n          uint
	)

	for {
		n++
		statusCode, ret, err = inc.request(fromAt...)
		// Успешный ответ
		if err == nil && (statusCode > 199 && statusCode < 300) {
			break
		}
		// Если выключены повторы или попытки кончились
		if inc.retryTimeout == 0 || inc.retryMax <= n {
			break
		}
		// Если было выполнено прерывание через контекст
		if err = inc.ctx.Err(); err != nil {
			break
		}
		// Ожидание прерывания или таймаута между повторами
		select {
		case <-time.After(inc.retryTimeout * time.Duration(n)):
		case <-inc.ctx.Done():
			err = inc.ctx.Err()
			break
		}
	}

	return
}

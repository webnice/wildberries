package sales

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"time"

	"github.com/webnice/wildberries/modules/communication"
	wildberriesTypes "github.com/webnice/wildberries/types"
)

const wildberriesNonRFC3339TimeFormat = `2006-01-02T15:04:05`

// Interface is an interface of package
type Interface interface {
	// WithContext Using context to interrupt requests to service
	WithContext(ctx context.Context) Interface

	// From Set of date and time of the beginning of the period for data request
	From(fromAt time.Time) Interface

	// Report Load report data from the service.
	// The onThisDay parameter indicates that data for the selected day is requested.
	// If not set the fromAt parameter, then the data will be loaded for the current day
	// or starting from the date and time set by the From function.
	// PriceWithDisc calculation formula:
	//   Pricewithdisc = totalprice*((100 – discountPercent)/100 ) *((100 – promoCodeDiscount)/100 ) *((100 – spp)/100 )
	Report(onThisDay bool, fromAt ...time.Time) (ret []*wildberriesTypes.Sale, err error)
}

// impl is an implementation of package
type impl struct {
	fromAt    time.Time               // Дата и время начала периода для запроса данных
	com       communication.Interface // Интерфейс коммуникации с сервисом
	ctx       context.Context         // Интерфейс контекста
	apiKey    string                  // Ключ API
	serverURI string                  // URI адрес сервиса
}

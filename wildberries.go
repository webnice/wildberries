package wildberries

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"time"

	"github.com/webnice/wildberries/v1/models/wildberries/incomes"
	monthSale "github.com/webnice/wildberries/v1/models/wildberries/month_detail_sale"
	"github.com/webnice/wildberries/v1/models/wildberries/orders"
	"github.com/webnice/wildberries/v1/models/wildberries/sales"
	"github.com/webnice/wildberries/v1/models/wildberries/stocks"
	"github.com/webnice/wildberries/v1/modules/communication"
)

// New Creates an new object of package and return interface
func New(apiKey string) Interface {
	var (
		com    communication.Interface
		wbs    *impl
		uri    string
		fromAt time.Time
	)

	com, fromAt = communication.New(), time.Now().Truncate(time.Hour*24)
	uri = fmt.Sprintf(serviceURNv1, serviceURL)
	wbs = &impl{
		apiKey: apiKey,
		inc:    incomes.New(com, uri, apiKey, fromAt),
		ods:    orders.New(com, uri, apiKey, fromAt),
		sle:    sales.New(com, uri, apiKey, fromAt),
		stk:    stocks.New(com, uri, apiKey, fromAt),
		mds:    monthSale.New(com, uri, apiKey, fromAt),
		com:    com,
		fromAt: fromAt,
	}

	return wbs
}

// From Set of date and time of the beginning of the period for data request
func (wbs *impl) From(from time.Time) Interface {
	if from.IsZero() {
		return wbs
	}
	wbs.fromAt = from
	wbs.inc.From(wbs.fromAt)
	wbs.ods.From(wbs.fromAt)
	wbs.sle.From(wbs.fromAt)
	wbs.stk.From(wbs.fromAt)
	wbs.mds.From(wbs.fromAt)

	return wbs
}

// Incomes methods of reports about supply
func (wbs *impl) Incomes() incomes.Interface { return wbs.inc }

// Orders methods of reports about orders
func (wbs *impl) Orders() orders.Interface { return wbs.ods }

// Sales methods of reports about sales
func (wbs *impl) Sales() sales.Interface { return wbs.sle }

// Stocks methods of reports about warehouse
func (wbs *impl) Stocks() stocks.Interface { return wbs.stk }

// MonthDetailSale methods of reports about monthly sales
func (wbs *impl) MonthDetailSale() monthSale.Interface { return wbs.mds }

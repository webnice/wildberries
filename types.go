package wildberries

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	"github.com/webnice/wildberries/models/wildberries/incomes"
	monthSale "github.com/webnice/wildberries/models/wildberries/month_detail_sale"
	"github.com/webnice/wildberries/models/wildberries/orders"
	"github.com/webnice/wildberries/models/wildberries/sales"
	"github.com/webnice/wildberries/models/wildberries/stocks"
	"github.com/webnice/wildberries/modules/communication"
)

const (
	serviceURL   = ` https://suppliers-stats.wildberries.ru` // URL сервиса wildberries
	serviceURNv1 = `%s/api/v1/supplier`                      // URN адрес ресурса - версия ресурса
)

// Interface is an interface of package
type Interface interface {
	// Incomes methods of reports about supply
	Incomes() incomes.Interface

	// Orders methods of reports about orders
	Orders() orders.Interface

	// Sales methods of reports about sales
	Sales() sales.Interface

	// Stocks methods of reports about warehouse
	Stocks() stocks.Interface

	// MonthDetailSale methods of reports about monthly sales
	MonthDetailSale() monthSale.Interface
}

// impl is an implementation of package
type impl struct {
	fromAt time.Time
	apiKey string
	inc    incomes.Interface
	ods    orders.Interface
	sle    sales.Interface
	stk    stocks.Interface
	mds    monthSale.Interface
	com    communication.Interface
}

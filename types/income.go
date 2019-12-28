package types

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"
)

// Income Структура данных отчёта валберис "Поставки"
type Income struct {
	ID            uint64    `json:"incomeId"`        // Уникальный идентификатор поставки
	WbID          uint64    `json:"nmId"`            // Код валберис, он же номенклатура валберис, он же код 1С
	CreateAt      time.Time `json:"date"`            // Дата поставки
	ChangeAt      time.Time `json:"lastChangeDate"`  // Дата и время последнего обновления информации отчёта в сервисе
	AcceptAt      time.Time `json:"dateClose"`       // Дата и время принятия (закрытия) в валберис
	Quantity      uint64    `json:"quantity"`        // Количество
	Price         float64   `json:"totalPrice"`      // Цена товара из УПД
	Number        string    `json:"number"`          // Номер УПД
	VendorCode    string    `json:"supplierArticle"` // Артикул товара поставщика
	TechSize      string    `json:"techSize"`        // Технический размер
	Barcode       string    `json:"barcode"`         // Штрихкод
	WarehouseName string    `json:"warehouseName"`   // Название склада
}

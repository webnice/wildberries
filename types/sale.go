package types

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// Sale Структура данных отчёта валберис "Продажи"
type Sale struct {
	ID                string          `json:"number"`            // Номер документа
	CreateAt          WildberriesTime `json:"date"`              // Дата продажи
	ChangeAt          WildberriesTime `json:"lastChangeDate"`    // Дата и время последнего обновления информации отчёта в сервисе
	VendorCode        string          `json:"supplierArticle"`   // Артикул товара поставщика
	TechSize          string          `json:"techSize"`          // Технический размер
	Barcode           string          `json:"barcode"`           // Штрихкод
	Quantity          uint64          `json:"quantity"`          // Количество
	Price             float64         `json:"totalPrice"`        // Цена товара из УПД
	DiscountPercent   float64         `json:"discountPercent"`   // Согласованная итоговая скидка в процентах
	IsSupply          bool            `json:"isSupply"`          // Договор поставки
	IsRealization     bool            `json:"isRealization"`     // Договор реализации
	OrderID           uint64          `json:"orderId"`           // Уникальный идентификатор заказа - номер заказа из сервиса "заказы"
	DiscountPromoCode float64         `json:"promoCodeDiscount"` // Согласованная скидка по промо коду
	WarehouseName     string          `json:"warehouseName"`     // Название склада отгрузки товара
	CountryName       string          `json:"countryName"`       // Страна
	AreaName          string          `json:"oblastOkrugName"`   // Область или округ
	RegionName        string          `json:"regionName"`        // Регион
	IncomeID          uint64          `json:"incomeID"`          // Уникальный идентификатор поставки
	SaleID            string          `json:"saleID"`            // Уникальный идентификатор продажи или возврата S-продажа, R-возврат, D-доплата
	PositionID        uint64          `json:"odid"`              // Уникальный идентификатор позиции заказа
	Spp               float64         `json:"spp"`               // Согласованная скидка постоянного покупателя (СПП)
	Forpay            float64         `json:"forPay"`            // Сумма к перечислению поставщику
	FinishedPrice     float64         `json:"finishedPrice"`     // Фактическая цена из заказа с учётом всех скидок включая скидки валберис
	PriceWithDisc     float64         `json:"priceWithDisc"`     // Цена, от которойсчитается вознаграждение поставщика forpay, с учётом всех согласованных скидок
	WbID              uint64          `json:"nmId"`              // Код валберис, он же номенклатура валберис, он же код 1С
	Name              string          `json:"subject"`           // Предмет или название товара
	Category          string          `json:"category"`          // Категория
	BrandName         string          `json:"brand"`             // Бренд
}

package types

// Order Структура данных отчёта валберис "Заказы"
type Order struct {
	ID              uint64          `json:"number"`          // Уникальный идентификатор заказа
	CreateAt        WildberriesTime `json:"date"`            // Дата заказа
	ChangeAt        WildberriesTime `json:"lastChangeDate"`  // Дата и время последнего обновления информации отчёта в сервисе
	VendorCode      string          `json:"supplierArticle"` // Артикул товара поставщика
	TechSize        string          `json:"techSize"`        // Технический размер
	Barcode         string          `json:"barcode"`         // Штрихкод
	Quantity        uint64          `json:"quantity"`        // Количество
	Price           float64         `json:"totalPrice"`      // Цена товара из УПД
	DiscountPercent float64         `json:"discountPercent"` // Согласованная итоговая скидка в процентах
	WarehouseName   string          `json:"warehouseName"`   // Название склада отгрузки товара
	AreaName        string          `json:"oblast"`          // Область
	IncomeID        uint64          `json:"incomeID"`        // Уникальный идентификатор поставки
	PositionID      uint64          `json:"odid"`            // Уникальный идентификатор позиции заказа
	WbID            uint64          `json:"nmId"`            // Код валберис, он же номенклатура валберис, он же код 1С
	Name            string          `json:"subject"`         // Предмет или название товара
	Category        string          `json:"category"`        // Категория
	BrandName       string          `json:"brand"`           // Бренд
	IsCancel        bool            `json:"isCancel"`        // ???
}

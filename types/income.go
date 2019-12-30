package types

// Income Структура данных отчёта валберис "Поставки"
type Income struct {
	ID            uint64          `json:"incomeId"`        // Уникальный идентификатор поставки
	Number        string          `json:"number"`          // Номер УПД
	CreateAt      WildberriesTime `json:"date"`            // Дата поставки
	ChangeAt      WildberriesTime `json:"lastChangeDate"`  // Дата и время последнего обновления информации отчёта в сервисе
	VendorCode    string          `json:"supplierArticle"` // Артикул товара поставщика
	TechSize      string          `json:"techSize"`        // Технический размер
	Barcode       string          `json:"barcode"`         // Штрихкод
	Quantity      int64           `json:"quantity"`        // Количество
	Price         float64         `json:"totalPrice"`      // Цена товара из УПД
	AcceptAt      WildberriesTime `json:"dateClose"`       // Дата и время принятия (закрытия) в валберис
	WarehouseName string          `json:"warehouseName"`   // Название склада
	WbID          uint64          `json:"nmId"`            // Код валберис, он же номенклатура валберис, он же код 1С
}

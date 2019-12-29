package types

// Stock Структура данных отчёта валберис "Склад"
type Stock struct {
	ChangeAt            WildberriesTime `json:"lastChangeDate"`      // Дата и время последнего обновления информации отчёта в сервисе
	VendorCode          string          `json:"supplierArticle"`     // Артикул товара поставщика
	TechSize            string          `json:"techSize"`            // Технический размер
	Barcode             string          `json:"barcode"`             // Штрихкод
	Quantity            uint64          `json:"quantity"`            // Количество доступное для продажи - доступно на сайте, можно добавить в корзину
	IsSupply            bool            `json:"isSupply"`            // Договор поставки
	IsRealization       bool            `json:"isRealization"`       // Договор реализации
	QuantityFull        uint64          `json:"quantityFull"`        // Количество полное - то, что не продано (числится на складе)
	QuantityNotInOrders uint64          `json:"quantityNotInOrders"` // Количество не в заказе - числится на складе, и при этом не числится в незавершенном заказе
	WarehouseName       string          `json:"warehouseName"`       // Название склада
	InWayToClient       uint64          `json:"inWayToClient"`       // В пути к клиенту. штук
	InWayFromClient     uint64          `json:"inWayFromClient"`     // В пути от клиента, штук
	WbID                uint64          `json:"nmId"`                // Код валберис, он же номенклатура валберис, он же код 1С
	Name                string          `json:"subject"`             // Предмет или название товара
	Category            string          `json:"category"`            // Категория
	DaysOnSite          uint64          `json:"daysOnSite"`          // Количество дней на сайте
	BrandName           string          `json:"brand"`               // Бренд
}

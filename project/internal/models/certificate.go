package models

type Certificate struct {
	ID                int    `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	RealPrice         int    `json:"real_price" db:"real_price"`
	Discount          int    `json:"discount" db:"discount"`
	PriceWithDiscount int    `json:"price_with_discount" db:"price_with_discount"`
	Count             int    `json:"count" db:"count"`
	CountOfSales      int    `json:"count_of_sales" db:"count_of_sales"`
}

func (c *Certificate) AddSales() {
	c.Count--
	c.CountOfSales++
}

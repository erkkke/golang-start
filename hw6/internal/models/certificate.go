package models

type Certificate struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	RealPrice         int    `json:"real_price"`
	Discount          int    `json:"discount"`
	PriceWithDiscount int    `json:"price_with_discount"`
	Count             int    `json:"count"`
	CountOfSales      int    `json:"count_of_sales"`
}

func (c *Certificate) CalculatePriceWithDiscount() {
	c.PriceWithDiscount = c.RealPrice * (1 - c.Discount/100)
}

func (c *Certificate) CalculateDiscount() {
	c.Discount = (c.PriceWithDiscount / c.RealPrice) * 100
}

func (c *Certificate) AddSales() {
	c.Count--
	c.CountOfSales++
}

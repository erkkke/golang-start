package models

type Certificate struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	RealPrice         int    `json:"real_price"`
	Discount          int    `json:"discount"`
	PriceWithDiscount int    `json:"price_with_discount"`
	Count             int    `json:"count"`
	CountOfSales      int    `json:"count_of_sales"`
}

//func NewCertificate(name string, realPrice, discount, count int) *Certificate {
//	c := new(Certificate)
//	c.Name = name
//	c.RealPrice = realPrice
//	c.Discount = discount
//	c.SetDiscount(discount)
//	c.Count = count
//	return c
//}
//
//func (c *Certificate) SetDiscount(discount int) {
//	c.Discount = discount
//	c.PriceWithDiscount = c.RealPrice * (1 - discount / 100)
//}
//
//func (c *Certificate) AddSales() {
//	c.Count--
//	c.CountOfSales++
//}

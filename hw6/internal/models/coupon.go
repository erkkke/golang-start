package models

var couponNextID = 0

type Coupon struct {
	ID           int           `json:"id"`
	Category     string        `json:"category"`
	Company      string        `json:"company"`
	Name         string        `json:"name"`
	Certificates []Certificate `json:"certificates"`
	Address      string        `json:"address"`
	PhoneNumbers []string      `json:"phone_numbers"`
}

func (c *Coupon) NextID() {
	c.ID = couponNextID
	couponNextID++
}
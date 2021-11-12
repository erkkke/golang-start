package models

var couponNextID = 0

type Coupon struct {
	ID           int            `json:"id" db:"id"`
	Category     string         `json:"category" db:"category"`
	Company      string         `json:"company" db:"company"`
	Name         string         `json:"name" db:"name"`
	Description  string         `json:"description" db:"description"`
	Certificates []*Certificate `json:"certificates" db:"certificates"`
	Address      string         `json:"address" db:"address"`
	PhoneNumber  string         `json:"phone_numbers" db:"phone_numbers"`
}

func (c *Coupon) NextID() {
	c.ID = couponNextID
	couponNextID++
}

func (c *Coupon) GenerateCertificateID() {
	for i, certificate := range c.Certificates {
		certificate.ID = i
	}
}

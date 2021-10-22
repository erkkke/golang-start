package models

type Coupon struct {
	Id           int           `json:"id"`
	Category     string        `json:"category"`
	Company      string        `json:"company"`
	Name         string        `json:"name"`
	Certificates []Certificate `json:"certificates"`
}

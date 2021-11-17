package models

type Coupon struct {
	ID           int            `json:"id" db:"id"`
	Category     int            `json:"category" db:"category"`
	Company      string         `json:"company" db:"company"`
	Name         string         `json:"name" db:"name"`
	Description  string         `json:"description" db:"description"`
	Certificates []*Certificate `json:"certificates" db:"certificates"`
	Address      string         `json:"address" db:"address"`
	PhoneNumber  string         `json:"phone_number" db:"phone_number"`
}

type CouponsFilter struct {
	Query *string `json:"query"`
}

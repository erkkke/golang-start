package models

type Status string

const (
	Pending  Status = "pending"
	Approved Status = "approved"
	Canceled Status = "canceled"
)

type Order struct {
	Id            int    `json:"id" db:"id"`
	UserId        int    `json:"user_id" db:"user_id"`
	CouponId      int    `json:"coupon_id" db:"coupon_id"`
	CertificateId int    `json:"certificate_id" db:"certificate_id"`
	Status        Status `json:"status" db:"status"`
}

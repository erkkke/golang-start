package postgres

import (
	"context"
	"github.com/erkkke/golang-start/hw6/internal/models"
	"github.com/erkkke/golang-start/hw6/internal/store"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Coupons() store.CouponsRepository {
	if db.coupons == nil {
		db.coupons = NewCouponsRepository(db.conn)
	}

	return db.coupons
}

type CouponsRepository struct {
	conn *sqlx.DB
}

func NewCouponsRepository(conn *sqlx.DB) store.CouponsRepository {
	return &CouponsRepository{conn: conn}
}

func (c CouponsRepository) Create(ctx context.Context, coupon *models.Coupon) error {
	_, err := c.conn.ExecContext(ctx,
		"INSERT INTO coupons(category, company, name, description, address, phone_number) VALUES ($1, $2, 3$, 4$, $5, $6);",
		coupon.Category, coupon.Company, coupon.Name, coupon.Description, coupon.Address, coupon.PhoneNumber,
	)
	if err != nil {
		return err
	}

	var res []int
	if err = c.conn.Select(&res, "SELECT currval(pg_get_serial_sequence('coupons', 'id'))"); err != nil {
		return err
	}

	id := res[0]
	for _, v := range coupon.Certificates {
		_, err = c.conn.ExecContext(ctx, "INSERT INTO certificates(coupon_id, name, real_price, discount, price_with_discount, count, count_of_sales)" +
			"VALUES ($1, $2, 3$, 4$, $5, $6, $7)",
			id, v.Name, v.RealPrice, v.Discount, v.PriceWithDiscount, v.Count, v.CountOfSales)
	}

	return nil
}

func (c CouponsRepository) All(ctx context.Context) ([]*models.Coupon, error) {
	coupons := make([]*models.Coupon, 0)

	err := c.conn.SelectContext(ctx, &coupons,
		"SELECT coupons.id, category, company, coupons.name, description, certificate.id, certificate.name, certificate.real_price, certificate.discount, certificate.price_with_discount, certificate.count, certificate.count_of_sales, address, phone_number" +
		"FROM coupons" +
		"JOIN certificates" +
		"ON coupons.id = certificates.coupon_id")
	if err != nil {
		return nil, err
	}

	return coupons, nil
}

func (c CouponsRepository) ByID(ctx context.Context, id int) (*models.Coupon, error) {
	coupon := new(models.Coupon)
	if err := c.conn.GetContext(ctx, coupon, "SELECT * FROM coupons WHERE id=$1", id); err != nil {
		return nil, err
	}

	return coupon, nil
}

func (c CouponsRepository) Update(ctx context.Context, coupon *models.Coupon) error {
	_, err := c.conn.ExecContext(ctx, "UPDATE coupons SET category = $1, company = $2, name = $3, description = $4, adderess = $5, phone_number = $6 WHERE id=$7",
		coupon.Category, coupon.Company, coupon.Name, coupon.Description, coupon.Address, coupon.PhoneNumber, coupon.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c CouponsRepository) Delete(ctx context.Context, id int) error {
	if _, err := c.conn.ExecContext(ctx, "DELETE FROM coupons WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
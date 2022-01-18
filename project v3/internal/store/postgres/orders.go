package postgres

import (
	"context"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Orders() store.OrdersRepository {
	if db.orders == nil {
		db.orders = NewOrdersRepository(db.conn)
	}

	return db.orders
}

type OrdersRepository struct {
	conn *sqlx.DB
}

func NewOrdersRepository(conn *sqlx.DB) store.OrdersRepository {
	return &OrdersRepository{conn: conn}
}

func (o *OrdersRepository) Create(ctx context.Context, order *models.Order) error {
	_, err := o.conn.ExecContext(ctx, "INSERT INTO orders(user_id, coupon_id, certificate_id, status) VALUES ($1, $2, $3, $4)", order.UserId, order.CouponId, order.CertificateId, order.Status)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrdersRepository) All(ctx context.Context, filter *models.NameFilter) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)
	basicQuery := "SELECT * FROM orders"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE coupon_id=$1", basicQuery)
		if err := o.conn.SelectContext(ctx, &orders, basicQuery, filter.Query); err != nil {
			return nil, err
		}

		return orders, nil
	}

	if err := o.conn.SelectContext(ctx, &orders, basicQuery); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrdersRepository) AllOfUsers(ctx context.Context, userId int) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)

	err := o.conn.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE user_id=$1", userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrdersRepository) ById(ctx context.Context, id int) (*models.Order, error) {
	order := new(models.Order)

	if err := o.conn.SelectContext(ctx, order, "SELECT * FROM orders WHERE id=$1", id); err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrdersRepository) ChangeStatus(ctx context.Context, orderStatusDTO *models.OrderStatusDTO) error {
	_, err := o.conn.ExecContext(ctx, "UPDATE orders SET status=$1 WHERE id=$2", orderStatusDTO.Status, orderStatusDTO.Id)
	if err != nil {
		return err
	}

	return nil
}
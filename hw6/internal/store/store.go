package store

import (
	"context"
	"github.com/erkkke/golang-start/hw6/internal/models"
)

type Store interface {
	Create(ctx context.Context, coupon *models.Coupon) error
	All(ctx context.Context) ([]*models.Coupon, error)
	ByID(ctx context.Context, id int) (*models.Coupon, error)
	Update(ctx context.Context, coupon *models.Coupon) error
	Delete(ctx context.Context, id int) error
}

package store

import (
	"context"
	"github.com/erkkke/golang-start/hw6/internal/models"
)

type Store interface {
	Coupons() CouponsRepository
	Users() UserRepository
}

type CouponsRepository interface {
	Create(ctx context.Context, coupon *models.Coupon) error
	All(ctx context.Context) ([]*models.Coupon, error)
	ByID(ctx context.Context, id int) (*models.Coupon, error)
	Update(ctx context.Context, coupon *models.Coupon) error
	Delete(ctx context.Context, id int) error
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context) ([]*models.User, error)
	ByEmail(ctx context.Context, email string) (*models.User, error)
}

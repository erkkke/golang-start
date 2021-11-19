package cache

import (
	"context"
	"github.com/erkkke/golang-start/project/internal/models"
)

type Cache interface {
	Close() error

	Coupons() CouponsCacheRepo
	Users() UsersCacheRepo
	Categories() CategoriesCacheRepo

	DeleteAll(ctx context.Context) error
}

type CouponsCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Coupon) error
	Get(ctx context.Context, key string) ([]*models.Coupon, error)
}

type UsersCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.User) error
	Get(ctx context.Context, key string) ([]*models.User, error)
}

type CategoriesCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Category) error
	Get(ctx context.Context, key string) ([]*models.Category, error)
}
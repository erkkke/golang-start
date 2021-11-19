package store

import (
	"context"
	"github.com/erkkke/golang-start/project/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error

	Categories() CategoriesRepository
	Coupons() CouponsRepository
	Users() UsersRepository
}

type CategoriesRepository interface {
	Create(ctx context.Context, category *models.Category) error
	All(ctx context.Context, filter *models.CategoriesFilter) ([]*models.Category, error)
	ByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type CouponsRepository interface {
	Create(ctx context.Context, coupon *models.Coupon) error
	All(ctx context.Context, filter *models.NameFilter) ([]*models.Coupon, error)
	ByID(ctx context.Context, id int) (*models.Coupon, error)
	Update(ctx context.Context, coupon *models.Coupon) error
	Delete(ctx context.Context, id int) error
}

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context, filter *models.NameFilter) ([]*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

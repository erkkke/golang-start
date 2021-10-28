package inmemory

import (
	"context"
	"fmt"
	"github.com/erkkke/golang-start/hw6/internal/models"
	"sync"
)

type CouponsRepo struct {
	data map[int]*models.Coupon
	mu   *sync.RWMutex
}

func (db *CouponsRepo) Create(ctx context.Context, coupon *models.Coupon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	coupon.NextID()
	db.data[coupon.ID] = coupon
	return nil
}

func (db *CouponsRepo) All(ctx context.Context) ([]*models.Coupon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	coupons := make([]*models.Coupon, 0, len(db.data))
	for _, coupon := range db.data {
		coupons = append(coupons, coupon)
	}

	return coupons, nil
}

func (db *CouponsRepo) ByID(ctx context.Context, id int) (*models.Coupon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	coupon, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no coupon with id %d", id)
	}

	return coupon, nil
}

func (db *CouponsRepo) Update(ctx context.Context, coupon *models.Coupon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[coupon.ID] = coupon
	return nil
}

func (db *CouponsRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}

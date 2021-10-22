package inmemory

import (
	"context"
	"fmt"
	"github.com/erkkke/golang-start/hw6/internal/models"
	"github.com/erkkke/golang-start/hw6/internal/store"
	"sync"
)

type DB struct {
	data map[int]*models.Coupon
	mu   *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.Coupon),
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, coupon *models.Coupon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[coupon.Id] = coupon
	return nil
}

func (db *DB) All(ctx context.Context) ([]*models.Coupon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	coupons := make([]*models.Coupon, 0, len(db.data))
	for _, coupon := range db.data {
		coupons = append(coupons, coupon)
	}

	return coupons, nil
}

func (db *DB) ByID(ctx context.Context, id int) (*models.Coupon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	coupon, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no coupon with id %d", id)
	}

	return coupon, nil
}

func (db *DB) Update(ctx context.Context, coupon *models.Coupon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[coupon.Id] = coupon
	return nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}

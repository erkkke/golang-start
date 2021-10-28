package inmemory

import (
	"github.com/erkkke/golang-start/hw6/internal/models"
	"github.com/erkkke/golang-start/hw6/internal/store"
	"sync"
)

type DB struct {
	usersRepo   store.UserRepository
	couponsRepo store.CouponsRepository

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Users() store.UserRepository {
	if db.usersRepo == nil {
		db.usersRepo = &UsersRepo{
			data: make(map[string]*models.User),
			mu:   new(sync.RWMutex),
		}
	}
	return db.usersRepo
}

func (db *DB) Coupons() store.CouponsRepository {
	if db.couponsRepo == nil {
		db.couponsRepo = &CouponsRepo{
			data: make(map[int]*models.Coupon),
			mu:   new(sync.RWMutex),
		}
	}
	return db.couponsRepo
}

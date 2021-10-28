package inmemory

import (
	"context"
	"errors"
	"github.com/erkkke/golang-start/hw6/internal/models"
	"sync"
)

type UsersRepo struct {
	data map[string]*models.User
	mu   *sync.RWMutex
}

func (db *UsersRepo) Create(ctx context.Context, user *models.User) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeCreating(); err != nil {
		return err
	}
	user.Sanitize()
	user.NextId()
	db.data[user.Email] = user

	return nil
}

func (db *UsersRepo) All(ctx context.Context) ([]*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]*models.User, 0, len(db.data))
	for _, user := range db.data {
		users = append(users, user)
	}

	return users, nil
}

func (db *UsersRepo) ByEmail(ctx context.Context, email string) (*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	u, ok := db.data[email]
	if !ok {
		return nil, errors.New("not found")
	}

	return u, nil
}

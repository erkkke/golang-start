package postgres

import (
	"context"
	"github.com/erkkke/golang-start/hw6/internal/models"
	"github.com/erkkke/golang-start/hw6/internal/store"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Users() store.UsersRepository {
	if db.users == nil {
		db.users = NewUsersRepository(db.conn)
	}

	return db.users
}

type UsersRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) store.UsersRepository {
	return &UsersRepository{conn: conn}
}

func (u UsersRepository) Create(ctx context.Context, user *models.User) error {
	_, err := u.conn.ExecContext(ctx, "INSERT INTO users(email, phone_number, password, name, surname, birth_date, city) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.Email, user.PhoneNumber, user.Password, user.Name, user.Surname, user.BirthDate, user.City)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersRepository) All(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0)
	if err := u.conn.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UsersRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	if err := u.conn.GetContext(ctx, user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		return nil, err
	}

	return user, nil
}

func (u UsersRepository) Update(ctx context.Context, user *models.User) error {
	_, err := u.conn.ExecContext(ctx, "UPDATE users SET phone_number = $1, password = $2, name = $3, surname = $4, birth_date = $5, city = $6 WHERE id=$7",
		user.PhoneNumber, user.Password, user.Name, user.Surname, user.BirthDate, user.City, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersRepository) Delete(ctx context.Context, id int) error {
	if _, err := u.conn.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
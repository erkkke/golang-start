package postgres

import (
	"context"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
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
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeCreating(); err != nil {
		return err
	}

	_, err := u.conn.ExecContext(ctx, "INSERT INTO users(email, phone_number, password, name, surname, birth_date, city, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		user.Email, user.PhoneNumber, user.EncryptedPassword, user.Name, user.Surname, user.BirthDate, user.City, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersRepository) All(ctx context.Context, filter *models.NameFilter) ([]*models.User, error) {
	users := make([]*models.User, 0)
	basicQuery := "SELECT * FROM users"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)
		if err := u.conn.SelectContext(ctx, &users, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return users, nil
	}

	if err := u.conn.SelectContext(ctx, &users, basicQuery); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UsersRepository) Find(ctx context.Context, id int) (*models.User, error) {
	user := new(models.User)
	if err := u.conn.GetContext(ctx, user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return nil, err
	}

	return user, nil
}

func (u UsersRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	if err := u.conn.GetContext(ctx, user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		return nil, err
	}

	return user, nil
}

func (u UsersRepository) Update(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeCreating(); err != nil {
		return err
	}

	_, err := u.conn.ExecContext(ctx, "UPDATE users SET phone_number = $1, password = $2, name = $3, surname = $4, birth_date = $5, city = $6, role = $7 WHERE id=$8",
		user.PhoneNumber, user.EncryptedPassword, user.Name, user.Surname, user.BirthDate, user.City, user.Role, user.ID)
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
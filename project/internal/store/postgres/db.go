package postgres

import (
	"github.com/erkkke/golang-start/project/internal/store"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn *sqlx.DB

	coupons store.CouponsRepository
	categories store.CategoriesRepository
	users store.UsersRepository
}

func (db *DB) Connect(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}

	db.conn = conn
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func NewDB() store.Store {
	return &DB{}
}

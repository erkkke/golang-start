package postgres

import (
	"context"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Categories() store.CategoriesRepository {
	if db.categories == nil {
		db.categories = NewCategoriesRepository(db.conn)
	}

	return db.categories
}

type CategoriesRepository struct {
	conn *sqlx.DB
}

func NewCategoriesRepository(conn *sqlx.DB) store.CategoriesRepository {
	return &CategoriesRepository{conn: conn}
}

func (c *CategoriesRepository) Create(ctx context.Context, category *models.Category) error {
	_, err := c.conn.ExecContext(ctx, "INSERT INTO categories VALUES (default, $1)", category.Name)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoriesRepository) All(ctx context.Context, filter *models.CategoriesFilter) ([]*models.Category, error) {
	categories := make([]*models.Category, 0)
	basicQuery := "SELECT * FROM categories"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)

		if err := c.conn.Select(&categories, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}

		return categories, nil
	}

	if err := c.conn.Select(&categories, basicQuery); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CategoriesRepository) ById(ctx context.Context, id int) (*models.Category, error) {
	category := new(models.Category)
	if err := c.conn.GetContext(ctx, category, "SELECT id, name FROM categories WHERE id=$1", id); err != nil {
		return nil, err
	}

	return category, nil
}

func (c *CategoriesRepository) Update(ctx context.Context, category *models.Category) error {
	_, err := c.conn.ExecContext(ctx, "UPDATE categories SET name = $1 WHERE id=$2", category.Name, category.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoriesRepository) Delete(ctx context.Context, id int) error {
	if _, err := c.conn.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}

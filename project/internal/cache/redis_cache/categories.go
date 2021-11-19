package redis_cache

import (
	"context"
	"encoding/json"
	"github.com/erkkke/golang-start/project/internal/cache"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

func (rc RedisCache) Categories() cache.CategoriesCacheRepo {
	if rc.categories == nil {
		rc.categories = newCategoriesRepo(rc.client, rc.expires)
	}

	return rc.categories
}

type CategoriesRepo struct {
	client  *redis.Client
	expires time.Duration
}

func newCategoriesRepo(client *redis.Client, exp time.Duration) cache.CategoriesCacheRepo {
	return &CategoriesRepo{
		client:  client,
		expires: exp,
	}
}

func (c CategoriesRepo) Set(ctx context.Context, key string, value []*models.Category) error {
	categoryBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = c.client.Set(ctx, key, categoryBytes, c.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (c CategoriesRepo) Get(ctx context.Context, key string) ([]*models.Category, error) {
	result, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	categories := make([]*models.Category, 0)
	if err = json.Unmarshal([]byte(result), &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

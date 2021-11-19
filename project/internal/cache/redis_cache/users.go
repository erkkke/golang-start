package redis_cache

import (
	"context"
	"encoding/json"
	"github.com/erkkke/golang-start/project/internal/cache"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

func (rc *RedisCache) Users() cache.UsersCacheRepo {
	if rc.users == nil {
		rc.users = newUsersRepo(rc.client, rc.expires)
	}

	return rc.users
}

type UsersRepo struct {
	client  *redis.Client
	expires time.Duration
}

func newUsersRepo(client *redis.Client, expires time.Duration) cache.UsersCacheRepo {
	return &UsersRepo{
		client:  client,
		expires: expires,
	}
}

func (u UsersRepo) Set(ctx context.Context, key string, value []*models.User) error {
	couponBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = u.client.Set(ctx, key, couponBytes, u.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (u UsersRepo) Get(ctx context.Context, key string) ([]*models.User, error) {
	result, err := u.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	users := make([]*models.User, 0)
	if err = json.Unmarshal([]byte(result), &users); err != nil {
		return nil, err
	}

	return users, nil
}
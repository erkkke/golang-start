package redis_cache

import (
	"context"
	"encoding/json"
	"github.com/erkkke/golang-start/project/internal/cache"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

func (rc *RedisCache) Coupons() cache.CouponsCacheRepo {
	if rc.coupons == nil {
		rc.coupons = newCouponsRepo(rc.client, rc.expires)
	}

	return rc.coupons
}

type CouponsRepo struct {
	client  *redis.Client
	expires time.Duration
}

func newCouponsRepo(client *redis.Client, exp time.Duration) cache.CouponsCacheRepo {
	return &CouponsRepo{client: client, expires: exp}
}

func (c *CouponsRepo) Set(ctx context.Context, key string, value []*models.Coupon) error {
	couponBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = c.client.Set(ctx, key, couponBytes, c.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (c *CouponsRepo) Get(ctx context.Context, key string) ([]*models.Coupon, error) {
	result, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	coupons := make([]*models.Coupon, 0)
	if err = json.Unmarshal([]byte(result), &coupons); err != nil {
		return nil, err
	}

	return coupons, nil
}



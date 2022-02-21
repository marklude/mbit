package datastore

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type redisClient struct {
	ctx context.Context
	rdb *redis.Client
}

type Redis interface {
	SetPrice(price float64) (string, error)
}

func NewRedisClient(ctx context.Context, addr, password string) (Redis, error) {
	rdbClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &redisClient{ctx: ctx, rdb: rdbClient}, nil
}

func (r *redisClient) SetPrice(price float64) (result string, err error) {
	fPrice := strconv.FormatFloat(price, 'E', -1, 64)
	return r.rdb.Get(r.ctx, fPrice).Result()
}

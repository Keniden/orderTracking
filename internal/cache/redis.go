package cache

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

func NewRedis(addr string, password string, db int) *redis.Client {
    return redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
}

func Ping(ctx context.Context, rdb *redis.Client) error {
    _, err := rdb.Ping(ctx).Result()
    return err
}

func SetNX(ctx context.Context, rdb *redis.Client, key string, val string, ttl time.Duration) (bool, error) {
    return rdb.SetNX(ctx, key, val, ttl).Result()
}


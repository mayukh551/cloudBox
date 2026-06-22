package utils

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

type redisClient struct {
	client *redis.Client
}

func InitRedis() {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0 // default value
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})
}

func PutCacheObject(ctx context.Context, rdb *redis.Client, key string, payload any, ttl time.Duration) error {
	return rdb.HSet(ctx, key, payload).Err()
}

func GetCacheObject(ctx context.Context, rdb *redis.Client, key string, payload any, ttl time.Duration) error {
	return rdb.HSet(ctx, key, payload).Err()
}

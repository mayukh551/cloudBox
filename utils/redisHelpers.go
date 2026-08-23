package utils

import (
	"context"
	"encoding/json"
	"fmt"
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

	addr := os.Getenv("REDIS_ADDR")
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: addr + ":" + port,
		DB:   redisDB,
	})

	fmt.Println("Connected to Redis!")
}

func PutCacheObject(ctx context.Context, key string, payload any, ttl time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	return RDB.Set(ctx, key, data, ttl).Err()
}

func PutCacheHashFields(ctx context.Context, hashKey string, key string, payload any, ttl time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	RDB.HSet(ctx, hashKey, key, data, ttl)

	return nil
}

func GetCacheObject(ctx context.Context, key string, dest any) error {
	data, err := RDB.Get(ctx, key).Bytes()
	if err != nil {
		return err // redis.Nil on cache miss
	}
	return json.Unmarshal(data, dest)
}

func GetCacheValue(ctx context.Context, key string) (any, error) {
	return RDB.Get(ctx, key).Result()
}

func DelCache(ctx context.Context, key string) error {
	return RDB.Del(ctx, key).Err()
}

func DelAllCache(ctx context.Context, keys []string) error {
	return RDB.Del(ctx, keys...).Err()
}

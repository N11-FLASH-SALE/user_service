package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func ConnectDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redisemail:6379",
		Password: "root",
		DB:       0,
	})

	return rdb
}
func StoreCodes(ctx context.Context, code, email string) error {
	rdb := ConnectDB()

	err := rdb.Set(ctx, email, code, 10*time.Minute).Err()
	if err != nil {
		return errors.Wrap(err, "failed to set code in Redis")
	}

	return nil
}

func GetCodes(ctx context.Context, email string) (string, error) {
	rdb := ConnectDB()
	code, err := rdb.Get(ctx, email).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("no code found for email: %s", email)
		}
		return "", errors.Wrap(err, "failed to get code from Redis")
	}
	return code, nil
}

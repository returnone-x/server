package redisDB

import (
	"context"
	"github.com/returnone-x/server/config"
	"time"
)

func CreateStringData(key string, vaule string, exp time.Duration) error {
	ctx := context.Background()

	err := config.Redis.Set(ctx, key, vaule, exp).Err()
	if err != nil {
		return err
	}

	return nil
}
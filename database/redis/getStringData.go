package redisDB

import (
	"context"
	"github.com/returnone-x/server/config"
)

func GetStrigData(key string) (string, error){
	ctx := context.Background()

	value, err := config.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, err
}
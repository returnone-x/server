package redis

import (
	"context"
	"returnone/config"
)

func DeleteStringData(key string) (int64, error){
	ctx := context.Background()

	result, err := config.Redis.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, err
}
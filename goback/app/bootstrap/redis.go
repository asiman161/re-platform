package bootstrap

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func InitRD(dsn string) (*redis.Client, error) {
	dsn = strings.Trim(dsn, " ")
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "can't init redis")
	}

	return rdb, nil
}

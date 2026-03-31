package server

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/sig-agro/services/user-service/domain"
)

func InitRedis(conf *domain.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password:     conf.Redis.Pass,
		DB:           conf.Redis.Database,
		PoolSize:     conf.Redis.Pool,
		MinIdleConns: conf.Redis.MinIdleConn,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Msgf("Failed to connect to Redis: %v", err)
	}

	return rdb
}

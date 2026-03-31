package server

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/sig-agro/services/user-service/domain"
)

func InitPostgresDb(conf *domain.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		conf.Database.User, conf.Database.Pass, conf.Database.Host, conf.Database.Port, conf.Database.DBName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal().Msgf("Failed to parse postgres URL: %v", err)
	}

	config.MaxConns = int32(conf.Database.MaxConnPool)
	config.MinConns = int32(conf.Database.MaxIdleConn)
	config.MaxConnLifetime = time.Duration(conf.Database.ConnLifeTime) * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal().Msgf("Failed to create pgxpool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal().Msgf("Failed to ping database: %v", err)
	}

	return pool
}

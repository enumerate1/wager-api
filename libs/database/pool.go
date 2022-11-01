package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/wager-api/libs/configs"
	"github.com/wager-api/libs/try"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

var defaultMaxConnIdleTime = time.Minute * 5

// NewConnectionPool creates a new pool of connections to the database.
// It panics in case of error.
func NewConnectionPool(ctx context.Context, logger *zap.Logger, cfg configs.Postgres) *pgxpool.Pool {
	poolCfg, err := pgxpool.ParseConfig(cfg.Connection)
	if err != nil {
		logger.Panic("failed to parse database.connection: %s", zap.Error(err))
	}
	poolCfg.ConnConfig.LogLevel, err = pgx.LogLevelFromString(cfg.LogLevel)
	if err != nil {
		logger.Error("failed to parse database.log_level: %s", zap.Error(err))
		poolCfg.ConnConfig.LogLevel = pgx.LogLevelNone
	}
	if cfg.MaxConnIdleTime.Seconds() <= 0 {
		cfg.MaxConnIdleTime = defaultMaxConnIdleTime
	}
	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.ConnConfig.Logger = zapadapter.NewLogger(logger)
	logger.Info("database configuration",
		zap.Int32("MaxConns", cfg.MaxConns),
		zap.Int("RetryCount", cfg.RetryCount),
		zap.Float64("RetryInterval (seconds)", cfg.RetryInterval.Seconds()),
	)
	var pool *pgxpool.Pool
	err = try.Do(func(attempt int) (retry bool, err error) {
		pool, err = pgxpool.ConnectConfig(ctx, poolCfg)
		if err != nil {
			logger.Error("failed to connect to database", zap.Error(err), zap.Int("attempt", attempt))
			time.Sleep(cfg.RetryInterval)
			return attempt < cfg.RetryCount, err
		}
		return false, nil
	})
	if err != nil {
		if u, err2 := url.Parse(cfg.Connection); err2 != nil {
			logger.Panic("cannot create new connection to Postgres (failed to parse URI)", zap.Error(err))
		} else {
			logger.Panic(fmt.Sprintf("cannot create new connection to %q", u.Redacted()), zap.Error(err))
		}
	}
	return pool
}

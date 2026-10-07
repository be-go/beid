package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB 封装数据库连接池及其配置
type DB struct {
	Config *pgxpool.Config
	Pool   *pgxpool.Pool
}

// Close DB
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Open DB
func OpenDB(ctx context.Context, dsn string) (*DB, error) {
	var err error

	// Parse DSH params
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Other pool configs
	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute

	// [NOTICE] conn from pool is lazy connect
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// To connect once by Ping
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	// Init table "users" from script "users_init.sql", which requires idempotent
	_, err = pool.Exec(ctx, UsersInitSql)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{Config: config, Pool: pool}, nil
}

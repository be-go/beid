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

// NewDB 解析 DSN、创建连接池并验证连通性。
// 注意:NewWithConfig 本身是懒加载的(不会立即建连),这里通过 Ping 主动验证一次。
func NewDB(ctx context.Context, dsn string) (*DB, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// 连接池参数,可按需调整
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{Config: config, Pool: pool}, nil
}

// Close 关闭连接池,释放所有连接
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

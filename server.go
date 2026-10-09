package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db     *pgxpool.Pool // pgx pool
	router *gin.Engine   // gin router
	cfg    *Config       // server config
}

func NewServer(ctx context.Context) (*Server, error) {
	var err error

	// Load config
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// Init gin router
	gin.SetMode(cfg.Gin.Mode)
	router := gin.Default()

	// Parse DSN params
	pgconfig, err := pgxpool.ParseConfig(cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}

	// Other pool configs
	pgconfig.MaxConns = 20
	pgconfig.MinConns = 2
	pgconfig.MaxConnLifetime = time.Hour
	pgconfig.MaxConnIdleTime = time.Minute * 30
	pgconfig.HealthCheckPeriod = time.Minute

	// [NOTICE] conn from pool is lazy connect
	pool, err := pgxpool.NewWithConfig(ctx, pgconfig)
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
	_, err = pool.Exec(ctx, kUsersInitSql)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &Server{db: pool, router: router, cfg: cfg}, nil
}

func (svr *Server) Run(addr ...string) error {
	err := svr.router.Run(addr...)
	return err
}

func (svr *Server) Close() {
	if svr.db != nil {
		svr.db.Close()
	}
}

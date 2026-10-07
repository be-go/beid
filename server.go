package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
}

func NewServer(ctx context.Context) (*Server, error) {
	var err error

	// Init gin router
	router := gin.Default()

	// Parse DSH params
	dsn := "postgres://user:password@localhost:5432/beid?sslmode=disable"
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

	return &Server{router: router, db: pool}, nil
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

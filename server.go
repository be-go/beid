package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	ip     string        // server ip
	port   int           // server port
	cfg    *Config       // server's config
	router *gin.Engine   // gin router
	db     *pgxpool.Pool // pgx pool
}

func NewServer(ctx context.Context) (*Server, error) {
	var err error

	// Load config
	cfg, err := LoadConfigFrom("config.json")
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// Init gin router
	gin.SetMode(cfg.Gin.Mode)
	router := gin.Default()

	// Init pgx pool
	pgconfig, err := pgxpool.ParseConfig(cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}
	pgconfig.MaxConns = 20
	pgconfig.MinConns = 2
	pgconfig.MaxConnLifetime = time.Hour
	pgconfig.MaxConnIdleTime = time.Minute * 30
	pgconfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, pgconfig)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	// Init db from scripts, which requires idempotent
	_, err = pool.Exec(ctx, kUsersInitSql)
	if err != nil {
		pool.Close()
		return nil, err
	}

	// Construct and Return
	return &Server{
		cfg:    cfg,
		ip:     cfg.Server.IP,
		port:   cfg.Server.Port,
		router: router,
		db:     pool,
	}, nil
}

func (svr *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", svr.ip, svr.port)
	log.Printf("[Server] [Info] server running at '%s'.\n", addr)

	defer svr.Close()
	return svr.router.Run(addr)
}

func (svr *Server) Close() {
	if svr.db != nil {
		svr.db.Close()
	}
}

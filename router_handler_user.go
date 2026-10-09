package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func handle_post_user_register(svr *Server) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required,min=2,max=64"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6,max=256"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
		if err != nil {
			log.Printf("[Server] [ERROR] handle_post_user_register: hash password: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		var uid string
		err = svr.db.QueryRow(ctx.Request.Context(),
			`INSERT INTO users (username, email, password_hash)
			 VALUES ($1, $2, $3)
			 RETURNING id`,
			req.Username, req.Email, hash,
		).Scan(&uid)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				ctx.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
				return
			}

			log.Printf("[Server] [ERROR] handle_post_user_register: insert user: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"uid": uid})
	}
}

func handle_post_user_login(svr *Server) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required_without=Email"`
			Email    string `json:"email" binding:"required_without=Username"`
			Password string `json:"password" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Either identifier is accepted; email wins when both are supplied.
		query := `SELECT id, password_hash FROM users WHERE username = $1`
		account := req.Username
		if req.Email != "" {
			query = `SELECT id, password_hash FROM users WHERE email = $1`
			account = req.Email
		}

		var uid, hash string
		err := svr.db.QueryRow(ctx.Request.Context(), query, account).Scan(&uid, &hash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
				return
			}

			log.Printf("[Server] [ERROR] handle_post_user_login: query user: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		match, err := argon2id.ComparePasswordAndHash(req.Password, hash)
		if err != nil {
			log.Printf("[Server] [ERROR] handle_post_user_login: compare password: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		if !match {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"uid": uid})
	}
}

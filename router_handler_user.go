package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func handle_post_user_register(svr *Server) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req UserRegisterRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := ComputePasswdHash(req.Password)
		if err != nil {
			log.Printf("[Server] [Error] handle_post_user_register: hash password: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		var user_id string
		const query = "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id"
		err = svr.db.QueryRow(ctx.Request.Context(), query, req.Username, req.Email, hash).Scan(&user_id)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					ctx.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
					return
				}
			}
			log.Printf("[Server] [Error] handle_post_user_register: insert user: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"uid": user_id})
	}
}

func handle_post_user_login(svr *Server) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req UserLoginRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var query string
		var args []any
		if req.Username != "" && req.Email != "" {
			query = `SELECT id, password_hash FROM users WHERE username = $1 AND email = $2`
			args = []any{req.Username, req.Email}
		} else if req.Email != "" {
			query = `SELECT id, password_hash FROM users WHERE email = $1`
			args = []any{req.Email}
		} else {
			query = `SELECT id, password_hash FROM users WHERE username = $1`
			args = []any{req.Username}
		}

		var user_id, hash string
		err = svr.db.QueryRow(ctx.Request.Context(), query, args...).Scan(&user_id, &hash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
				return
			}
			log.Printf("[Server] [Error] handle_post_user_login: query user: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		match, err := ComparePasswdHash(req.Password, hash)
		if err != nil {
			log.Printf("[Server] [Error] handle_post_user_login: compare password: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		if !match {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"uid": user_id})
	}
}

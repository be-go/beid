package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=256"`
}

func handle_post_user_register(svr *Server) gin.HandlerFunc {
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

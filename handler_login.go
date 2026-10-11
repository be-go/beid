package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func handle_post_user_login(svr *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UserLoginRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user_id, hash string
		const query string = "SELECT id, password_hash FROM users WHERE email = $1"
		err = svr.db.QueryRow(ctx.Request.Context(), query, req.Email).Scan(&user_id, &hash)
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

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handle_get_test(svr *Server) func(*gin.Context) {
	return func(ctx *gin.Context) {
		status := http.StatusOK
		data := make(map[string]any)

		if err := svr.db.Ping(ctx.Request.Context()); err != nil {
			data["db"] = "down"
			data["db_error"] = err.Error()
			status = http.StatusServiceUnavailable
		} else {
			data["db"] = "up"
			status = http.StatusOK
		}

		ctx.JSON(status, data)
	}
}

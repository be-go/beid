package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	server, err := NewServer(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	server.router.GET("/test", func(c *gin.Context) {
		status := http.StatusOK
		data := make(map[string]any)

		if err := server.db.Ping(c.Request.Context()); err != nil {
			data["db"] = "down"
			data["db_error"] = err.Error()
			status = http.StatusServiceUnavailable
		} else {
			data["db"] = "up"
			status = http.StatusOK
		}

		c.JSON(status, data)
	})

	server.Run(":8080")
	server.Close()
}

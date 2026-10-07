package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	GlobalInit()
	gin.SetMode(gin.DebugMode)
	server := gin.Default()

	server.GET("/test", func(c *gin.Context) {
		status := http.StatusOK
		data := make(map[string]any)

		if err := gDatabase.Pool.Ping(c.Request.Context()); err != nil {
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
	GlobalQuit()
}

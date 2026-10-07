package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// db 全局数据库句柄:main 中初始化后,整个进程内的 handler 都可以直接使用。
// pgxpool.Pool 本身是并发安全的;全局变量在启动时写入一次,之后只读,不存在数据竞争。
var db *DB

func main() {
	var err error
	db, err = NewDB(context.Background(), "postgres://user:password@localhost:5432/beid?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"phrase": "OK", "test": true})
	})

	// 示例:在 handler 中通过全局 db 使用连接池
	router.GET("/health", func(c *gin.Context) {
		if err := db.Pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"db": "down", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"db": "up"})
	})

	router.Run(":8080")
}

package main

import (
	"context"
	"log"
)

func main() {
	server, err := NewServer(context.Background())
	if err != nil {
		log.Fatalf("[Server] [Fatal] %v\n", err)
	}

	test_router(server) // For test, should remove in release

	server.router.POST("/user/register", handle_post_user_register(server))
	server.router.POST("/user/login", handle_post_user_login(server))

	server.Run()
}

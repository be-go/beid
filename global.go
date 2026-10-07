package main

import (
	"context"
	"log"
)

var gDatabase *DB // 全局 Postgres 数据库句柄

func GlobalInit() {
	var err error
	ctx := context.Background()

	gDatabase, err = OpenDB(ctx, "postgres://user:password@localhost:5432/beid?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

}

func GlobalQuit() {
	if gDatabase != nil {
		gDatabase.Close()
	}
}

package main

import (
	_ "embed"
)

//go:embed .env/secret.data
var kHS256Secret string

//go:embed sql/users_init.sql
var kUsersInitSql string

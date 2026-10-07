package main

import (
	_ "embed"
)

//go:embed sql/users_init.sql
var UsersInitSql string

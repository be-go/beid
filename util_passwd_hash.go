package main

import "github.com/alexedwards/argon2id"

func ComputePasswdHash(passwd string) (string, error) {
	return argon2id.CreateHash(passwd, argon2id.DefaultParams)
}

func ComparePasswdHash(passwd string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(passwd, hash)
}

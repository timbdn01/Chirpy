package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	//takes a password and returns a hashed version of the password using argon2id
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	//takes a password and a hash and returns true if the password matches the hash, false otherwise
	return argon2id.ComparePasswordAndHash(password, hash)
}
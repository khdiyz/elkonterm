package helper

import (
	"crypto/sha1"
	"elkonterm/config"
	"errors"
	"fmt"
)

func GenerateHash(cfg *config.Config, password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot empty")
	}

	salt := cfg.HashKey
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))), nil
}

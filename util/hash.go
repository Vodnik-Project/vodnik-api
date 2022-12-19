package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func PassHash(password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't generate password: %v", err)
	}
	return string(passHash), nil
}

package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func PassHash(password string) string {
	passHash := sha256.Sum256([]byte(password))
	passHashHex := hex.EncodeToString(passHash[:])
	return passHashHex
}

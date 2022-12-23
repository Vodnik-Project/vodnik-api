package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// generate sessionID hash
func GetSessionID(userAgent, acceptLang string) string {
	sessionIDHash := sha256.Sum256([]byte(userAgent + acceptLang))
	sessionID := hex.EncodeToString(sessionIDHash[:])
	return sessionID
}

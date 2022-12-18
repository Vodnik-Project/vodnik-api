package util

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"

	"github.com/labstack/echo/v4"
)

// generate sessionID hash
func GetSessionID(userAgent, acceptLang string) string {
	sessionIDHash := sha256.Sum256([]byte(userAgent + acceptLang))
	sessionID := hex.EncodeToString(sessionIDHash[:])
	return sessionID
}

// get field from access token payload parsed from jwt middleware
func GetFieldFromPayload(c echo.Context, field string) string {
	payload := reflect.ValueOf(c.Get("user")).Elem()
	claims := payload.FieldByName("Claims").Elem()
	username := claims.Elem().FieldByName(field)

	return username.String()
}

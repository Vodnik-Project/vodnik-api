package auth

import (
	"fmt"
	"time"
)

type AccessTokenPayload struct {
	UserID    string
	Username  string
	IssuedAt  int64
	ExpiresAt int64
}

type RefreshTokenPayload struct {
	IssuedAt  int64
	ExpiresAt int64
}

func NewAccessTokenPayload(userid string, username string, duration time.Duration) AccessTokenPayload {
	p := AccessTokenPayload{
		UserID:    userid,
		Username:  username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	}
	return p
}
func NewRefreshTokenPayload(duration time.Duration) RefreshTokenPayload {
	p := RefreshTokenPayload{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	}
	return p
}

func (p RefreshTokenPayload) Valid() error {
	if time.Now().Unix() > p.ExpiresAt {
		return fmt.Errorf("payload is not valid")
	}
	return nil
}

func (p AccessTokenPayload) Valid() error {
	if time.Now().Unix() > p.ExpiresAt {
		return fmt.Errorf("payload is not valid")
	}
	return nil
}

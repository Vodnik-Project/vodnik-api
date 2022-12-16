package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccessTokenPayload struct {
	User_ID   uuid.UUID
	IssuedAt  int64
	ExpiresAt int64
}

type RefreshTokenPayload struct {
	ID        string
	IssuedAt  int64
	ExpiresAt int64
}

func NewAccessTokenPayload(userID uuid.UUID, duration time.Duration) AccessTokenPayload {
	p := AccessTokenPayload{
		User_ID:   userID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	}
	return p
}
func NewRefreshTokenPayload(sessionID string, duration time.Duration) RefreshTokenPayload {
	p := RefreshTokenPayload{
		ID:        sessionID,
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

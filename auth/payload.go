package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccessTokenPayload struct {
	ID        uuid.UUID
	User_ID   uuid.UUID
	IssuedAt  int64
	ExpiresAt int64
}

type RefreshTokenPayload struct {
	ID        uuid.UUID
	IssuedAt  int64
	ExpiresAt int64
}

func NewAccessTokenPayload(userID uuid.UUID, duration time.Duration) AccessTokenPayload {
	p := AccessTokenPayload{
		ID:        uuid.New(),
		User_ID:   userID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	}
	return p
}
func NewRefreshTokenPayload(duration time.Duration) RefreshTokenPayload {
	p := RefreshTokenPayload{
		ID:        uuid.New(),
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

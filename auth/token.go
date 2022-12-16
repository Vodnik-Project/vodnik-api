package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Token struct {
	Secret               []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type TokenMaker interface {
	CreateAccessToken(userID uuid.UUID) (string, error)
	CreateRefreshToken(sessionID string) (string, error)
}

func NewTokenMaker(token Token) Token {
	t := Token{
		Secret:               token.Secret,
		AccessTokenDuration:  token.AccessTokenDuration,
		RefreshTokenDuration: token.RefreshTokenDuration,
	}
	return t
}

func (t Token) CreateAccessToken(userID uuid.UUID) (string, error) {
	p := NewAccessTokenPayload(userID, t.AccessTokenDuration)
	if err := p.Valid(); err != nil {
		return "", fmt.Errorf("payload is not valid: %v", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	tk, err := token.SignedString(t.Secret)
	if err != nil {
		return "", fmt.Errorf("can't create token: %v", err)
	}
	return tk, nil
}

func (t Token) CreateRefreshToken(sessionID string) (string, error) {
	p := NewRefreshTokenPayload(sessionID, t.RefreshTokenDuration)
	if err := p.Valid(); err != nil {
		return "", fmt.Errorf("payload is not valid: %v", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	tk, err := token.SignedString(t.Secret)
	if err != nil {
		return "", fmt.Errorf("can't create token: %v", err)
	}
	return tk, nil
}

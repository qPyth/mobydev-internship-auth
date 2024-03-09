package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenManager interface {
	NewJWT(userID uint) (string, error)
}

type Manager struct {
	secretKey string
	tokenTTL  time.Duration
}

func NewManager(secretKey string, tokenTTL time.Duration) *Manager {
	return &Manager{secretKey: secretKey, tokenTTL: tokenTTL}
}

func (m *Manager) NewJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(m.tokenTTL).Unix(),
	})

	return token.SignedString([]byte(m.secretKey))
}

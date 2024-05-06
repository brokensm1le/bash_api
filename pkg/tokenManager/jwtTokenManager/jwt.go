package jwtTokenManager

import (
	"bash_api/pkg/tokenManager"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

type Manager struct {
	signingKey string
}

func NewManger(signingKey string) (tokenManager.TokenManager, error) {
	if signingKey == "" {
		return nil, errors.New("no signing key")
	}
	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(data *tokenManager.Data, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenManager.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id:   data.Id,
		Role: data.Role,
	})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (*tokenManager.Data, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenManager.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid singing method")
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return &tokenManager.Data{}, err
	}

	claims, ok := token.Claims.(*tokenManager.CustomClaims)
	if !ok {
		return &tokenManager.Data{}, fmt.Errorf("invalid claims type")
	}

	return &tokenManager.Data{Id: claims.Id, Role: claims.Role}, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	r := rand.New(rand.NewSource(time.Now().Unix()))

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

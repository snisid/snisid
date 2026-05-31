package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snisid/platform/backend/internal/config"
)

type Claims struct {
	Role   string `json:"role"`
	Agency string `json:"agency"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret     string
	ttl        time.Duration
	refreshTTL time.Duration
}

func NewJWTService(cfg config.JWTConfig) *JWTService {
	return &JWTService{
		secret:     cfg.Secret,
		ttl:        cfg.TTL,
		refreshTTL: cfg.RefreshTTL,
	}
}

func (s *JWTService) SignToken(subject, role, agency string) (string, error) {
	now := time.Now()
	claims := Claims{
		Role:   role,
		Agency: agency,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.secret))
}

func (s *JWTService) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

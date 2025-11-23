package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
)

// Service handles JWT encoding/decoding.
type Service struct {
	secret string
	issuer string
	ttl    time.Duration
}

// Claims carry user metadata inside the token.
type Claims struct {
	UserID string              `json:"userId"`
	Role   entity.UserRole     `json:"role"`
	Exam   entity.ExamCategory `json:"exam,omitempty"`
	jwtlib.RegisteredClaims
}

// NewService creates a JWT helper.
func NewService(secret, issuer string, ttl time.Duration) *Service {
	return &Service{secret: secret, issuer: issuer, ttl: ttl}
}

// Generate signs claims into a token.
func (s *Service) Generate(claims Claims) (string, error) {
	if claims.RegisteredClaims.Issuer == "" {
		claims.RegisteredClaims.Issuer = s.issuer
	}

	if claims.RegisteredClaims.ExpiresAt == nil {
		claims.RegisteredClaims.ExpiresAt = jwtlib.NewNumericDate(time.Now().Add(s.ttl))
	}

	if claims.RegisteredClaims.IssuedAt == nil {
		claims.RegisteredClaims.IssuedAt = jwtlib.NewNumericDate(time.Now())
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// Parse validates the token and returns claims.
func (s *Service) Parse(token string) (*Claims, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}

	parsed, err := jwtlib.ParseWithClaims(token, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token claims")
	}

	if _, err := uuid.Parse(claims.UserID); err != nil {
		return nil, fmt.Errorf("invalid user id in token: %w", err)
	}

	return claims, nil
}

// internal/adapters/services/auth_service.go
package services

import (
	"context"
	"iam-service/internal/core/domain"
	"iam-service/internal/core/ports"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	secretKey []byte
	redis     ports.CacheRepository
}

func NewAuthService(redis ports.CacheRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	return &AuthService{
		secretKey: []byte(secret),
		redis:     redis,
	}
}

func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *AuthService) InvalidateToken(ctx context.Context, token string) error {
	// Add token to blacklist with 24h expiry
	return s.redis.SetOTP(ctx, "blacklist:"+token, "invalid", time.Hour*24)
}

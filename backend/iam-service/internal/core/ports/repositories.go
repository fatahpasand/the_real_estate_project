// internal/core/ports/repositories.go
package ports

import (
	"context"
	"iam-service/internal/core/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByPhone(ctx context.Context, phone string) (*domain.User, error) // New method
	Update(ctx context.Context, user *domain.User) error
	UpdateVerificationStatus(ctx context.Context, userID uint, verified bool) error
}

type CacheRepository interface {
	SetOTP(ctx context.Context, key, value string, expiration time.Duration) error
	GetOTP(ctx context.Context, key string) (string, error)
	DeleteOTP(ctx context.Context, key string) error
}

type AuditRepository interface {
	LogAudit(ctx context.Context, log *domain.AuditLog) error
}

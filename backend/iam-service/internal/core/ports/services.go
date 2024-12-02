// internal/core/ports/services.go
package ports

import (
	"context"
	"iam-service/internal/core/domain"
)

type EmailService interface {
	SendVerificationEmail(to, otp string) error
	SendLoginAlert(to, ip, userAgent string) error
}

type AuthService interface {
	GenerateToken(user *domain.User) (string, error)
	InvalidateToken(ctx context.Context, token string) error
}

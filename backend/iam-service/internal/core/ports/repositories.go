// internal/core/ports/repositories.go
package ports

import "github.com/fatahpasand/backend/iam-service/internal/core/domain"

type UserRepository interface {
    Create(user *domain.User) error
    FindByEmail(email string) (*domain.User, error)
    FindByID(id uint) (*domain.User, error)
    UpdateVerificationStatus(userID uint, verified bool) error
    StoreVerificationOTP(email, otp string) error
    VerifyOTP(email, otp string) bool
}
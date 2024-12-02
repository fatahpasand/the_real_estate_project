// internal/core/usecases/user_usecase.go
package usecases

import (
	"context"
	"errors"
	"strconv"
	"time"

	"iam-service/internal/core/domain"
	"iam-service/internal/core/ports"
	"iam-service/pkg/utils"
)

// UserUseCase interface defines the contract for user operations
type UserUseCase interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	VerifyEmail(ctx context.Context, email, otp string) error
	UpdateProfile(ctx context.Context, userID string, name string, phone string) error
	GetCache() ports.CacheRepository
}

// userUseCaseImpl implements UserUseCase interface
type userUseCaseImpl struct {
	userRepo ports.UserRepository
	cache    ports.CacheRepository
	audit    ports.AuditRepository
	email    ports.EmailService
	auth     ports.AuthService
}

// NewUserUseCase creates a new UserUseCase instance
func NewUserUseCase(ur ports.UserRepository, cr ports.CacheRepository, ar ports.AuditRepository,
	es ports.EmailService, as ports.AuthService) UserUseCase {
	return &userUseCaseImpl{
		userRepo: ur,
		cache:    cr,
		audit:    ar,
		email:    es,
		auth:     as,
	}
}

func (uc *userUseCaseImpl) GetCache() ports.CacheRepository {
	return uc.cache
}

func (uc *userUseCaseImpl) Register(ctx context.Context, user *domain.User) error {
	// Check if email exists
	if _, err := uc.userRepo.FindByEmail(ctx, user.Email); err == nil {
		return errors.New("email already exists")
	}

	// Check if phone exists (if provided)
	if user.Phone != "" {
		if _, err := uc.userRepo.FindByPhone(ctx, user.Phone); err == nil {
			return errors.New("phone number already exists")
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Generate OTP
	otp := utils.GenerateOTP()
	if err := uc.cache.SetOTP(ctx, user.Email, otp, time.Minute*15); err != nil {
		return err
	}

	// Create user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Send verification email
	return uc.email.SendVerificationEmail(user.Email, otp)
}

func (uc *userUseCaseImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := utils.ComparePasswords(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.Verified {
		return "", errors.New("email not verified")
	}

	// Generate JWT token
	token, err := uc.auth.GenerateToken(user)
	if err != nil {
		return "", err
	}

	// Log successful login
	go uc.LogAudit(ctx, user.ID, "login", "success", "", "")

	return token, nil
}

func (uc *userUseCaseImpl) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Assuming we need to add FindByID to UserRepository interface and implement it
	user, err := uc.userRepo.FindByID(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	// Clear sensitive data before returning
	user.Password = ""
	return user, nil
}

func (uc *userUseCaseImpl) VerifyEmail(ctx context.Context, email, otp string) error {
	storedOTP, err := uc.cache.GetOTP(ctx, email)
	if err != nil || storedOTP != otp {
		return errors.New("invalid OTP")
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	return uc.userRepo.UpdateVerificationStatus(ctx, user.ID, true)
}

func (uc *userUseCaseImpl) LogAudit(ctx context.Context, userID uint, action, status, ip, userAgent string) {
	log := &domain.AuditLog{
		UserID:    userID,
		Action:    action,
		Status:    status,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	// Log asynchronously
	go func() {
		if err := uc.audit.LogAudit(context.Background(), log); err != nil {
			// Log error to monitoring system
		}
	}()
}

func (uc *userUseCaseImpl) UpdateProfile(ctx context.Context, userID string, name string, phone string) error {
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := uc.userRepo.FindByID(ctx, uint(id))
	if err != nil {
		return err
	}

	// Check if new phone number already exists
	if phone != user.Phone && phone != "" {
		if _, err := uc.userRepo.FindByPhone(ctx, phone); err == nil {
			return errors.New("phone number already exists")
		}
	}

	user.Name = name
	user.Phone = phone
	return uc.userRepo.Update(ctx, user)
}

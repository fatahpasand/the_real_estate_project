// internal/adapters/repositories/mysql/user_repository.go
package mysql

import (
	"context"
	"time"

	"iam-service/internal/core/domain"

	"gorm.io/gorm"
)

type mysqlUserRepository struct {
	db *gorm.DB
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func NewMySQLUserRepository(db *gorm.DB) *mysqlUserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mysqlUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mysqlUserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *mysqlUserRepository) UpdateVerificationStatus(ctx context.Context, userID uint, verified bool) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"verified":   verified,
			"updated_at": time.Now(),
		}).Error
}

func (r *mysqlUserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Implement other repository methods...

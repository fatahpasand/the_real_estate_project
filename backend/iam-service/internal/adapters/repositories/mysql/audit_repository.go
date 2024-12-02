// internal/adapters/repositories/mysql/audit_repository.go
package mysql

import (
	"context"
	"iam-service/internal/core/domain"

	"gorm.io/gorm"
)

type mysqlAuditRepository struct {
	db *gorm.DB
}

func NewMySQLAuditRepository(db *gorm.DB) *mysqlAuditRepository {
	return &mysqlAuditRepository{db: db}
}

func (r *mysqlAuditRepository) LogAudit(ctx context.Context, log *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

package repository

import (
	"context"

	"github.com/snisid/platform/internal/domain/authorization/entity"
	"gorm.io/gorm"
)

type PolicyRepository interface {
	GetActivePolicies(ctx context.Context) ([]entity.Policy, error)
	GetRoleGrants(ctx context.Context) ([]entity.RoleGrant, error)
}

type postgresPolicyRepo struct {
	db *gorm.DB
}

func NewPostgresPolicyRepository(db *gorm.DB) PolicyRepository {
	return &postgresPolicyRepo{db: db}
}

func (r *postgresPolicyRepo) GetActivePolicies(ctx context.Context) ([]entity.Policy, error) {
	var policies []entity.Policy
	err := r.db.WithContext(ctx).Where("enabled = ?", true).Find(&policies).Error
	return policies, err
}

func (r *postgresPolicyRepo) GetRoleGrants(ctx context.Context) ([]entity.RoleGrant, error) {
	var grants []entity.RoleGrant
	err := r.db.WithContext(ctx).Find(&grants).Error
	return grants, err
}

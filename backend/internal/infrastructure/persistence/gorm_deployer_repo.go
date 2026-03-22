package persistence

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormDeployerRepository struct {
	db *gorm.DB
}

func NewGormDeployerRepository(db *gorm.DB) *GormDeployerRepository {
	return &GormDeployerRepository{db: db}
}

func (r *GormDeployerRepository) IsDeployer(ctx context.Context, tenantID, userID uuid.UUID) (bool, error) {
	tx := GetTx(ctx, r.db)
	var count int64
	if err := tx.Model(&DeployerModel{}).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormDeployerRepository) AddDeployer(ctx context.Context, tenantID, userID uuid.UUID) error {
	tx := GetTx(ctx, r.db)
	model := &DeployerModel{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
	}
	return tx.Save(model).Error
}

func (r *GormDeployerRepository) RemoveDeployer(ctx context.Context, tenantID, userID uuid.UUID) error {
	tx := GetTx(ctx, r.db)
	return tx.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Delete(&DeployerModel{}).Error
}

func (r *GormDeployerRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]uuid.UUID, error) {
	tx := GetTx(ctx, r.db)
	var models []DeployerModel
	if err := tx.Where("tenant_id = ?", tenantID).Find(&models).Error; err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, len(models))
	for i, m := range models {
		ids[i] = m.UserID
	}
	return ids, nil
}

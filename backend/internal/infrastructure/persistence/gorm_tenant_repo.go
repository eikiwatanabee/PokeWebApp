package persistence

import (
	"context"
	"errors"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormTenantRepository struct {
	db *gorm.DB
}

func NewGormTenantRepository(db *gorm.DB) *GormTenantRepository {
	return &GormTenantRepository{db: db}
}

func (r *GormTenantRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error) {
	tx := GetTx(ctx, r.db)
	var model TenantModel
	if err := tx.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}
	return toTenantEntity(&model), nil
}

func (r *GormTenantRepository) Save(ctx context.Context, tenant *entity.Tenant) error {
	tx := GetTx(ctx, r.db)
	model := toTenantModel(tenant)
	return tx.Save(model).Error
}

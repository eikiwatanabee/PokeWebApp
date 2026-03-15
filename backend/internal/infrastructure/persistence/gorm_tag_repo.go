package persistence

import (
	"context"
	"errors"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormTagRepository struct {
	db *gorm.DB
}

func NewGormTagRepository(db *gorm.DB) *GormTagRepository {
	return &GormTagRepository{db: db}
}

func (r *GormTagRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error) {
	tx := GetTx(ctx, r.db)
	var model TagModel
	if err := tx.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return toTagEntity(&model), nil
}

func (r *GormTagRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.Tag, error) {
	tx := GetTx(ctx, r.db)
	var models []TagModel
	if err := tx.Where("tenant_id = ?", tenantID).Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.Tag, len(models))
	for i := range models {
		entities[i] = toTagEntity(&models[i])
	}
	return entities, nil
}

func (r *GormTagRepository) FindByName(ctx context.Context, tenantID uuid.UUID, name string) (*entity.Tag, error) {
	tx := GetTx(ctx, r.db)
	var model TagModel
	if err := tx.Where("tenant_id = ? AND name = ?", tenantID, name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toTagEntity(&model), nil
}

func (r *GormTagRepository) Save(ctx context.Context, tag *entity.Tag) error {
	tx := GetTx(ctx, r.db)
	model := toTagModel(tag)
	return tx.Save(model).Error
}

func (r *GormTagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx := GetTx(ctx, r.db)
	return tx.Delete(&TagModel{}, "id = ?", id).Error
}

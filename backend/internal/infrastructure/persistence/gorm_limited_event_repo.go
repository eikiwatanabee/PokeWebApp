package persistence

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormLimitedEventRepository struct {
	db *gorm.DB
}

func NewGormLimitedEventRepository(db *gorm.DB) *GormLimitedEventRepository {
	return &GormLimitedEventRepository{db: db}
}

func (r *GormLimitedEventRepository) FindActiveByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.LimitedEvent, error) {
	tx := GetTx(ctx, r.db)
	var models []LimitedEventModel
	now := time.Now()
	if err := tx.Where("tenant_id = ? AND starts_at <= ? AND ends_at > ?", tenantID, now, now).
		Order("ends_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.LimitedEvent, len(models))
	for i := range models {
		entities[i] = toLimitedEventEntity(&models[i])
	}
	return entities, nil
}

func (r *GormLimitedEventRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.LimitedEvent, error) {
	tx := GetTx(ctx, r.db)
	var models []LimitedEventModel
	if err := tx.Where("tenant_id = ?", tenantID).Order("starts_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.LimitedEvent, len(models))
	for i := range models {
		entities[i] = toLimitedEventEntity(&models[i])
	}
	return entities, nil
}

func (r *GormLimitedEventRepository) Save(ctx context.Context, event *entity.LimitedEvent) error {
	tx := GetTx(ctx, r.db)
	model := toLimitedEventModel(event)
	return tx.Save(model).Error
}

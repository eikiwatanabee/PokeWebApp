package persistence

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormTradeRepository struct {
	db *gorm.DB
}

func NewGormTradeRepository(db *gorm.DB) *GormTradeRepository {
	return &GormTradeRepository{db: db}
}

func (r *GormTradeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Trade, error) {
	tx := GetTx(ctx, r.db)
	var model TradeModel
	if err := tx.Where("id = ?", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toTradeEntity(&model), nil
}

func (r *GormTradeRepository) FindOpenByTenantID(ctx context.Context, tenantID uuid.UUID, limit int) ([]*entity.Trade, error) {
	tx := GetTx(ctx, r.db)
	var models []TradeModel
	q := tx.Where("tenant_id = ? AND status = ?", tenantID, string(entity.TradeStatusOpen)).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.Trade, len(models))
	for i := range models {
		entities[i] = toTradeEntity(&models[i])
	}
	return entities, nil
}

func (r *GormTradeRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Trade, error) {
	tx := GetTx(ctx, r.db)
	var models []TradeModel
	if err := tx.Where("offerer_id = ? OR accepter_id = ?", userID, userID).Order("created_at DESC").Limit(50).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.Trade, len(models))
	for i := range models {
		entities[i] = toTradeEntity(&models[i])
	}
	return entities, nil
}

func (r *GormTradeRepository) Save(ctx context.Context, trade *entity.Trade) error {
	tx := GetTx(ctx, r.db)
	model := toTradeModel(trade)
	return tx.Save(model).Error
}

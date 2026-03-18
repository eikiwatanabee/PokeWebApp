package persistence

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormWeeklyEventRepository struct {
	db *gorm.DB
}

func NewGormWeeklyEventRepository(db *gorm.DB) *GormWeeklyEventRepository {
	return &GormWeeklyEventRepository{db: db}
}

func (r *GormWeeklyEventRepository) FindByTenantIDAndWeek(ctx context.Context, tenantID uuid.UUID, weekStart time.Time) (*entity.WeeklyEvent, error) {
	tx := GetTx(ctx, r.db)
	var model WeeklyEventModel
	err := tx.Where("tenant_id = ? AND week_start = ?", tenantID, weekStart).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toWeeklyEventEntity(&model), nil
}

func (r *GormWeeklyEventRepository) Save(ctx context.Context, event *entity.WeeklyEvent) error {
	tx := GetTx(ctx, r.db)
	model := toWeeklyEventModel(event)
	return tx.Save(model).Error
}

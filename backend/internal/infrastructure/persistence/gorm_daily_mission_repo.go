package persistence

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormDailyMissionRepository struct {
	db *gorm.DB
}

func NewGormDailyMissionRepository(db *gorm.DB) *GormDailyMissionRepository {
	return &GormDailyMissionRepository{db: db}
}

func (r *GormDailyMissionRepository) FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*entity.DailyMission, error) {
	tx := GetTx(ctx, r.db)
	truncated := date.Truncate(24 * time.Hour)
	var models []DailyMissionModel
	if err := tx.Where("user_id = ? AND date = ?", userID, truncated).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.DailyMission, len(models))
	for i := range models {
		entities[i] = toDailyMissionEntity(&models[i])
	}
	return entities, nil
}

func (r *GormDailyMissionRepository) Save(ctx context.Context, mission *entity.DailyMission) error {
	tx := GetTx(ctx, r.db)
	model := toDailyMissionModel(mission)
	return tx.Save(model).Error
}

func (r *GormDailyMissionRepository) SaveAll(ctx context.Context, missions []*entity.DailyMission) error {
	tx := GetTx(ctx, r.db)
	models := make([]DailyMissionModel, len(missions))
	for i, m := range missions {
		models[i] = *toDailyMissionModel(m)
	}
	return tx.Create(&models).Error
}

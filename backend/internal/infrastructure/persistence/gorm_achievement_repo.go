package persistence

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormAchievementRepository struct {
	db *gorm.DB
}

func NewGormAchievementRepository(db *gorm.DB) *GormAchievementRepository {
	return &GormAchievementRepository{db: db}
}

func (r *GormAchievementRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserAchievement, error) {
	tx := GetTx(ctx, r.db)
	var models []UserAchievementModel
	if err := tx.Where("user_id = ?", userID).Order("unlocked_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.UserAchievement, len(models))
	for i := range models {
		entities[i] = toUserAchievementEntity(&models[i])
	}
	return entities, nil
}

func (r *GormAchievementRepository) HasAchievement(ctx context.Context, userID uuid.UUID, achievementType entity.AchievementType) (bool, error) {
	tx := GetTx(ctx, r.db)
	var count int64
	if err := tx.Model(&UserAchievementModel{}).Where("user_id = ? AND achievement_type = ?", userID, string(achievementType)).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormAchievementRepository) Save(ctx context.Context, achievement *entity.UserAchievement) error {
	tx := GetTx(ctx, r.db)
	model := toUserAchievementModel(achievement)
	return tx.Save(model).Error
}

package persistence

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormGitHubActivityRepository struct {
	db *gorm.DB
}

func NewGormGitHubActivityRepository(db *gorm.DB) *GormGitHubActivityRepository {
	return &GormGitHubActivityRepository{db: db}
}

func (r *GormGitHubActivityRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.GitHubActivity, error) {
	tx := GetTx(ctx, r.db)
	var models []GitHubActivityModel
	q := tx.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GitHubActivity, len(models))
	for i := range models {
		entities[i] = toGitHubActivityEntity(&models[i])
	}
	return entities, nil
}

func (r *GormGitHubActivityRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	tx := GetTx(ctx, r.db)
	var count int64
	if err := tx.Model(&GitHubActivityModel{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormGitHubActivityRepository) CountByUserIDAndType(ctx context.Context, userID uuid.UUID, eventType entity.GitHubEventType) (int64, error) {
	tx := GetTx(ctx, r.db)
	var count int64
	if err := tx.Model(&GitHubActivityModel{}).Where("user_id = ? AND event_type = ?", userID, string(eventType)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormGitHubActivityRepository) TotalXPByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	tx := GetTx(ctx, r.db)
	var total *int
	if err := tx.Model(&GitHubActivityModel{}).Where("user_id = ?", userID).Select("COALESCE(SUM(xp), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}

func (r *GormGitHubActivityRepository) Save(ctx context.Context, activity *entity.GitHubActivity) error {
	tx := GetTx(ctx, r.db)
	model := toGitHubActivityModel(activity)
	return tx.Save(model).Error
}

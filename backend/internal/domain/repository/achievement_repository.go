package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type AchievementRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserAchievement, error)
	HasAchievement(ctx context.Context, userID uuid.UUID, achievementType entity.AchievementType) (bool, error)
	Save(ctx context.Context, achievement *entity.UserAchievement) error
}

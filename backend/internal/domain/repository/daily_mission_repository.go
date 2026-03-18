package repository

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type DailyMissionRepository interface {
	FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*entity.DailyMission, error)
	Save(ctx context.Context, mission *entity.DailyMission) error
	SaveAll(ctx context.Context, missions []*entity.DailyMission) error
}

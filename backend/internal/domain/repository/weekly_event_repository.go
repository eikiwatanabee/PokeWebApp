package repository

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type WeeklyEventRepository interface {
	FindByTenantIDAndWeek(ctx context.Context, tenantID uuid.UUID, weekStart time.Time) (*entity.WeeklyEvent, error)
	Save(ctx context.Context, event *entity.WeeklyEvent) error
}

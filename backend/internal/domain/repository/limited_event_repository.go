package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type LimitedEventRepository interface {
	FindActiveByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.LimitedEvent, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.LimitedEvent, error)
	Save(ctx context.Context, event *entity.LimitedEvent) error
}

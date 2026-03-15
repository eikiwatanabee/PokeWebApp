package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type TeamRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.Team, error)
	Save(ctx context.Context, team *entity.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
}

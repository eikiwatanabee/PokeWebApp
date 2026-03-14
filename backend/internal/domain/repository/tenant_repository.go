package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type TenantRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
	Save(ctx context.Context, tenant *entity.Tenant) error
}

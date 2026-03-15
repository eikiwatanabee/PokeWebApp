package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type TagRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.Tag, error)
	FindByName(ctx context.Context, tenantID uuid.UUID, name string) (*entity.Tag, error)
	Save(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
}

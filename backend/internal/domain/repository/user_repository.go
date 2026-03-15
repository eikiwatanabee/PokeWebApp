package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}

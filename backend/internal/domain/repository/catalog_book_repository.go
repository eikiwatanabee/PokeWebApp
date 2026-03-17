package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type CatalogBookRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.CatalogBook, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID, search string, page, limit int) ([]*entity.CatalogBook, int64, error)
	Save(ctx context.Context, book *entity.CatalogBook) error
	Update(ctx context.Context, book *entity.CatalogBook) error
	Delete(ctx context.Context, id uuid.UUID) error
}

package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type BookRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Book, error)
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.Book, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, status string, tagID string, page, limit int) ([]*entity.Book, int64, error)
	Save(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}

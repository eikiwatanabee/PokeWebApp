package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type MemoRepository interface {
	FindByBookID(ctx context.Context, bookID uuid.UUID) ([]*entity.Memo, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Memo, error)
	Save(ctx context.Context, memo *entity.Memo) error
	Delete(ctx context.Context, id uuid.UUID) error
}

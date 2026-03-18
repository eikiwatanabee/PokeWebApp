package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type TradeRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Trade, error)
	FindOpenByTenantID(ctx context.Context, tenantID uuid.UUID, limit int) ([]*entity.Trade, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Trade, error)
	Save(ctx context.Context, trade *entity.Trade) error
}

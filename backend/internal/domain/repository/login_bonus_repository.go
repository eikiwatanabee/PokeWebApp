package repository

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type LoginBonusRepository interface {
	FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*entity.LoginBonus, error)
	FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.LoginBonus, error)
	Save(ctx context.Context, bonus *entity.LoginBonus) error
}

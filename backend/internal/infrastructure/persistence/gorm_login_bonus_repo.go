package persistence

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormLoginBonusRepository struct {
	db *gorm.DB
}

func NewGormLoginBonusRepository(db *gorm.DB) *GormLoginBonusRepository {
	return &GormLoginBonusRepository{db: db}
}

func (r *GormLoginBonusRepository) FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*entity.LoginBonus, error) {
	tx := GetTx(ctx, r.db)
	truncated := date.Truncate(24 * time.Hour)
	var model LoginBonusModel
	if err := tx.Where("user_id = ? AND date = ?", userID, truncated).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toLoginBonusEntity(&model), nil
}

func (r *GormLoginBonusRepository) FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.LoginBonus, error) {
	tx := GetTx(ctx, r.db)
	var model LoginBonusModel
	if err := tx.Where("user_id = ?", userID).Order("date DESC").First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toLoginBonusEntity(&model), nil
}

func (r *GormLoginBonusRepository) Save(ctx context.Context, bonus *entity.LoginBonus) error {
	tx := GetTx(ctx, r.db)
	model := toLoginBonusModel(bonus)
	return tx.Save(model).Error
}

package persistence

import (
	"context"
	"errors"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	tx := GetTx(ctx, r.db)
	var model UserModel
	if err := tx.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toUserEntity(&model), nil
}

func (r *GormUserRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	tx := GetTx(ctx, r.db)
	var model UserModel
	if err := tx.First(&model, "google_id = ?", googleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error for OAuth flow
		}
		return nil, err
	}
	return toUserEntity(&model), nil
}

func (r *GormUserRepository) FindByGitHubUsername(ctx context.Context, username string) (*entity.User, error) {
	tx := GetTx(ctx, r.db)
	var model UserModel
	if err := tx.First(&model, "git_hub_username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toUserEntity(&model), nil
}

func (r *GormUserRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error) {
	tx := GetTx(ctx, r.db)
	var models []UserModel
	if err := tx.Where("tenant_id = ?", tenantID).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.User, len(models))
	for i := range models {
		entities[i] = toUserEntity(&models[i])
	}
	return entities, nil
}

func (r *GormUserRepository) Save(ctx context.Context, user *entity.User) error {
	tx := GetTx(ctx, r.db)
	model := toUserModel(user)
	return tx.Save(model).Error
}

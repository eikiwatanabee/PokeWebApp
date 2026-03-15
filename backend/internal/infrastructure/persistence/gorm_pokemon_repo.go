package persistence

import (
	"context"
	"errors"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormPokemonRepository struct {
	db *gorm.DB
}

func NewGormPokemonRepository(db *gorm.DB) *GormPokemonRepository {
	return &GormPokemonRepository{db: db}
}

func (r *GormPokemonRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserPokemon, error) {
	tx := GetTx(ctx, r.db)
	var models []UserPokemonModel
	if err := tx.Where("user_id = ?", userID).Order("caught_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.UserPokemon, len(models))
	for i := range models {
		entities[i] = toUserPokemonEntity(&models[i])
	}
	return entities, nil
}

func (r *GormPokemonRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.UserPokemon, error) {
	tx := GetTx(ctx, r.db)
	var model UserPokemonModel
	if err := tx.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pokemon not found")
		}
		return nil, err
	}
	return toUserPokemonEntity(&model), nil
}

func (r *GormPokemonRepository) Save(ctx context.Context, pokemon *entity.UserPokemon) error {
	tx := GetTx(ctx, r.db)
	model := toUserPokemonModel(pokemon)
	return tx.Save(model).Error
}

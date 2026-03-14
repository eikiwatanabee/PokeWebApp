package repository

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type PokemonRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserPokemon, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.UserPokemon, error)
	Save(ctx context.Context, pokemon *entity.UserPokemon) error
}

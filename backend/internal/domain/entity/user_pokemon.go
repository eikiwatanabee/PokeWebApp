package entity

import (
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

type UserPokemon struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	BookID   uuid.UUID
	Pokemon  valueobject.PokemonInfo
	CaughtAt time.Time
}

func NewUserPokemon(userID, bookID uuid.UUID, pokemon valueobject.PokemonInfo) *UserPokemon {
	return &UserPokemon{
		ID:       uuid.New(),
		UserID:   userID,
		BookID:   bookID,
		Pokemon:  pokemon,
		CaughtAt: time.Now(),
	}
}

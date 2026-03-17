package entity

import (
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

type UserPokemon struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	BookID     uuid.UUID // kept for DB compatibility, now stores ActivityID
	ActivityID uuid.UUID
	Pokemon    valueobject.PokemonInfo
	CaughtAt   time.Time
}

func NewUserPokemon(userID, bookID uuid.UUID, pokemon valueobject.PokemonInfo) *UserPokemon {
	return &UserPokemon{
		ID:         uuid.New(),
		UserID:     userID,
		BookID:     bookID,
		ActivityID: bookID,
		Pokemon:    pokemon,
		CaughtAt:   time.Now(),
	}
}

func NewUserPokemonFromActivity(userID, activityID uuid.UUID, pokemon valueobject.PokemonInfo) *UserPokemon {
	return &UserPokemon{
		ID:         uuid.New(),
		UserID:     userID,
		BookID:     activityID,
		ActivityID: activityID,
		Pokemon:    pokemon,
		CaughtAt:   time.Now(),
	}
}

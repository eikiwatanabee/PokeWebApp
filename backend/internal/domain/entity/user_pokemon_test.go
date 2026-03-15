package entity

import (
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

func TestNewUserPokemon(t *testing.T) {
	userID := uuid.New()
	bookID := uuid.New()
	pokemon := valueobject.PokemonInfo{
		PokemonID: 25,
		Name:      "pikachu",
		SpriteURL: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png",
		Types:     []string{"electric"},
	}

	up := NewUserPokemon(userID, bookID, pokemon)

	if up.ID == uuid.Nil {
		t.Error("NewUserPokemon() ID should not be nil")
	}
	if up.UserID != userID {
		t.Errorf("NewUserPokemon() UserID = %v, want %v", up.UserID, userID)
	}
	if up.BookID != bookID {
		t.Errorf("NewUserPokemon() BookID = %v, want %v", up.BookID, bookID)
	}
	if up.Pokemon.PokemonID != 25 {
		t.Errorf("NewUserPokemon() PokemonID = %v, want 25", up.Pokemon.PokemonID)
	}
	if up.Pokemon.Name != "pikachu" {
		t.Errorf("NewUserPokemon() Name = %v, want pikachu", up.Pokemon.Name)
	}
	if up.CaughtAt.IsZero() {
		t.Error("NewUserPokemon() CaughtAt should not be zero")
	}
}

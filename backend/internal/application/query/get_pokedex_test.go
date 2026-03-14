package query

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestGetPokedexHandler_Handle(t *testing.T) {
	t.Run("returns pokemon for user", func(t *testing.T) {
		pokemonRepo := newMockPokemonRepo()
		userID := uuid.New()
		bookID := uuid.New()
		up := entity.NewUserPokemon(userID, bookID, newTestPokemonInfo())
		pokemonRepo.pokemons[up.ID] = up

		handler := NewGetPokedexHandler(pokemonRepo)
		result, err := handler.Handle(context.Background(), &GetPokedexQuery{UserID: userID.String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.Total != 1 {
			t.Errorf("Total = %v, want 1", result.Total)
		}
		if result.Pokemon[0].PokemonName != "pikachu" {
			t.Errorf("PokemonName = %v, want pikachu", result.Pokemon[0].PokemonName)
		}
	})

	t.Run("empty pokedex", func(t *testing.T) {
		handler := NewGetPokedexHandler(newMockPokemonRepo())
		result, err := handler.Handle(context.Background(), &GetPokedexQuery{UserID: uuid.New().String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.Total != 0 {
			t.Errorf("Total = %v, want 0", result.Total)
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		handler := NewGetPokedexHandler(newMockPokemonRepo())
		_, err := handler.Handle(context.Background(), &GetPokedexQuery{UserID: "invalid"})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

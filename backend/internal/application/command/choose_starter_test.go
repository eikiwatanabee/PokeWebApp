package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/service"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

func TestChooseStarterHandler_Handle(t *testing.T) {
	bulbasaur := &valueobject.PokemonInfo{
		PokemonID: 1,
		Name:      "bulbasaur",
		SpriteURL: "https://example.com/bulbasaur.png",
		Types:     []string{"grass", "poison"},
	}

	t.Run("success - choose bulbasaur", func(t *testing.T) {
		pokemonRepo := newMockPokemonRepo()
		fetcher := &mockPokemonFetcher{pokemon: bulbasaur}
		gachaSvc := service.NewPokemonGachaService(fetcher)

		handler := NewChooseStarterHandler(&mockUoW{}, pokemonRepo, gachaSvc)
		result, err := handler.Handle(context.Background(), &ChooseStarterCommand{
			UserID:    uuid.New().String(),
			PokemonID: 1,
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.PokemonName != "bulbasaur" {
			t.Errorf("PokemonName = %v, want bulbasaur", result.PokemonName)
		}
		if len(pokemonRepo.pokemons) != 1 {
			t.Errorf("expected 1 pokemon saved, got %d", len(pokemonRepo.pokemons))
		}
	})

	t.Run("invalid starter ID", func(t *testing.T) {
		handler := NewChooseStarterHandler(&mockUoW{}, newMockPokemonRepo(), service.NewPokemonGachaService(&mockPokemonFetcher{}))
		_, err := handler.Handle(context.Background(), &ChooseStarterCommand{
			UserID:    uuid.New().String(),
			PokemonID: 25, // pikachu - not a starter
		})
		if err != ErrInvalidStarter {
			t.Errorf("Handle() error = %v, want %v", err, ErrInvalidStarter)
		}
	})

	t.Run("already has starter", func(t *testing.T) {
		pokemonRepo := newMockPokemonRepo()
		userID := uuid.New()

		// Add existing pokemon
		existing := &valueobject.PokemonInfo{PokemonID: 4, Name: "charmander"}
		up := newTestUserPokemon(userID, *existing)
		pokemonRepo.pokemons[up.ID] = up

		handler := NewChooseStarterHandler(&mockUoW{}, pokemonRepo, service.NewPokemonGachaService(&mockPokemonFetcher{pokemon: bulbasaur}))
		_, err := handler.Handle(context.Background(), &ChooseStarterCommand{
			UserID:    userID.String(),
			PokemonID: 1,
		})
		if err != ErrAlreadyHasStarter {
			t.Errorf("Handle() error = %v, want %v", err, ErrAlreadyHasStarter)
		}
	})
}

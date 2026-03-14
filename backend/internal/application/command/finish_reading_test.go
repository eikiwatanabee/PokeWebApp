package command

import (
	"context"
	"errors"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/service"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

func TestFinishReadingHandler_Handle(t *testing.T) {
	pikachu := &valueobject.PokemonInfo{
		PokemonID: 25,
		Name:      "pikachu",
		SpriteURL: "https://example.com/pikachu.png",
		Types:     []string{"electric"},
	}

	t.Run("success - finish reading and get pokemon", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		pokemonRepo := newMockPokemonRepo()
		fetcher := &mockPokemonFetcher{pokemon: pikachu}
		gachaSvc := service.NewPokemonGachaService(fetcher)

		userID := uuid.New()
		book, _ := entity.NewBook(userID, "DDD", "Eric Evans")
		book.Status = valueobject.Reading
		bookRepo.books[book.ID] = book

		handler := NewFinishReadingHandler(&mockUoW{}, bookRepo, pokemonRepo, gachaSvc)
		result, err := handler.Handle(context.Background(), &FinishReadingCommand{
			BookID: book.ID.String(),
			UserID: userID.String(),
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.PokemonID != 25 {
			t.Errorf("PokemonID = %v, want 25", result.PokemonID)
		}
		if result.PokemonName != "pikachu" {
			t.Errorf("PokemonName = %v, want pikachu", result.PokemonName)
		}
		if bookRepo.books[book.ID].Status != valueobject.Finished {
			t.Error("book should be finished")
		}
		if len(pokemonRepo.pokemons) != 1 {
			t.Errorf("expected 1 pokemon saved, got %d", len(pokemonRepo.pokemons))
		}
	})

	t.Run("book not in reading status", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		fetcher := &mockPokemonFetcher{pokemon: pikachu}
		gachaSvc := service.NewPokemonGachaService(fetcher)

		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		bookRepo.books[book.ID] = book

		handler := NewFinishReadingHandler(&mockUoW{}, bookRepo, newMockPokemonRepo(), gachaSvc)
		_, err := handler.Handle(context.Background(), &FinishReadingCommand{
			BookID: book.ID.String(),
			UserID: uuid.New().String(),
		})
		if err == nil {
			t.Error("Handle() should return error for unread book")
		}
	})

	t.Run("gacha service error", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		fetchErr := errors.New("api error")
		fetcher := &mockPokemonFetcher{err: fetchErr}
		gachaSvc := service.NewPokemonGachaService(fetcher)

		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		book.Status = valueobject.Reading
		bookRepo.books[book.ID] = book

		handler := NewFinishReadingHandler(&mockUoW{}, bookRepo, newMockPokemonRepo(), gachaSvc)
		_, err := handler.Handle(context.Background(), &FinishReadingCommand{
			BookID: book.ID.String(),
			UserID: uuid.New().String(),
		})
		if !errors.Is(err, fetchErr) {
			t.Errorf("Handle() error = %v, want %v", err, fetchErr)
		}
	})
}

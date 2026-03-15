package service

import (
	"context"
	"errors"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
)

type mockFetcher struct {
	pokemon *valueobject.PokemonInfo
	err     error
}

func (m *mockFetcher) FetchRandom(ctx context.Context) (*valueobject.PokemonInfo, error) {
	return m.pokemon, m.err
}

func (m *mockFetcher) FetchByID(ctx context.Context, id int) (*valueobject.PokemonInfo, error) {
	return m.pokemon, m.err
}

func TestPokemonGachaService_Draw(t *testing.T) {
	pikachu := &valueobject.PokemonInfo{
		PokemonID: 25,
		Name:      "pikachu",
		SpriteURL: "https://example.com/pikachu.png",
		Types:     []string{"electric"},
	}

	t.Run("success", func(t *testing.T) {
		svc := NewPokemonGachaService(&mockFetcher{pokemon: pikachu})
		got, err := svc.Draw(context.Background())
		if err != nil {
			t.Fatalf("Draw() error = %v", err)
		}
		if got.PokemonID != 25 {
			t.Errorf("Draw() PokemonID = %v, want 25", got.PokemonID)
		}
		if got.Name != "pikachu" {
			t.Errorf("Draw() Name = %v, want pikachu", got.Name)
		}
	})

	t.Run("fetch error", func(t *testing.T) {
		fetchErr := errors.New("api error")
		svc := NewPokemonGachaService(&mockFetcher{err: fetchErr})
		_, err := svc.Draw(context.Background())
		if !errors.Is(err, fetchErr) {
			t.Errorf("Draw() error = %v, want %v", err, fetchErr)
		}
	})
}

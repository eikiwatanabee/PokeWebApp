package service

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
)

// PokemonFetcher is the interface for fetching Pokemon data from an external source.
// Implemented in infrastructure layer (PokeAPI client).
type PokemonFetcher interface {
	FetchRandom(ctx context.Context) (*valueobject.PokemonInfo, error)
	FetchByID(ctx context.Context, id int) (*valueobject.PokemonInfo, error)
}

// PokemonGachaService handles the random Pokemon assignment logic.
type PokemonGachaService struct {
	fetcher PokemonFetcher
}

func NewPokemonGachaService(fetcher PokemonFetcher) *PokemonGachaService {
	return &PokemonGachaService{fetcher: fetcher}
}

// Draw fetches a random Pokemon from the external API.
func (s *PokemonGachaService) Draw(ctx context.Context) (*valueobject.PokemonInfo, error) {
	return s.fetcher.FetchRandom(ctx)
}

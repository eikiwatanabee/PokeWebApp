package service

import (
	"context"
	"math/rand"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
)

// PokemonFetcher is the interface for fetching Pokemon data from an external source.
type PokemonFetcher interface {
	FetchRandom(ctx context.Context) (*valueobject.PokemonInfo, error)
	FetchByID(ctx context.Context, id int) (*valueobject.PokemonInfo, error)
}

// PokemonGachaService handles the random Pokemon assignment with rarity.
type PokemonGachaService struct {
	fetcher PokemonFetcher
}

func NewPokemonGachaService(fetcher PokemonFetcher) *PokemonGachaService {
	return &PokemonGachaService{fetcher: fetcher}
}

// Draw fetches a random Pokemon. Rarity is determined by PokeAPI species data.
func (s *PokemonGachaService) Draw(ctx context.Context) (*valueobject.PokemonInfo, error) {
	return s.DrawWithBoost(ctx, 1.0)
}

// DrawWithBoost fetches a random Pokemon with streak-based rarity upgrade chance.
// The streakMultiplier increases the chance of upgrading the Pokemon's natural rarity.
func (s *PokemonGachaService) DrawWithBoost(ctx context.Context, streakMultiplier float64) (*valueobject.PokemonInfo, error) {
	info, err := s.fetcher.FetchRandom(ctx)
	if err != nil {
		return nil, err
	}

	// Apply streak boost: chance to upgrade rarity tier
	if streakMultiplier > 1.0 {
		info.Rarity = applyStreakBoost(info.Rarity, streakMultiplier)
	}

	return info, nil
}

// FetchByID fetches a specific Pokemon by its ID.
func (s *PokemonGachaService) FetchByID(ctx context.Context, id int) (*valueobject.PokemonInfo, error) {
	return s.fetcher.FetchByID(ctx, id)
}

// applyStreakBoost gives a chance to upgrade rarity based on streak multiplier.
// Higher multiplier = higher chance to upgrade.
func applyStreakBoost(rarity valueobject.PokemonRarity, multiplier float64) valueobject.PokemonRarity {
	// Boost chance: x1.5→25%, x2.0→50%, x3.0→100%
	boostChance := (multiplier - 1.0) * 0.5
	if rand.Float64() < boostChance {
		upgraded := upgradeRarity(rarity)
		// At x3.0, try a second upgrade
		if multiplier >= 3.0 && rand.Float64() < 0.25 {
			upgraded = upgradeRarity(upgraded)
		}
		return upgraded
	}
	return rarity
}

func upgradeRarity(rarity valueobject.PokemonRarity) valueobject.PokemonRarity {
	switch rarity {
	case valueobject.RarityCommon:
		return valueobject.RarityUncommon
	case valueobject.RarityUncommon:
		return valueobject.RarityRare
	case valueobject.RarityRare:
		return valueobject.RarityEpic
	case valueobject.RarityEpic:
		return valueobject.RarityLegendary
	default:
		return rarity // Already legendary, can't upgrade
	}
}

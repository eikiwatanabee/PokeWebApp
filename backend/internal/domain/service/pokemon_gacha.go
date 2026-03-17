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

// Rarity tier Pokemon ID ranges
var (
	// Legendary: mythical/legendary Pokemon
	legendaryIDs = []int{
		150, 151, // Mewtwo, Mew
		249, 250, 251, // Lugia, Ho-Oh, Celebi
		382, 383, 384, 385, 386, // Kyogre, Groudon, Rayquaza, Jirachi, Deoxys
		483, 484, 487, 491, 492, 493, // Dialga, Palkia, Giratina, Darkrai, Shaymin, Arceus
	}
	// Epic: legendary birds/beasts, sub-legendaries
	epicIDs = []int{
		144, 145, 146, // Articuno, Zapdos, Moltres
		243, 244, 245, // Raikou, Entei, Suicune
		377, 378, 379, 380, 381, // Regis, Latias, Latios
		480, 481, 482, 485, 486, // Lake trio, Heatran, Regigigas
		638, 639, 640, 641, 642, 643, 644, 645, 646, // Swords of Justice, Forces of Nature, Tao trio
	}
)

// PokemonGachaService handles the random Pokemon assignment with rarity.
type PokemonGachaService struct {
	fetcher PokemonFetcher
}

func NewPokemonGachaService(fetcher PokemonFetcher) *PokemonGachaService {
	return &PokemonGachaService{fetcher: fetcher}
}

// Draw fetches a random Pokemon with rarity-weighted selection.
// streakMultiplier boosts rare chances (1.0 = normal, 1.5+ = boosted).
func (s *PokemonGachaService) Draw(ctx context.Context) (*valueobject.PokemonInfo, error) {
	return s.DrawWithBoost(ctx, 1.0)
}

// DrawWithBoost fetches a Pokemon with rarity boosted by streak multiplier.
func (s *PokemonGachaService) DrawWithBoost(ctx context.Context, streakMultiplier float64) (*valueobject.PokemonInfo, error) {
	rarity := rollRarity(streakMultiplier)
	pokemonID := pickPokemonIDForRarity(rarity)

	info, err := s.fetcher.FetchByID(ctx, pokemonID)
	if err != nil {
		return nil, err
	}
	info.Rarity = rarity
	return info, nil
}

// FetchByID fetches a specific Pokemon by its ID.
func (s *PokemonGachaService) FetchByID(ctx context.Context, id int) (*valueobject.PokemonInfo, error) {
	info, err := s.fetcher.FetchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	info.Rarity = determineRarity(id)
	return info, nil
}

func rollRarity(boost float64) valueobject.PokemonRarity {
	roll := rand.Float64() * 100

	// Base rates: Common 60%, Uncommon 25%, Rare 12%, Epic 2.5%, Legendary 0.5%
	// Boost reduces common chance, increases rare+
	legendaryThreshold := 0.5 * boost
	epicThreshold := legendaryThreshold + 2.5*boost
	rareThreshold := epicThreshold + 12.0*boost
	uncommonThreshold := rareThreshold + 25.0

	switch {
	case roll < legendaryThreshold:
		return valueobject.RarityLegendary
	case roll < epicThreshold:
		return valueobject.RarityEpic
	case roll < rareThreshold:
		return valueobject.RarityRare
	case roll < uncommonThreshold:
		return valueobject.RarityUncommon
	default:
		return valueobject.RarityCommon
	}
}

func pickPokemonIDForRarity(rarity valueobject.PokemonRarity) int {
	switch rarity {
	case valueobject.RarityLegendary:
		return legendaryIDs[rand.Intn(len(legendaryIDs))]
	case valueobject.RarityEpic:
		return epicIDs[rand.Intn(len(epicIDs))]
	case valueobject.RarityRare:
		// Final evolution forms (high ID pokemon within each gen)
		// Pick from ranges that tend to be final evolutions
		rareRanges := []int{3, 6, 9, 12, 15, 18, 31, 34, 38, 45, 47, 49, 51, 55, 59, 62, 65, 68, 71, 76, 78, 80, 82, 85, 87, 89, 91, 94, 97, 99, 101, 103, 105, 106, 107, 108, 110, 112, 113, 114, 115, 119, 121, 122, 123, 124, 125, 126, 127, 128, 130, 131, 132, 134, 135, 136, 137, 139, 141, 142, 143}
		return rareRanges[rand.Intn(len(rareRanges))]
	case valueobject.RarityUncommon:
		// Mid-evolution pokemon (2nd stage)
		uncommonRanges := []int{2, 5, 8, 11, 14, 17, 20, 22, 24, 25, 26, 28, 30, 33, 36, 40, 42, 44, 49, 57, 61, 64, 67, 70, 73, 75, 85, 93, 95, 110, 117, 119, 148, 153, 156, 159, 162, 164, 166, 168, 171, 176, 178, 180, 184, 186, 189, 195, 199}
		return uncommonRanges[rand.Intn(len(uncommonRanges))]
	default:
		// Common: basic Pokemon from Gen 1-8 (first stage)
		commonIDs := []int{1, 4, 7, 10, 13, 16, 19, 21, 23, 27, 29, 32, 35, 37, 39, 41, 43, 46, 48, 50, 52, 54, 56, 58, 60, 63, 66, 69, 72, 74, 77, 79, 81, 83, 84, 86, 88, 90, 92, 96, 98, 100, 102, 104, 109, 111, 116, 118, 120, 129, 133, 138, 140, 147, 152, 155, 158, 161, 163, 165, 167}
		return commonIDs[rand.Intn(len(commonIDs))]
	}
}

func determineRarity(pokemonID int) valueobject.PokemonRarity {
	for _, id := range legendaryIDs {
		if id == pokemonID {
			return valueobject.RarityLegendary
		}
	}
	for _, id := range epicIDs {
		if id == pokemonID {
			return valueobject.RarityEpic
		}
	}
	// Simplified: higher evolution = rarer
	return valueobject.RarityCommon
}

package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetTrainerCardQuery struct {
	UserID string
}

type TrainerCardPokemonDTO struct {
	PokemonName string `json:"pokemon_name"`
	SpriteURL   string `json:"sprite_url"`
	Rarity      string `json:"rarity"`
}

type TrainerCardAchievementDTO struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type GetTrainerCardResult struct {
	UserID         string                      `json:"user_id"`
	Name           string                      `json:"name"`
	AvatarURL      string                      `json:"avatar_url"`
	GitHubUsername string                      `json:"github_username"`
	Level          int                         `json:"level"`
	TotalXP        int                         `json:"total_xp"`
	CurrentStreak  int                         `json:"current_streak"`
	MaxStreak      int                         `json:"max_streak"`
	PokemonCount   int                         `json:"pokemon_count"`
	FeaturedPokemon []TrainerCardPokemonDTO    `json:"featured_pokemon"`
	Achievements    []TrainerCardAchievementDTO `json:"achievements"`
	AchievementCount int                       `json:"achievement_count"`
}

type GetTrainerCardHandler struct {
	userRepo        repository.UserRepository
	pokemonRepo     repository.PokemonRepository
	achievementRepo repository.AchievementRepository
}

func NewGetTrainerCardHandler(
	userRepo repository.UserRepository,
	pokemonRepo repository.PokemonRepository,
	achievementRepo repository.AchievementRepository,
) *GetTrainerCardHandler {
	return &GetTrainerCardHandler{
		userRepo:        userRepo,
		pokemonRepo:     pokemonRepo,
		achievementRepo: achievementRepo,
	}
}

func (h *GetTrainerCardHandler) Handle(ctx context.Context, q *GetTrainerCardQuery) (*GetTrainerCardResult, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	pokemon, err := h.pokemonRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	achievements, err := h.achievementRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Featured pokemon: up to 6, prioritizing rare/epic/legendary
	featured := selectFeaturedPokemon(pokemon, 6)
	featuredDTOs := make([]TrainerCardPokemonDTO, len(featured))
	for i, p := range featured {
		featuredDTOs[i] = TrainerCardPokemonDTO{
			PokemonName: p.Pokemon.Name,
			SpriteURL:   p.Pokemon.SpriteURL,
			Rarity:      string(p.Pokemon.Rarity),
		}
	}

	// Achievement definitions lookup
	defMap := make(map[entity.AchievementType]entity.AchievementDefinition)
	for _, d := range entity.AchievementDefinitions {
		defMap[d.Type] = d
	}

	achievementDTOs := make([]TrainerCardAchievementDTO, 0, len(achievements))
	for _, a := range achievements {
		if def, ok := defMap[a.AchievementType]; ok {
			achievementDTOs = append(achievementDTOs, TrainerCardAchievementDTO{
				Name: def.Name,
				Icon: def.Icon,
			})
		}
	}

	return &GetTrainerCardResult{
		UserID:           user.ID.String(),
		Name:             user.Name,
		AvatarURL:        user.AvatarURL,
		GitHubUsername:   user.GitHubUsername,
		Level:            user.Level,
		TotalXP:          user.TotalXP,
		CurrentStreak:    user.CurrentStreak,
		MaxStreak:        user.MaxStreak,
		PokemonCount:     len(pokemon),
		FeaturedPokemon:  featuredDTOs,
		Achievements:     achievementDTOs,
		AchievementCount: len(achievements),
	}, nil
}

func selectFeaturedPokemon(pokemon []*entity.UserPokemon, limit int) []*entity.UserPokemon {
	if len(pokemon) <= limit {
		return pokemon
	}

	rarityOrder := map[string]int{
		"legendary": 0,
		"epic":      1,
		"rare":      2,
		"uncommon":  3,
		"common":    4,
	}

	// Sort by rarity (best first), then by caught_at (newest first)
	sorted := make([]*entity.UserPokemon, len(pokemon))
	copy(sorted, pokemon)

	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			ri := rarityOrder[string(sorted[i].Pokemon.Rarity)]
			rj := rarityOrder[string(sorted[j].Pokemon.Rarity)]
			if rj < ri || (rj == ri && sorted[j].CaughtAt.After(sorted[i].CaughtAt)) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted[:limit]
}

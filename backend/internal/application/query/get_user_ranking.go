package query

import (
	"context"
	"sort"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetUserRankingQuery struct {
	TenantID string
	SortBy   string // "xp", "level", "pokemon", "streak"
}

type UserRankDTO struct {
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
	GitHubUsername string `json:"github_username"`
	Level          int    `json:"level"`
	TotalXP        int    `json:"total_xp"`
	PokemonCount   int    `json:"pokemon_count"`
	CurrentStreak  int    `json:"current_streak"`
	MaxStreak      int    `json:"max_streak"`
	Rank           int    `json:"rank"`
}

type GetUserRankingResult struct {
	Rankings []UserRankDTO `json:"rankings"`
	SortBy   string        `json:"sort_by"`
}

type GetUserRankingHandler struct {
	userRepo    repository.UserRepository
	pokemonRepo repository.PokemonRepository
}

func NewGetUserRankingHandler(
	userRepo repository.UserRepository,
	pokemonRepo repository.PokemonRepository,
) *GetUserRankingHandler {
	return &GetUserRankingHandler{
		userRepo:    userRepo,
		pokemonRepo: pokemonRepo,
	}
}

func (h *GetUserRankingHandler) Handle(ctx context.Context, q *GetUserRankingQuery) (*GetUserRankingResult, error) {
	tenantID, err := uuid.Parse(q.TenantID)
	if err != nil {
		return nil, err
	}

	users, err := h.userRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	rankings := make([]UserRankDTO, 0, len(users))
	for _, u := range users {
		pokemon, err := h.pokemonRepo.FindByUserID(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		rankings = append(rankings, UserRankDTO{
			UserID:         u.ID.String(),
			Name:           u.Name,
			AvatarURL:      u.AvatarURL,
			GitHubUsername: u.GitHubUsername,
			Level:          u.Level,
			TotalXP:        u.TotalXP,
			PokemonCount:   len(pokemon),
			CurrentStreak:  u.CurrentStreak,
			MaxStreak:      u.MaxStreak,
		})
	}

	sortBy := q.SortBy
	if sortBy == "" {
		sortBy = "xp"
	}

	sort.Slice(rankings, func(i, j int) bool {
		switch sortBy {
		case "level":
			if rankings[i].Level == rankings[j].Level {
				return rankings[i].TotalXP > rankings[j].TotalXP
			}
			return rankings[i].Level > rankings[j].Level
		case "pokemon":
			return rankings[i].PokemonCount > rankings[j].PokemonCount
		case "streak":
			if rankings[i].CurrentStreak == rankings[j].CurrentStreak {
				return rankings[i].MaxStreak > rankings[j].MaxStreak
			}
			return rankings[i].CurrentStreak > rankings[j].CurrentStreak
		default: // "xp"
			return rankings[i].TotalXP > rankings[j].TotalXP
		}
	})

	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	return &GetUserRankingResult{
		Rankings: rankings,
		SortBy:   sortBy,
	}, nil
}

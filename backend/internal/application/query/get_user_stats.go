package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetUserStatsQuery struct {
	UserID string
}

type UserStatsResult struct {
	TotalXP       int `json:"total_xp"`
	Level         int `json:"level"`
	XPToNextLevel int `json:"xp_to_next_level"`
	TotalCommits  int `json:"total_commits"`
	TotalMerges   int `json:"total_merges"`
	PokemonCount  int `json:"pokemon_count"`
}

type GetUserStatsHandler struct {
	userRepo     repository.UserRepository
	activityRepo repository.GitHubActivityRepository
	pokemonRepo  repository.PokemonRepository
}

func NewGetUserStatsHandler(
	userRepo repository.UserRepository,
	activityRepo repository.GitHubActivityRepository,
	pokemonRepo repository.PokemonRepository,
) *GetUserStatsHandler {
	return &GetUserStatsHandler{
		userRepo:     userRepo,
		activityRepo: activityRepo,
		pokemonRepo:  pokemonRepo,
	}
}

func (h *GetUserStatsHandler) Handle(ctx context.Context, q *GetUserStatsQuery) (*UserStatsResult, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalActivities, err := h.activityRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	pokemon, err := h.pokemonRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserStatsResult{
		TotalXP:       user.TotalXP,
		Level:         user.Level,
		XPToNextLevel: user.XPToNextLevel(),
		TotalCommits:  int(totalActivities),
		TotalMerges:   len(pokemon), // each merge = 1 pokemon
		PokemonCount:  len(pokemon),
	}, nil
}

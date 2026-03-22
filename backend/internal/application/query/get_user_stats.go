package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetUserStatsQuery struct {
	UserID string
}

type AchievementDTO struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	UnlockedAt  string `json:"unlocked_at"`
}

type UserStatsResult struct {
	TotalXP         int              `json:"total_xp"`
	Level           int              `json:"level"`
	XPToNextLevel   int              `json:"xp_to_next_level"`
	TotalCommits    int              `json:"total_commits"`
	TotalMerges     int              `json:"total_merges"`
	PokemonCount    int              `json:"pokemon_count"`
	CurrentStreak   int              `json:"current_streak"`
	MaxStreak       int              `json:"max_streak"`
	StreakMultiplier float64          `json:"streak_multiplier"`
	Achievements    []AchievementDTO `json:"achievements"`
}

type GetUserStatsHandler struct {
	userRepo        repository.UserRepository
	activityRepo    repository.GitHubActivityRepository
	pokemonRepo     repository.PokemonRepository
	achievementRepo repository.AchievementRepository
}

func NewGetUserStatsHandler(
	userRepo repository.UserRepository,
	activityRepo repository.GitHubActivityRepository,
	pokemonRepo repository.PokemonRepository,
	achievementRepo repository.AchievementRepository,
) *GetUserStatsHandler {
	return &GetUserStatsHandler{
		userRepo:        userRepo,
		activityRepo:    activityRepo,
		pokemonRepo:     pokemonRepo,
		achievementRepo: achievementRepo,
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

	achievements, err := h.achievementRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build achievement definitions map
	defMap := make(map[entity.AchievementType]entity.AchievementDefinition)
	for _, d := range entity.AchievementDefinitions {
		defMap[d.Type] = d
	}

	achievementDTOs := make([]AchievementDTO, 0, len(achievements))
	for _, a := range achievements {
		if def, ok := defMap[a.AchievementType]; ok {
			achievementDTOs = append(achievementDTOs, AchievementDTO{
				Type:        string(a.AchievementType),
				Name:        def.Name,
				Description: def.Description,
				Icon:        def.Icon,
				Category:    def.Category,
				UnlockedAt:  a.UnlockedAt.Format("2006-01-02T15:04:05Z"),
			})
		}
	}

	return &UserStatsResult{
		TotalXP:         user.TotalXP,
		Level:           user.Level,
		XPToNextLevel:   user.XPToNextLevel(),
		TotalCommits:    int(totalActivities),
		TotalMerges:     len(pokemon),
		PokemonCount:    len(pokemon),
		CurrentStreak:   user.CurrentStreak,
		MaxStreak:       user.MaxStreak,
		StreakMultiplier: user.StreakMultiplier(),
		Achievements:    achievementDTOs,
	}, nil
}

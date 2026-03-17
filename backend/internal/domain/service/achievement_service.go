package service

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

type AchievementService struct {
	achievementRepo repository.AchievementRepository
	activityRepo    repository.GitHubActivityRepository
	pokemonRepo     repository.PokemonRepository
}

func NewAchievementService(
	achievementRepo repository.AchievementRepository,
	activityRepo repository.GitHubActivityRepository,
	pokemonRepo repository.PokemonRepository,
) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		activityRepo:    activityRepo,
		pokemonRepo:     pokemonRepo,
	}
}

// CheckAndUnlock checks all achievement conditions and unlocks any newly earned ones.
func (s *AchievementService) CheckAndUnlock(ctx context.Context, user *entity.User, caughtRarity valueobject.PokemonRarity) ([]*entity.UserAchievement, error) {
	var newAchievements []*entity.UserAchievement

	// Lazy-load counters only when needed
	var (
		commitCount  *int64
		mergeCount   *int64
		reviewCount  *int64
		issueCount   *int64
		pokemonCount *int
	)

	getCommitCount := func() (int64, error) {
		if commitCount == nil {
			c, err := s.activityRepo.CountByUserID(ctx, user.ID)
			if err != nil {
				return 0, err
			}
			commitCount = &c
		}
		return *commitCount, nil
	}
	getMergeCount := func() (int64, error) {
		if mergeCount == nil {
			c, err := s.activityRepo.CountByUserIDAndType(ctx, user.ID, entity.EventPRMerge)
			if err != nil {
				return 0, err
			}
			mergeCount = &c
		}
		return *mergeCount, nil
	}
	getReviewCount := func() (int64, error) {
		if reviewCount == nil {
			c, err := s.activityRepo.CountByUserIDAndType(ctx, user.ID, entity.EventReview)
			if err != nil {
				return 0, err
			}
			reviewCount = &c
		}
		return *reviewCount, nil
	}
	getIssueCount := func() (int64, error) {
		if issueCount == nil {
			c, err := s.activityRepo.CountByUserIDAndType(ctx, user.ID, entity.EventIssueClose)
			if err != nil {
				return 0, err
			}
			issueCount = &c
		}
		return *issueCount, nil
	}
	getPokemonCount := func() (int, error) {
		if pokemonCount == nil {
			p, err := s.pokemonRepo.FindByUserID(ctx, user.ID)
			if err != nil {
				return 0, err
			}
			c := len(p)
			pokemonCount = &c
		}
		return *pokemonCount, nil
	}

	checks := []struct {
		achievementType entity.AchievementType
		condition       func() (bool, error)
	}{
		// --- コミット系 ---
		{entity.AchievementFirstCommit, func() (bool, error) {
			c, err := getCommitCount()
			return c >= 1, err
		}},
		{entity.AchievementCommit50, func() (bool, error) {
			c, err := getCommitCount()
			return c >= 50, err
		}},
		{entity.AchievementCommit100, func() (bool, error) {
			c, err := getCommitCount()
			return c >= 100, err
		}},
		{entity.AchievementCommit500, func() (bool, error) {
			c, err := getCommitCount()
			return c >= 500, err
		}},
		{entity.AchievementCommit1000, func() (bool, error) {
			c, err := getCommitCount()
			return c >= 1000, err
		}},

		// --- ストリーク系 ---
		{entity.AchievementStreak3, func() (bool, error) {
			return user.CurrentStreak >= 3, nil
		}},
		{entity.AchievementStreak7, func() (bool, error) {
			return user.CurrentStreak >= 7, nil
		}},
		{entity.AchievementStreak14, func() (bool, error) {
			return user.CurrentStreak >= 14, nil
		}},
		{entity.AchievementStreak30, func() (bool, error) {
			return user.CurrentStreak >= 30, nil
		}},
		{entity.AchievementStreak60, func() (bool, error) {
			return user.CurrentStreak >= 60, nil
		}},
		{entity.AchievementStreak100, func() (bool, error) {
			return user.CurrentStreak >= 100, nil
		}},
		{entity.AchievementStreak365, func() (bool, error) {
			return user.CurrentStreak >= 365, nil
		}},

		// --- レベル系 ---
		{entity.AchievementLevel5, func() (bool, error) {
			return user.Level >= 5, nil
		}},
		{entity.AchievementLevel10, func() (bool, error) {
			return user.Level >= 10, nil
		}},
		{entity.AchievementLevel25, func() (bool, error) {
			return user.Level >= 25, nil
		}},
		{entity.AchievementLevel50, func() (bool, error) {
			return user.Level >= 50, nil
		}},
		{entity.AchievementLevel100, func() (bool, error) {
			return user.Level >= 100, nil
		}},

		// --- XP系 ---
		{entity.AchievementXP1000, func() (bool, error) {
			return user.TotalXP >= 1000, nil
		}},
		{entity.AchievementXP5000, func() (bool, error) {
			return user.TotalXP >= 5000, nil
		}},
		{entity.AchievementXP10000, func() (bool, error) {
			return user.TotalXP >= 10000, nil
		}},
		{entity.AchievementXP50000, func() (bool, error) {
			return user.TotalXP >= 50000, nil
		}},

		// --- ポケモン系 ---
		{entity.AchievementPokemon1, func() (bool, error) {
			c, err := getPokemonCount()
			return c >= 1, err
		}},
		{entity.AchievementPokemon10, func() (bool, error) {
			c, err := getPokemonCount()
			return c >= 10, err
		}},
		{entity.AchievementPokemon25, func() (bool, error) {
			c, err := getPokemonCount()
			return c >= 25, err
		}},
		{entity.AchievementPokemon50, func() (bool, error) {
			c, err := getPokemonCount()
			return c >= 50, err
		}},
		{entity.AchievementPokemon100, func() (bool, error) {
			c, err := getPokemonCount()
			return c >= 100, err
		}},
		{entity.AchievementPokemon200, func() (bool, error) {
			c, err := getPokemonCount()
			return c >= 200, err
		}},
		{entity.AchievementPokemon500, func() (bool, error) {
			c, err := getPokemonCount()
			return c >= 500, err
		}},

		// --- レアリティ系 ---
		{entity.AchievementRareCatch, func() (bool, error) {
			return caughtRarity == valueobject.RarityRare, nil
		}},
		{entity.AchievementEpicCatch, func() (bool, error) {
			return caughtRarity == valueobject.RarityEpic, nil
		}},
		{entity.AchievementLegendary, func() (bool, error) {
			return caughtRarity == valueobject.RarityLegendary, nil
		}},

		// --- PR系 ---
		{entity.AchievementFirstMerge, func() (bool, error) {
			c, err := getMergeCount()
			return c >= 1, err
		}},
		{entity.AchievementMerge10, func() (bool, error) {
			c, err := getMergeCount()
			return c >= 10, err
		}},
		{entity.AchievementMerge25, func() (bool, error) {
			c, err := getMergeCount()
			return c >= 25, err
		}},
		{entity.AchievementMerge50, func() (bool, error) {
			c, err := getMergeCount()
			return c >= 50, err
		}},
		{entity.AchievementMerge100, func() (bool, error) {
			c, err := getMergeCount()
			return c >= 100, err
		}},

		// --- レビュー系 ---
		{entity.AchievementFirstReview, func() (bool, error) {
			c, err := getReviewCount()
			return c >= 1, err
		}},
		{entity.AchievementReview10, func() (bool, error) {
			c, err := getReviewCount()
			return c >= 10, err
		}},
		{entity.AchievementReview25, func() (bool, error) {
			c, err := getReviewCount()
			return c >= 25, err
		}},
		{entity.AchievementReview50, func() (bool, error) {
			c, err := getReviewCount()
			return c >= 50, err
		}},
		{entity.AchievementReview100, func() (bool, error) {
			c, err := getReviewCount()
			return c >= 100, err
		}},

		// --- Issue系 ---
		{entity.AchievementFirstIssue, func() (bool, error) {
			c, err := getIssueCount()
			return c >= 1, err
		}},
		{entity.AchievementIssue10, func() (bool, error) {
			c, err := getIssueCount()
			return c >= 10, err
		}},
		{entity.AchievementIssue50, func() (bool, error) {
			c, err := getIssueCount()
			return c >= 50, err
		}},

		// --- 特殊系 ---
		{entity.AchievementAllRounder, func() (bool, error) {
			commits, err := getCommitCount()
			if err != nil || commits < 1 {
				return false, err
			}
			merges, err := getMergeCount()
			if err != nil || merges < 1 {
				return false, err
			}
			reviews, err := getReviewCount()
			if err != nil || reviews < 1 {
				return false, err
			}
			issues, err := getIssueCount()
			if err != nil || issues < 1 {
				return false, err
			}
			return true, nil
		}},
	}

	for _, check := range checks {
		has, err := s.achievementRepo.HasAchievement(ctx, user.ID, check.achievementType)
		if err != nil {
			return nil, err
		}
		if has {
			continue
		}

		met, err := check.condition()
		if err != nil {
			return nil, err
		}
		if !met {
			continue
		}

		achievement := entity.NewUserAchievement(user.ID, check.achievementType)
		if err := s.achievementRepo.Save(ctx, achievement); err != nil {
			return nil, err
		}
		newAchievements = append(newAchievements, achievement)
	}

	return newAchievements, nil
}

func (s *AchievementService) checkPokemonCount(ctx context.Context, userID uuid.UUID, required int) (bool, error) {
	pokemon, err := s.pokemonRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	return len(pokemon) >= required, nil
}

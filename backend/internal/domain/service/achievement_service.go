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
// Returns the list of newly unlocked achievements.
func (s *AchievementService) CheckAndUnlock(ctx context.Context, user *entity.User, caughtRarity valueobject.PokemonRarity) ([]*entity.UserAchievement, error) {
	var newAchievements []*entity.UserAchievement

	checks := []struct {
		achievementType entity.AchievementType
		condition       func() (bool, error)
	}{
		{entity.AchievementFirstCommit, func() (bool, error) {
			count, err := s.activityRepo.CountByUserID(ctx, user.ID)
			return count >= 1, err
		}},
		{entity.AchievementStreak7, func() (bool, error) {
			return user.CurrentStreak >= 7, nil
		}},
		{entity.AchievementStreak30, func() (bool, error) {
			return user.CurrentStreak >= 30, nil
		}},
		{entity.AchievementStreak100, func() (bool, error) {
			return user.CurrentStreak >= 100, nil
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
		{entity.AchievementPokemon10, func() (bool, error) {
			return s.checkPokemonCount(ctx, user.ID, 10)
		}},
		{entity.AchievementPokemon50, func() (bool, error) {
			return s.checkPokemonCount(ctx, user.ID, 50)
		}},
		{entity.AchievementPokemon100, func() (bool, error) {
			return s.checkPokemonCount(ctx, user.ID, 100)
		}},
		{entity.AchievementReview10, func() (bool, error) {
			count, err := s.activityRepo.CountByUserIDAndType(ctx, user.ID, entity.EventReview)
			return count >= 10, err
		}},
		{entity.AchievementLegendary, func() (bool, error) {
			return caughtRarity == valueobject.RarityLegendary, nil
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

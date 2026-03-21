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

// achievementCheck pairs an achievement type with its unlock condition.
type achievementCheck struct {
	achievementType entity.AchievementType
	condition       func() (bool, error)
}

// thresholdChecks generates checks for a series of thresholds against a lazy-loaded counter.
func thresholdChecks(pairs []struct {
	t         entity.AchievementType
	threshold int64
}, getter func() (int64, error)) []achievementCheck {
	checks := make([]achievementCheck, len(pairs))
	for i, p := range pairs {
		threshold := p.threshold
		checks[i] = achievementCheck{p.t, func() (bool, error) {
			v, err := getter()
			return v >= threshold, err
		}}
	}
	return checks
}

// intThresholdChecks is like thresholdChecks but for int getters (e.g. pokemon count).
func intThresholdChecks(pairs []struct {
	t         entity.AchievementType
	threshold int
}, getter func() (int, error)) []achievementCheck {
	checks := make([]achievementCheck, len(pairs))
	for i, p := range pairs {
		threshold := p.threshold
		checks[i] = achievementCheck{p.t, func() (bool, error) {
			v, err := getter()
			return v >= threshold, err
		}}
	}
	return checks
}

// valueChecks generates checks that compare a user field against thresholds (no DB call).
func valueChecks(pairs []struct {
	t         entity.AchievementType
	threshold int
}, value int) []achievementCheck {
	checks := make([]achievementCheck, len(pairs))
	for i, p := range pairs {
		met := value >= p.threshold
		checks[i] = achievementCheck{p.t, func() (bool, error) { return met, nil }}
	}
	return checks
}

// CheckAndUnlock checks all achievement conditions and unlocks any newly earned ones.
func (s *AchievementService) CheckAndUnlock(ctx context.Context, user *entity.User, caughtRarity valueobject.PokemonRarity) ([]*entity.UserAchievement, error) {
	counters := newLazyCounters(ctx, user.ID, s.activityRepo, s.pokemonRepo)
	checks := s.buildChecks(user, caughtRarity, counters)
	return s.evaluateAndUnlock(ctx, user.ID, checks)
}

// buildChecks assembles all achievement checks organized by category.
func (s *AchievementService) buildChecks(user *entity.User, caughtRarity valueobject.PokemonRarity, c *lazyCounters) []achievementCheck {
	var all []achievementCheck

	// コミット系
	all = append(all, thresholdChecks([]struct {
		t         entity.AchievementType
		threshold int64
	}{
		{entity.AchievementFirstCommit, 1},
		{entity.AchievementCommit50, 50},
		{entity.AchievementCommit100, 100},
		{entity.AchievementCommit500, 500},
		{entity.AchievementCommit1000, 1000},
	}, c.commits)...)

	// ストリーク系
	all = append(all, valueChecks([]struct {
		t         entity.AchievementType
		threshold int
	}{
		{entity.AchievementStreak3, 3},
		{entity.AchievementStreak7, 7},
		{entity.AchievementStreak14, 14},
		{entity.AchievementStreak30, 30},
		{entity.AchievementStreak60, 60},
		{entity.AchievementStreak100, 100},
		{entity.AchievementStreak365, 365},
	}, user.CurrentStreak)...)

	// レベル系
	all = append(all, valueChecks([]struct {
		t         entity.AchievementType
		threshold int
	}{
		{entity.AchievementLevel5, 5},
		{entity.AchievementLevel10, 10},
		{entity.AchievementLevel25, 25},
		{entity.AchievementLevel50, 50},
		{entity.AchievementLevel100, 100},
	}, user.Level)...)

	// XP系
	all = append(all, valueChecks([]struct {
		t         entity.AchievementType
		threshold int
	}{
		{entity.AchievementXP1000, 1000},
		{entity.AchievementXP5000, 5000},
		{entity.AchievementXP10000, 10000},
		{entity.AchievementXP50000, 50000},
	}, user.TotalXP)...)

	// ポケモン系
	all = append(all, intThresholdChecks([]struct {
		t         entity.AchievementType
		threshold int
	}{
		{entity.AchievementPokemon1, 1},
		{entity.AchievementPokemon10, 10},
		{entity.AchievementPokemon25, 25},
		{entity.AchievementPokemon50, 50},
		{entity.AchievementPokemon100, 100},
		{entity.AchievementPokemon200, 200},
		{entity.AchievementPokemon500, 500},
	}, c.pokemon)...)

	// レアリティ系
	all = append(all, s.rarityChecks(caughtRarity)...)

	// PR系
	all = append(all, thresholdChecks([]struct {
		t         entity.AchievementType
		threshold int64
	}{
		{entity.AchievementFirstMerge, 1},
		{entity.AchievementMerge10, 10},
		{entity.AchievementMerge25, 25},
		{entity.AchievementMerge50, 50},
		{entity.AchievementMerge100, 100},
	}, c.merges)...)

	// レビュー系
	all = append(all, thresholdChecks([]struct {
		t         entity.AchievementType
		threshold int64
	}{
		{entity.AchievementFirstReview, 1},
		{entity.AchievementReview10, 10},
		{entity.AchievementReview25, 25},
		{entity.AchievementReview50, 50},
		{entity.AchievementReview100, 100},
	}, c.reviews)...)

	// Issue系
	all = append(all, thresholdChecks([]struct {
		t         entity.AchievementType
		threshold int64
	}{
		{entity.AchievementFirstIssue, 1},
		{entity.AchievementIssue10, 10},
		{entity.AchievementIssue50, 50},
	}, c.issues)...)

	// デプロイ系
	all = append(all, thresholdChecks([]struct {
		t         entity.AchievementType
		threshold int64
	}{
		{entity.AchievementFirstDeploy, 1},
		{entity.AchievementDeploy10, 10},
		{entity.AchievementDeploy50, 50},
		{entity.AchievementDeploy100, 100},
	}, c.deploys)...)

	// 特殊系
	all = append(all, s.specialChecks(c)...)

	return all
}

func (s *AchievementService) rarityChecks(rarity valueobject.PokemonRarity) []achievementCheck {
	return []achievementCheck{
		{entity.AchievementRareCatch, func() (bool, error) {
			return rarity == valueobject.RarityRare, nil
		}},
		{entity.AchievementEpicCatch, func() (bool, error) {
			return rarity == valueobject.RarityEpic, nil
		}},
		{entity.AchievementLegendary, func() (bool, error) {
			return rarity == valueobject.RarityLegendary, nil
		}},
	}
}

func (s *AchievementService) specialChecks(c *lazyCounters) []achievementCheck {
	return []achievementCheck{
		{entity.AchievementAllRounder, func() (bool, error) {
			commits, err := c.commits()
			if err != nil || commits < 1 {
				return false, err
			}
			merges, err := c.merges()
			if err != nil || merges < 1 {
				return false, err
			}
			reviews, err := c.reviews()
			if err != nil || reviews < 1 {
				return false, err
			}
			issues, err := c.issues()
			if err != nil || issues < 1 {
				return false, err
			}
			return true, nil
		}},
	}
}

// evaluateAndUnlock iterates checks, skips already-unlocked, and saves new achievements.
func (s *AchievementService) evaluateAndUnlock(ctx context.Context, userID uuid.UUID, checks []achievementCheck) ([]*entity.UserAchievement, error) {
	var newAchievements []*entity.UserAchievement

	for _, check := range checks {
		has, err := s.achievementRepo.HasAchievement(ctx, userID, check.achievementType)
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

		achievement := entity.NewUserAchievement(userID, check.achievementType)
		if err := s.achievementRepo.Save(ctx, achievement); err != nil {
			return nil, err
		}
		newAchievements = append(newAchievements, achievement)
	}

	return newAchievements, nil
}

// lazyCounters caches DB queries so each counter is fetched at most once.
type lazyCounters struct {
	ctx         context.Context
	userID      uuid.UUID
	activityRepo repository.GitHubActivityRepository
	pokemonRepo  repository.PokemonRepository

	commitCount  *int64
	mergeCount   *int64
	reviewCount  *int64
	issueCount   *int64
	deployCount  *int64
	pokemonCount *int
}

func newLazyCounters(ctx context.Context, userID uuid.UUID, activityRepo repository.GitHubActivityRepository, pokemonRepo repository.PokemonRepository) *lazyCounters {
	return &lazyCounters{
		ctx:          ctx,
		userID:       userID,
		activityRepo: activityRepo,
		pokemonRepo:  pokemonRepo,
	}
}

func (lc *lazyCounters) commits() (int64, error) {
	if lc.commitCount == nil {
		c, err := lc.activityRepo.CountByUserID(lc.ctx, lc.userID)
		if err != nil {
			return 0, err
		}
		lc.commitCount = &c
	}
	return *lc.commitCount, nil
}

func (lc *lazyCounters) merges() (int64, error) {
	if lc.mergeCount == nil {
		c, err := lc.activityRepo.CountByUserIDAndType(lc.ctx, lc.userID, entity.EventPRMerge)
		if err != nil {
			return 0, err
		}
		lc.mergeCount = &c
	}
	return *lc.mergeCount, nil
}

func (lc *lazyCounters) reviews() (int64, error) {
	if lc.reviewCount == nil {
		c, err := lc.activityRepo.CountByUserIDAndType(lc.ctx, lc.userID, entity.EventReview)
		if err != nil {
			return 0, err
		}
		lc.reviewCount = &c
	}
	return *lc.reviewCount, nil
}

func (lc *lazyCounters) issues() (int64, error) {
	if lc.issueCount == nil {
		c, err := lc.activityRepo.CountByUserIDAndType(lc.ctx, lc.userID, entity.EventIssueClose)
		if err != nil {
			return 0, err
		}
		lc.issueCount = &c
	}
	return *lc.issueCount, nil
}

func (lc *lazyCounters) deploys() (int64, error) {
	if lc.deployCount == nil {
		c, err := lc.activityRepo.CountByUserIDAndType(lc.ctx, lc.userID, entity.EventDeploy)
		if err != nil {
			return 0, err
		}
		lc.deployCount = &c
	}
	return *lc.deployCount, nil
}

func (lc *lazyCounters) pokemon() (int, error) {
	if lc.pokemonCount == nil {
		p, err := lc.pokemonRepo.FindByUserID(lc.ctx, lc.userID)
		if err != nil {
			return 0, err
		}
		c := len(p)
		lc.pokemonCount = &c
	}
	return *lc.pokemonCount, nil
}

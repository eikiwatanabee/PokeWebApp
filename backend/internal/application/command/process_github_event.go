package command

import (
	"context"
	"fmt"
	"math"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/service"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
)

type ProcessGitHubEventCommand struct {
	GitHubUsername string
	EventType     entity.GitHubEventType
	RepoName      string
	Title         string
	URL           string
}

type NewAchievementDTO struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type ProcessGitHubEventResult struct {
	XPGained        int                `json:"xp_gained"`
	MissionBonusXP  int                `json:"mission_bonus_xp,omitempty"`
	TotalXP         int                `json:"total_xp"`
	Level           int                `json:"level"`
	CurrentStreak   int                `json:"current_streak"`
	StreakMultiplier float64            `json:"streak_multiplier"`
	PokemonName     string             `json:"pokemon_name,omitempty"`
	SpriteURL       string             `json:"sprite_url,omitempty"`
	PokemonRarity   string             `json:"pokemon_rarity,omitempty"`
	NewAchievements []NewAchievementDTO `json:"new_achievements,omitempty"`
	DeployBoost     bool               `json:"deploy_boost,omitempty"`
}

type ProcessGitHubEventHandler struct {
	uow              uow.UnitOfWork
	userRepo         repository.UserRepository
	activityRepo     repository.GitHubActivityRepository
	pokemonRepo      repository.PokemonRepository
	limitedEventRepo repository.LimitedEventRepository
	deployerRepo     repository.DeployerRepository
	gachaSvc         *service.PokemonGachaService
	achievementSvc   *service.AchievementService
	missionSvc       *service.DailyMissionService
}

func NewProcessGitHubEventHandler(
	uow uow.UnitOfWork,
	userRepo repository.UserRepository,
	activityRepo repository.GitHubActivityRepository,
	pokemonRepo repository.PokemonRepository,
	gachaSvc *service.PokemonGachaService,
	achievementSvc *service.AchievementService,
	missionSvc *service.DailyMissionService,
	limitedEventRepo repository.LimitedEventRepository,
	deployerRepo repository.DeployerRepository,
) *ProcessGitHubEventHandler {
	return &ProcessGitHubEventHandler{
		uow:              uow,
		userRepo:         userRepo,
		activityRepo:     activityRepo,
		pokemonRepo:      pokemonRepo,
		limitedEventRepo: limitedEventRepo,
		deployerRepo:     deployerRepo,
		gachaSvc:         gachaSvc,
		achievementSvc:   achievementSvc,
		missionSvc:       missionSvc,
	}
}

func (h *ProcessGitHubEventHandler) Handle(ctx context.Context, cmd *ProcessGitHubEventCommand) (*ProcessGitHubEventResult, error) {
	user, err := h.userRepo.FindByGitHubUsername(ctx, cmd.GitHubUsername)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	// Deploy events require deployer permission
	if cmd.EventType == entity.EventDeploy && h.deployerRepo != nil {
		isDeployer, err := h.deployerRepo.IsDeployer(ctx, user.TenantID, user.ID)
		if err != nil {
			return nil, err
		}
		if !isDeployer {
			return nil, fmt.Errorf("user is not authorized as deployer")
		}
	}

	var result *ProcessGitHubEventResult
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		// Create activity record
		activity := entity.NewGitHubActivity(user.ID, cmd.EventType, cmd.RepoName, cmd.Title, cmd.URL)
		if err := h.activityRepo.Save(ctx, activity); err != nil {
			return err
		}

		// Update streak
		streakMultiplier := user.UpdateStreak(activity.CreatedAt)

		// Apply streak bonus to XP
		baseXP := activity.XP
		bonusXP := int(math.Round(float64(baseXP) * (streakMultiplier - 1)))
		totalXPGained := baseXP + bonusXP

		user.AddXP(totalXPGained)
		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		result = &ProcessGitHubEventResult{
			XPGained:        totalXPGained,
			TotalXP:         user.TotalXP,
			Level:           user.Level,
			CurrentStreak:   user.CurrentStreak,
			StreakMultiplier: streakMultiplier,
		}

		var caughtRarity valueobject.PokemonRarity

		// On PR merge, catch a Pokemon with rarity boost!
		if cmd.EventType == entity.EventPRMerge {
			// Apply limited event boost on top of streak multiplier
			effectiveMultiplier := streakMultiplier
			if h.limitedEventRepo != nil {
				activeEvents, err := h.limitedEventRepo.FindActiveByTenantID(ctx, user.TenantID)
				if err == nil {
					for _, ev := range activeEvents {
						effectiveMultiplier *= ev.RarityBoost
					}
				}
			}

			pokemonInfo, err := h.gachaSvc.DrawWithBoost(ctx, effectiveMultiplier)
			if err != nil {
				return err
			}

			userPokemon := entity.NewUserPokemonFromActivity(user.ID, activity.ID, *pokemonInfo)
			if err := h.pokemonRepo.Save(ctx, userPokemon); err != nil {
				return err
			}

			result.PokemonName = pokemonInfo.Name
			result.SpriteURL = pokemonInfo.SpriteURL
			result.PokemonRarity = string(pokemonInfo.Rarity)
			caughtRarity = pokemonInfo.Rarity
		}

		// On deploy, catch a Pokemon (200 XP bonus already applied above)
		if cmd.EventType == entity.EventDeploy {
			pokemonInfo, err := h.gachaSvc.DrawWithBoost(ctx, streakMultiplier)
			if err != nil {
				return err
			}

			userPokemon := entity.NewUserPokemonFromActivity(user.ID, activity.ID, *pokemonInfo)
			if err := h.pokemonRepo.Save(ctx, userPokemon); err != nil {
				return err
			}

			result.PokemonName = pokemonInfo.Name
			result.SpriteURL = pokemonInfo.SpriteURL
			result.PokemonRarity = string(pokemonInfo.Rarity)
			result.DeployBoost = true
			caughtRarity = pokemonInfo.Rarity
		}

		// Check achievements
		newAchievements, err := h.achievementSvc.CheckAndUnlock(ctx, user, caughtRarity)
		if err != nil {
			return err
		}

		if len(newAchievements) > 0 {
			defMap := make(map[entity.AchievementType]entity.AchievementDefinition)
			for _, d := range entity.AchievementDefinitions {
				defMap[d.Type] = d
			}
			for _, a := range newAchievements {
				if def, ok := defMap[a.AchievementType]; ok {
					result.NewAchievements = append(result.NewAchievements, NewAchievementDTO{
						Type: string(a.AchievementType),
						Name: def.Name,
						Icon: def.Icon,
					})
				}
			}
		}

		// Process daily missions
		if h.missionSvc != nil {
			missionBonus, err := h.missionSvc.ProcessActivity(ctx, user.ID, cmd.EventType)
			if err != nil {
				return err
			}
			if missionBonus > 0 {
				user.AddXP(missionBonus)
				if err := h.userRepo.Save(ctx, user); err != nil {
					return err
				}
				result.MissionBonusXP = missionBonus
				result.TotalXP = user.TotalXP
				result.Level = user.Level
			}
		}

		return nil
	})

	return result, err
}

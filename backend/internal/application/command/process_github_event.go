package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/service"
)

type ProcessGitHubEventCommand struct {
	GitHubUsername string
	EventType     entity.GitHubEventType
	RepoName      string
	Title         string
	URL           string
}

type ProcessGitHubEventResult struct {
	XPGained    int    `json:"xp_gained"`
	TotalXP     int    `json:"total_xp"`
	Level       int    `json:"level"`
	PokemonName string `json:"pokemon_name,omitempty"`
	SpriteURL   string `json:"sprite_url,omitempty"`
}

type ProcessGitHubEventHandler struct {
	uow          uow.UnitOfWork
	userRepo     repository.UserRepository
	activityRepo repository.GitHubActivityRepository
	pokemonRepo  repository.PokemonRepository
	gachaSvc     *service.PokemonGachaService
}

func NewProcessGitHubEventHandler(
	uow uow.UnitOfWork,
	userRepo repository.UserRepository,
	activityRepo repository.GitHubActivityRepository,
	pokemonRepo repository.PokemonRepository,
	gachaSvc *service.PokemonGachaService,
) *ProcessGitHubEventHandler {
	return &ProcessGitHubEventHandler{
		uow:          uow,
		userRepo:     userRepo,
		activityRepo: activityRepo,
		pokemonRepo:  pokemonRepo,
		gachaSvc:     gachaSvc,
	}
}

func (h *ProcessGitHubEventHandler) Handle(ctx context.Context, cmd *ProcessGitHubEventCommand) (*ProcessGitHubEventResult, error) {
	// Look up user by GitHub username
	user, err := h.userRepo.FindByGitHubUsername(ctx, cmd.GitHubUsername)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// User not registered, skip silently
		return nil, nil
	}

	var result *ProcessGitHubEventResult
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		// Create activity record
		activity := entity.NewGitHubActivity(user.ID, cmd.EventType, cmd.RepoName, cmd.Title, cmd.URL)
		if err := h.activityRepo.Save(ctx, activity); err != nil {
			return err
		}

		// Add XP to user
		user.AddXP(activity.XP)
		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		result = &ProcessGitHubEventResult{
			XPGained: activity.XP,
			TotalXP:  user.TotalXP,
			Level:    user.Level,
		}

		// On PR merge, catch a Pokemon!
		if cmd.EventType == entity.EventPRMerge {
			pokemonInfo, err := h.gachaSvc.Draw(ctx)
			if err != nil {
				return err
			}

			userPokemon := entity.NewUserPokemonFromActivity(user.ID, activity.ID, *pokemonInfo)
			if err := h.pokemonRepo.Save(ctx, userPokemon); err != nil {
				return err
			}

			result.PokemonName = pokemonInfo.Name
			result.SpriteURL = pokemonInfo.SpriteURL
		}

		return nil
	})

	return result, err
}

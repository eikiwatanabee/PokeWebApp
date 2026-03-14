package command

import (
	"context"
	"errors"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/service"
	"github.com/google/uuid"
)

var (
	ErrInvalidStarter    = errors.New("invalid starter pokemon: must be 1 (bulbasaur), 4 (charmander), or 7 (squirtle)")
	ErrAlreadyHasStarter = errors.New("user already has a starter pokemon")
)

var validStarters = map[int]bool{
	1: true, // Bulbasaur
	4: true, // Charmander
	7: true, // Squirtle
}

type ChooseStarterCommand struct {
	UserID    string
	PokemonID int
}

type ChooseStarterResult struct {
	PokemonID   int    `json:"pokemon_id"`
	PokemonName string `json:"pokemon_name"`
	SpriteURL   string `json:"sprite_url"`
}

type ChooseStarterHandler struct {
	uow         uow.UnitOfWork
	pokemonRepo repository.PokemonRepository
	gachaSvc    *service.PokemonGachaService
}

func NewChooseStarterHandler(
	uow uow.UnitOfWork,
	pokemonRepo repository.PokemonRepository,
	gachaSvc *service.PokemonGachaService,
) *ChooseStarterHandler {
	return &ChooseStarterHandler{uow: uow, pokemonRepo: pokemonRepo, gachaSvc: gachaSvc}
}

func (h *ChooseStarterHandler) Handle(ctx context.Context, cmd *ChooseStarterCommand) (*ChooseStarterResult, error) {
	if !validStarters[cmd.PokemonID] {
		return nil, ErrInvalidStarter
	}

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}

	var result *ChooseStarterResult
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		existing, err := h.pokemonRepo.FindByUserID(ctx, userID)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return ErrAlreadyHasStarter
		}

		pokemonInfo, err := h.gachaSvc.FetchByID(ctx, cmd.PokemonID)
		if err != nil {
			return err
		}

		// Starter is not tied to a book, use nil UUID
		userPokemon := entity.NewUserPokemon(userID, uuid.Nil, *pokemonInfo)
		if err := h.pokemonRepo.Save(ctx, userPokemon); err != nil {
			return err
		}

		result = &ChooseStarterResult{
			PokemonID:   pokemonInfo.PokemonID,
			PokemonName: pokemonInfo.Name,
			SpriteURL:   pokemonInfo.SpriteURL,
		}
		return nil
	})
	return result, err
}

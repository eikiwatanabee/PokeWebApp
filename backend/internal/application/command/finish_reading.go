package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/service"
	"github.com/google/uuid"
)

type FinishReadingCommand struct {
	BookID string
	UserID string
}

type FinishReadingResult struct {
	PokemonID   int    `json:"pokemon_id"`
	PokemonName string `json:"pokemon_name"`
	SpriteURL   string `json:"sprite_url"`
}

type FinishReadingHandler struct {
	uow         uow.UnitOfWork
	bookRepo    repository.BookRepository
	pokemonRepo repository.PokemonRepository
	gachaSvc    *service.PokemonGachaService
}

func NewFinishReadingHandler(
	uow uow.UnitOfWork,
	bookRepo repository.BookRepository,
	pokemonRepo repository.PokemonRepository,
	gachaSvc *service.PokemonGachaService,
) *FinishReadingHandler {
	return &FinishReadingHandler{
		uow: uow, bookRepo: bookRepo,
		pokemonRepo: pokemonRepo, gachaSvc: gachaSvc,
	}
}

func (h *FinishReadingHandler) Handle(ctx context.Context, cmd *FinishReadingCommand) (*FinishReadingResult, error) {
	bookID, err := uuid.Parse(cmd.BookID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}

	var result *FinishReadingResult
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		// Lock order: Book(3)
		book, err := h.bookRepo.FindByIDForUpdate(ctx, bookID)
		if err != nil {
			return err
		}

		// Domain logic
		if err := book.Finish(); err != nil {
			return err
		}
		if err := h.bookRepo.Save(ctx, book); err != nil {
			return err
		}

		// Lock order: UserPokemon(6) — after Book(3) ✓
		pokemonInfo, err := h.gachaSvc.Draw(ctx)
		if err != nil {
			return err
		}

		userPokemon := entity.NewUserPokemon(userID, bookID, *pokemonInfo)
		if err := h.pokemonRepo.Save(ctx, userPokemon); err != nil {
			return err
		}

		result = &FinishReadingResult{
			PokemonID:   pokemonInfo.PokemonID,
			PokemonName: pokemonInfo.Name,
			SpriteURL:   pokemonInfo.SpriteURL,
		}
		return nil
	})
	return result, err
}

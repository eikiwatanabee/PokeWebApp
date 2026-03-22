package command

import (
	"context"
	"fmt"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

// --- Create Trade ---

type CreateTradeCommand struct {
	TenantID             string
	UserID               string
	OfferedPokemonID     string
	RequestedPokemonName string
}

type CreateTradeHandler struct {
	uow         uow.UnitOfWork
	tradeRepo   repository.TradeRepository
	pokemonRepo repository.PokemonRepository
}

func NewCreateTradeHandler(uow uow.UnitOfWork, tradeRepo repository.TradeRepository, pokemonRepo repository.PokemonRepository) *CreateTradeHandler {
	return &CreateTradeHandler{uow: uow, tradeRepo: tradeRepo, pokemonRepo: pokemonRepo}
}

func (h *CreateTradeHandler) Handle(ctx context.Context, cmd *CreateTradeCommand) (string, error) {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}
	tenantID, err := uuid.Parse(cmd.TenantID)
	if err != nil {
		return "", fmt.Errorf("invalid tenant ID: %w", err)
	}
	pokemonID, err := uuid.Parse(cmd.OfferedPokemonID)
	if err != nil {
		return "", fmt.Errorf("invalid pokemon ID: %w", err)
	}

	// Verify ownership
	pokemon, err := h.pokemonRepo.FindByID(ctx, pokemonID)
	if err != nil {
		return "", err
	}
	if pokemon == nil || pokemon.UserID != userID {
		return "", entity.ErrPokemonNotOwned
	}

	trade := entity.NewTrade(tenantID, userID, pokemonID, cmd.RequestedPokemonName)

	var tradeID string
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		if err := h.tradeRepo.Save(ctx, trade); err != nil {
			return err
		}
		tradeID = trade.ID.String()
		return nil
	})

	return tradeID, err
}

// --- Accept Trade ---

type AcceptTradeCommand struct {
	TradeID         string
	UserID          string
	OfferedPokemonID string // Pokemon the accepter is giving
}

type AcceptTradeHandler struct {
	uow         uow.UnitOfWork
	tradeRepo   repository.TradeRepository
	pokemonRepo repository.PokemonRepository
}

func NewAcceptTradeHandler(uow uow.UnitOfWork, tradeRepo repository.TradeRepository, pokemonRepo repository.PokemonRepository) *AcceptTradeHandler {
	return &AcceptTradeHandler{uow: uow, tradeRepo: tradeRepo, pokemonRepo: pokemonRepo}
}

func (h *AcceptTradeHandler) Handle(ctx context.Context, cmd *AcceptTradeCommand) error {
	tradeID, err := uuid.Parse(cmd.TradeID)
	if err != nil {
		return fmt.Errorf("invalid trade ID: %w", err)
	}
	accepterID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	accepterPokemonID, err := uuid.Parse(cmd.OfferedPokemonID)
	if err != nil {
		return fmt.Errorf("invalid pokemon ID: %w", err)
	}

	return h.uow.Do(ctx, func(ctx context.Context) error {
		trade, err := h.tradeRepo.FindByID(ctx, tradeID)
		if err != nil {
			return err
		}
		if trade == nil {
			return fmt.Errorf("trade not found")
		}

		// Verify accepter owns the pokemon they're offering
		accepterPokemon, err := h.pokemonRepo.FindByID(ctx, accepterPokemonID)
		if err != nil {
			return err
		}
		if accepterPokemon == nil || accepterPokemon.UserID != accepterID {
			return entity.ErrPokemonNotOwned
		}

		// Verify offerer still owns the offered pokemon
		offeredPokemon, err := h.pokemonRepo.FindByID(ctx, trade.OfferedPokemonID)
		if err != nil {
			return err
		}
		if offeredPokemon == nil || offeredPokemon.UserID != trade.OffererID {
			return entity.ErrPokemonNotOwned
		}

		// Check requested pokemon name matches (if specified)
		if trade.RequestedPokemonName != "" && accepterPokemon.Pokemon.Name != trade.RequestedPokemonName {
			return fmt.Errorf("offered pokemon does not match requested: %s", trade.RequestedPokemonName)
		}

		// Accept the trade
		if err := trade.Accept(accepterID, accepterPokemonID); err != nil {
			return err
		}

		// Swap ownership
		offeredPokemon.UserID = accepterID
		accepterPokemon.UserID = trade.OffererID

		if err := h.pokemonRepo.Save(ctx, offeredPokemon); err != nil {
			return err
		}
		if err := h.pokemonRepo.Save(ctx, accepterPokemon); err != nil {
			return err
		}

		return h.tradeRepo.Save(ctx, trade)
	})
}

// --- Cancel Trade ---

type CancelTradeCommand struct {
	TradeID string
	UserID  string
}

type CancelTradeHandler struct {
	tradeRepo repository.TradeRepository
}

func NewCancelTradeHandler(tradeRepo repository.TradeRepository) *CancelTradeHandler {
	return &CancelTradeHandler{tradeRepo: tradeRepo}
}

func (h *CancelTradeHandler) Handle(ctx context.Context, cmd *CancelTradeCommand) error {
	tradeID, err := uuid.Parse(cmd.TradeID)
	if err != nil {
		return fmt.Errorf("invalid trade ID: %w", err)
	}
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	trade, err := h.tradeRepo.FindByID(ctx, tradeID)
	if err != nil {
		return err
	}
	if trade == nil {
		return fmt.Errorf("trade not found")
	}

	if err := trade.Cancel(userID); err != nil {
		return err
	}

	return h.tradeRepo.Save(ctx, trade)
}

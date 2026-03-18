package query

import (
	"context"
	"fmt"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetTradesQuery struct {
	TenantID string
	Limit    int
}

type TradeDTO struct {
	ID                   string  `json:"id"`
	OffererID            string  `json:"offerer_id"`
	OffererName          string  `json:"offerer_name"`
	OffererAvatarURL     string  `json:"offerer_avatar_url"`
	OfferedPokemonID     string  `json:"offered_pokemon_id"`
	OfferedPokemonName   string  `json:"offered_pokemon_name"`
	OfferedPokemonSprite string  `json:"offered_pokemon_sprite"`
	OfferedPokemonRarity string  `json:"offered_pokemon_rarity"`
	RequestedPokemonName string  `json:"requested_pokemon_name"`
	Status               string  `json:"status"`
	CreatedAt            string  `json:"created_at"`
}

type GetTradesResult struct {
	Trades []TradeDTO `json:"trades"`
}

type GetTradesHandler struct {
	tradeRepo   repository.TradeRepository
	pokemonRepo repository.PokemonRepository
	userRepo    repository.UserRepository
}

func NewGetTradesHandler(
	tradeRepo repository.TradeRepository,
	pokemonRepo repository.PokemonRepository,
	userRepo repository.UserRepository,
) *GetTradesHandler {
	return &GetTradesHandler{tradeRepo: tradeRepo, pokemonRepo: pokemonRepo, userRepo: userRepo}
}

func (h *GetTradesHandler) Handle(ctx context.Context, q *GetTradesQuery) (*GetTradesResult, error) {
	tenantID, err := uuid.Parse(q.TenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	limit := q.Limit
	if limit <= 0 {
		limit = 50
	}

	trades, err := h.tradeRepo.FindOpenByTenantID(ctx, tenantID, limit)
	if err != nil {
		return nil, err
	}

	// Collect user IDs and pokemon IDs for batch resolution
	users, err := h.userRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	userMap := make(map[uuid.UUID]struct {
		Name      string
		AvatarURL string
	})
	for _, u := range users {
		userMap[u.ID] = struct {
			Name      string
			AvatarURL string
		}{Name: u.Name, AvatarURL: u.AvatarURL}
	}

	dtos := make([]TradeDTO, len(trades))
	for i, t := range trades {
		offerer := userMap[t.OffererID]

		dto := TradeDTO{
			ID:                   t.ID.String(),
			OffererID:            t.OffererID.String(),
			OffererName:          offerer.Name,
			OffererAvatarURL:     offerer.AvatarURL,
			OfferedPokemonID:     t.OfferedPokemonID.String(),
			RequestedPokemonName: t.RequestedPokemonName,
			Status:               string(t.Status),
			CreatedAt:            t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// Get offered pokemon details
		pokemon, err := h.pokemonRepo.FindByID(ctx, t.OfferedPokemonID)
		if err == nil && pokemon != nil {
			dto.OfferedPokemonName = pokemon.Pokemon.Name
			dto.OfferedPokemonSprite = pokemon.Pokemon.SpriteURL
			dto.OfferedPokemonRarity = string(pokemon.Pokemon.Rarity)
		}

		dtos[i] = dto
	}

	return &GetTradesResult{Trades: dtos}, nil
}

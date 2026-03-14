package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetPokedexQuery struct {
	UserID string
}

type PokemonDTO struct {
	ID          string   `json:"id"`
	PokemonID   int      `json:"pokemon_id"`
	PokemonName string   `json:"pokemon_name"`
	SpriteURL   string   `json:"sprite_url"`
	Types       []string `json:"types"`
	BookID      string   `json:"book_id"`
	CaughtAt    string   `json:"caught_at"`
}

type GetPokedexResult struct {
	Pokemon []PokemonDTO `json:"pokemon"`
	Total   int          `json:"total"`
}

type GetPokedexHandler struct {
	pokemonRepo repository.PokemonRepository
}

func NewGetPokedexHandler(pokemonRepo repository.PokemonRepository) *GetPokedexHandler {
	return &GetPokedexHandler{pokemonRepo: pokemonRepo}
}

func (h *GetPokedexHandler) Handle(ctx context.Context, q *GetPokedexQuery) (*GetPokedexResult, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return nil, err
	}

	pokemons, err := h.pokemonRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]PokemonDTO, len(pokemons))
	for i, p := range pokemons {
		dtos[i] = PokemonDTO{
			ID:          p.ID.String(),
			PokemonID:   p.Pokemon.PokemonID,
			PokemonName: p.Pokemon.Name,
			SpriteURL:   p.Pokemon.SpriteURL,
			Types:       p.Pokemon.Types,
			BookID:      p.BookID.String(),
			CaughtAt:    p.CaughtAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &GetPokedexResult{
		Pokemon: dtos,
		Total:   len(dtos),
	}, nil
}

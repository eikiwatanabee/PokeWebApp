package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetTeamRankingQuery struct {
	TenantID string
}

type TeamRankingDTO struct {
	TeamID       string `json:"team_id"`
	TeamName     string `json:"team_name"`
	MemberCount  int    `json:"member_count"`
	PokemonCount int    `json:"pokemon_count"`
	BookCount    int    `json:"book_count"`
}

type GetTeamRankingResult struct {
	Teams []TeamRankingDTO `json:"teams"`
}

type GetTeamRankingHandler struct {
	teamRepo    repository.TeamRepository
	userRepo    repository.UserRepository
	pokemonRepo repository.PokemonRepository
	bookRepo    repository.BookRepository
}

func NewGetTeamRankingHandler(
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
	pokemonRepo repository.PokemonRepository,
	bookRepo repository.BookRepository,
) *GetTeamRankingHandler {
	return &GetTeamRankingHandler{
		teamRepo:    teamRepo,
		userRepo:    userRepo,
		pokemonRepo: pokemonRepo,
		bookRepo:    bookRepo,
	}
}

func (h *GetTeamRankingHandler) Handle(ctx context.Context, q *GetTeamRankingQuery) (*GetTeamRankingResult, error) {
	tenantID, err := uuid.Parse(q.TenantID)
	if err != nil {
		return nil, err
	}

	teams, err := h.teamRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	users, err := h.userRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Build team member map
	teamMembers := make(map[uuid.UUID][]uuid.UUID)
	for _, u := range users {
		if u.TeamID != nil {
			teamMembers[*u.TeamID] = append(teamMembers[*u.TeamID], u.ID)
		}
	}

	dtos := make([]TeamRankingDTO, len(teams))
	for i, team := range teams {
		memberIDs := teamMembers[team.ID]
		pokemonCount := 0
		bookCount := 0

		for _, memberID := range memberIDs {
			pokemons, err := h.pokemonRepo.FindByUserID(ctx, memberID)
			if err != nil {
				return nil, err
			}
			pokemonCount += len(pokemons)

			books, _, err := h.bookRepo.FindByUserID(ctx, memberID, "finished", "", 1, 9999)
			if err != nil {
				return nil, err
			}
			bookCount += len(books)
		}

		dtos[i] = TeamRankingDTO{
			TeamID:       team.ID.String(),
			TeamName:     team.Name,
			MemberCount:  len(memberIDs),
			PokemonCount: pokemonCount,
			BookCount:    bookCount,
		}
	}

	// Sort by pokemon count (simple bubble sort, teams are small)
	for i := 0; i < len(dtos); i++ {
		for j := i + 1; j < len(dtos); j++ {
			if dtos[j].PokemonCount > dtos[i].PokemonCount {
				dtos[i], dtos[j] = dtos[j], dtos[i]
			}
		}
	}

	return &GetTeamRankingResult{Teams: dtos}, nil
}

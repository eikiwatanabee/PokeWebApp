package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type CreateTeamCommand struct {
	TenantID string
	Name     string
}

type CreateTeamResult struct {
	TeamID string
}

type CreateTeamHandler struct {
	uow      uow.UnitOfWork
	teamRepo repository.TeamRepository
}

func NewCreateTeamHandler(uow uow.UnitOfWork, teamRepo repository.TeamRepository) *CreateTeamHandler {
	return &CreateTeamHandler{uow: uow, teamRepo: teamRepo}
}

func (h *CreateTeamHandler) Handle(ctx context.Context, cmd *CreateTeamCommand) (*CreateTeamResult, error) {
	tenantID, err := uuid.Parse(cmd.TenantID)
	if err != nil {
		return nil, err
	}

	var result *CreateTeamResult
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		team, err := entity.NewTeam(tenantID, cmd.Name)
		if err != nil {
			return err
		}
		if err := h.teamRepo.Save(ctx, team); err != nil {
			return err
		}
		result = &CreateTeamResult{TeamID: team.ID.String()}
		return nil
	})
	return result, err
}

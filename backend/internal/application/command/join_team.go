package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type JoinTeamCommand struct {
	UserID string
	TeamID string
}

type JoinTeamHandler struct {
	uow      uow.UnitOfWork
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewJoinTeamHandler(uow uow.UnitOfWork, userRepo repository.UserRepository, teamRepo repository.TeamRepository) *JoinTeamHandler {
	return &JoinTeamHandler{uow: uow, userRepo: userRepo, teamRepo: teamRepo}
}

func (h *JoinTeamHandler) Handle(ctx context.Context, cmd *JoinTeamCommand) error {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return err
	}
	teamID, err := uuid.Parse(cmd.TeamID)
	if err != nil {
		return err
	}

	return h.uow.Do(ctx, func(ctx context.Context) error {
		// Verify team exists
		if _, err := h.teamRepo.FindByID(ctx, teamID); err != nil {
			return err
		}

		user, err := h.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		user.JoinTeam(teamID)
		return h.userRepo.Save(ctx, user)
	})
}

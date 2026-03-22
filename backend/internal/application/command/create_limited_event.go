package command

import (
	"context"
	"fmt"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type CreateLimitedEventCommand struct {
	TenantID    string
	Title       string
	Description string
	Icon        string
	RarityBoost float64
	StartsAt    time.Time
	EndsAt      time.Time
}

type CreateLimitedEventHandler struct {
	repo repository.LimitedEventRepository
}

func NewCreateLimitedEventHandler(repo repository.LimitedEventRepository) *CreateLimitedEventHandler {
	return &CreateLimitedEventHandler{repo: repo}
}

func (h *CreateLimitedEventHandler) Handle(ctx context.Context, cmd *CreateLimitedEventCommand) (string, error) {
	tenantID, err := uuid.Parse(cmd.TenantID)
	if err != nil {
		return "", fmt.Errorf("invalid tenant ID: %w", err)
	}

	if cmd.RarityBoost <= 0 {
		cmd.RarityBoost = 1.5
	}

	event := entity.NewLimitedEvent(
		tenantID,
		cmd.Title,
		cmd.Description,
		cmd.Icon,
		cmd.RarityBoost,
		cmd.StartsAt,
		cmd.EndsAt,
	)

	if err := h.repo.Save(ctx, event); err != nil {
		return "", err
	}

	return event.ID.String(), nil
}

package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type CreateTagCommand struct {
	TenantID string
	Name     string
}

type CreateTagResult struct {
	TagID string
}

type CreateTagHandler struct {
	uow     uow.UnitOfWork
	tagRepo repository.TagRepository
}

func NewCreateTagHandler(uow uow.UnitOfWork, tagRepo repository.TagRepository) *CreateTagHandler {
	return &CreateTagHandler{uow: uow, tagRepo: tagRepo}
}

func (h *CreateTagHandler) Handle(ctx context.Context, cmd *CreateTagCommand) (*CreateTagResult, error) {
	tenantID, err := uuid.Parse(cmd.TenantID)
	if err != nil {
		return nil, err
	}

	var result *CreateTagResult
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		tag, err := entity.NewTag(tenantID, cmd.Name)
		if err != nil {
			return err
		}
		if err := h.tagRepo.Save(ctx, tag); err != nil {
			return err
		}
		result = &CreateTagResult{TagID: tag.ID.String()}
		return nil
	})
	return result, err
}

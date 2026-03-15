package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type DeleteTagCommand struct {
	TagID string
}

type DeleteTagHandler struct {
	uow     uow.UnitOfWork
	tagRepo repository.TagRepository
}

func NewDeleteTagHandler(uow uow.UnitOfWork, tagRepo repository.TagRepository) *DeleteTagHandler {
	return &DeleteTagHandler{uow: uow, tagRepo: tagRepo}
}

func (h *DeleteTagHandler) Handle(ctx context.Context, cmd *DeleteTagCommand) error {
	tagID, err := uuid.Parse(cmd.TagID)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(ctx context.Context) error {
		return h.tagRepo.Delete(ctx, tagID)
	})
}

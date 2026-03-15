package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type DeleteMemoCommand struct {
	MemoID string
}

type DeleteMemoHandler struct {
	uow      uow.UnitOfWork
	memoRepo repository.MemoRepository
}

func NewDeleteMemoHandler(uow uow.UnitOfWork, memoRepo repository.MemoRepository) *DeleteMemoHandler {
	return &DeleteMemoHandler{uow: uow, memoRepo: memoRepo}
}

func (h *DeleteMemoHandler) Handle(ctx context.Context, cmd *DeleteMemoCommand) error {
	memoID, err := uuid.Parse(cmd.MemoID)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(ctx context.Context) error {
		return h.memoRepo.Delete(ctx, memoID)
	})
}

package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type UpdateMemoCommand struct {
	MemoID  string
	Content string
}

type UpdateMemoHandler struct {
	uow      uow.UnitOfWork
	memoRepo repository.MemoRepository
}

func NewUpdateMemoHandler(uow uow.UnitOfWork, memoRepo repository.MemoRepository) *UpdateMemoHandler {
	return &UpdateMemoHandler{uow: uow, memoRepo: memoRepo}
}

func (h *UpdateMemoHandler) Handle(ctx context.Context, cmd *UpdateMemoCommand) error {
	memoID, err := uuid.Parse(cmd.MemoID)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(ctx context.Context) error {
		memo, err := h.memoRepo.FindByID(ctx, memoID)
		if err != nil {
			return err
		}
		if err := memo.UpdateContent(cmd.Content); err != nil {
			return err
		}
		return h.memoRepo.Save(ctx, memo)
	})
}

package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type AddMemoCommand struct {
	BookID  string
	Content string
}

type AddMemoResult struct {
	MemoID string
}

type AddMemoHandler struct {
	uow      uow.UnitOfWork
	bookRepo repository.BookRepository
	memoRepo repository.MemoRepository
}

func NewAddMemoHandler(uow uow.UnitOfWork, bookRepo repository.BookRepository, memoRepo repository.MemoRepository) *AddMemoHandler {
	return &AddMemoHandler{uow: uow, bookRepo: bookRepo, memoRepo: memoRepo}
}

func (h *AddMemoHandler) Handle(ctx context.Context, cmd *AddMemoCommand) (*AddMemoResult, error) {
	bookID, err := uuid.Parse(cmd.BookID)
	if err != nil {
		return nil, err
	}

	var result *AddMemoResult
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		// Verify book exists (lock order: Book=3)
		if _, err := h.bookRepo.FindByID(ctx, bookID); err != nil {
			return err
		}

		// Create memo (lock order: Memo=4, after Book=3 ✓)
		memo, err := entity.NewMemo(bookID, cmd.Content)
		if err != nil {
			return err
		}
		if err := h.memoRepo.Save(ctx, memo); err != nil {
			return err
		}

		result = &AddMemoResult{MemoID: memo.ID.String()}
		return nil
	})
	return result, err
}

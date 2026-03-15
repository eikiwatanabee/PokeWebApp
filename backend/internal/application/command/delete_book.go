package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type DeleteBookCommand struct {
	BookID string
}

type DeleteBookHandler struct {
	uow      uow.UnitOfWork
	bookRepo repository.BookRepository
}

func NewDeleteBookHandler(uow uow.UnitOfWork, bookRepo repository.BookRepository) *DeleteBookHandler {
	return &DeleteBookHandler{uow: uow, bookRepo: bookRepo}
}

func (h *DeleteBookHandler) Handle(ctx context.Context, cmd *DeleteBookCommand) error {
	bookID, err := uuid.Parse(cmd.BookID)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(ctx context.Context) error {
		return h.bookRepo.Delete(ctx, bookID)
	})
}

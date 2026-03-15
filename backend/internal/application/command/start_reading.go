package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type StartReadingCommand struct {
	BookID string
}

type StartReadingHandler struct {
	uow      uow.UnitOfWork
	bookRepo repository.BookRepository
}

func NewStartReadingHandler(uow uow.UnitOfWork, bookRepo repository.BookRepository) *StartReadingHandler {
	return &StartReadingHandler{uow: uow, bookRepo: bookRepo}
}

func (h *StartReadingHandler) Handle(ctx context.Context, cmd *StartReadingCommand) error {
	bookID, err := uuid.Parse(cmd.BookID)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(ctx context.Context) error {
		book, err := h.bookRepo.FindByIDForUpdate(ctx, bookID)
		if err != nil {
			return err
		}
		if err := book.StartReading(); err != nil {
			return err
		}
		return h.bookRepo.Save(ctx, book)
	})
}

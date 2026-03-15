package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type UpdateBookCommand struct {
	BookID string
	Title  string
	Author string
	TagIDs []string
}

type UpdateBookHandler struct {
	uow      uow.UnitOfWork
	bookRepo repository.BookRepository
	tagRepo  repository.TagRepository
}

func NewUpdateBookHandler(uow uow.UnitOfWork, bookRepo repository.BookRepository, tagRepo repository.TagRepository) *UpdateBookHandler {
	return &UpdateBookHandler{uow: uow, bookRepo: bookRepo, tagRepo: tagRepo}
}

func (h *UpdateBookHandler) Handle(ctx context.Context, cmd *UpdateBookCommand) error {
	bookID, err := uuid.Parse(cmd.BookID)
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(ctx context.Context) error {
		// Lock order: Book(3)
		book, err := h.bookRepo.FindByIDForUpdate(ctx, bookID)
		if err != nil {
			return err
		}
		if cmd.Title != "" {
			book.Title = cmd.Title
		}
		if cmd.Author != "" {
			book.Author = cmd.Author
		}
		// Lock order: Tag(5), after Book(3) ✓
		if cmd.TagIDs != nil {
			book.Tags = nil
			for _, tagIDStr := range cmd.TagIDs {
				tagID, err := uuid.Parse(tagIDStr)
				if err != nil {
					return err
				}
				tag, err := h.tagRepo.FindByID(ctx, tagID)
				if err != nil {
					return err
				}
				book.Tags = append(book.Tags, *tag)
			}
		}
		return h.bookRepo.Save(ctx, book)
	})
}

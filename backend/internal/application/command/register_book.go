package command

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/uow"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type RegisterBookCommand struct {
	UserID string
	Title  string
	Author string
	TagIDs []string
}

type RegisterBookResult struct {
	BookID string
}

type RegisterBookHandler struct {
	uow      uow.UnitOfWork
	bookRepo repository.BookRepository
	tagRepo  repository.TagRepository
}

func NewRegisterBookHandler(uow uow.UnitOfWork, bookRepo repository.BookRepository, tagRepo repository.TagRepository) *RegisterBookHandler {
	return &RegisterBookHandler{uow: uow, bookRepo: bookRepo, tagRepo: tagRepo}
}

func (h *RegisterBookHandler) Handle(ctx context.Context, cmd *RegisterBookCommand) (*RegisterBookResult, error) {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}

	var result *RegisterBookResult
	err = h.uow.Do(ctx, func(ctx context.Context) error {
		book, err := entity.NewBook(userID, cmd.Title, cmd.Author)
		if err != nil {
			return err
		}

		// Attach tags (lock order: Book=3, then Tag=5 ✓)
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

		if err := h.bookRepo.Save(ctx, book); err != nil {
			return err
		}

		result = &RegisterBookResult{BookID: book.ID.String()}
		return nil
	})
	return result, err
}

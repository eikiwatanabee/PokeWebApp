package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestDeleteBookHandler_Handle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		bookRepo.books[book.ID] = book

		handler := NewDeleteBookHandler(&mockUoW{}, bookRepo)
		err := handler.Handle(context.Background(), &DeleteBookCommand{BookID: book.ID.String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if len(bookRepo.books) != 0 {
			t.Error("book should be deleted")
		}
	})

	t.Run("invalid book ID", func(t *testing.T) {
		handler := NewDeleteBookHandler(&mockUoW{}, newMockBookRepo())
		err := handler.Handle(context.Background(), &DeleteBookCommand{BookID: "invalid"})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

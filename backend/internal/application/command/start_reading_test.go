package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

func TestStartReadingHandler_Handle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		bookRepo.books[book.ID] = book

		handler := NewStartReadingHandler(&mockUoW{}, bookRepo)
		err := handler.Handle(context.Background(), &StartReadingCommand{
			BookID: book.ID.String(),
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if bookRepo.books[book.ID].Status != valueobject.Reading {
			t.Errorf("expected Reading status, got %v", bookRepo.books[book.ID].Status)
		}
	})

	t.Run("book already reading", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		book.Status = valueobject.Reading
		bookRepo.books[book.ID] = book

		handler := NewStartReadingHandler(&mockUoW{}, bookRepo)
		err := handler.Handle(context.Background(), &StartReadingCommand{
			BookID: book.ID.String(),
		})
		if err == nil {
			t.Error("Handle() should return error for already reading book")
		}
	})

	t.Run("invalid book ID", func(t *testing.T) {
		handler := NewStartReadingHandler(&mockUoW{}, newMockBookRepo())
		err := handler.Handle(context.Background(), &StartReadingCommand{
			BookID: "invalid",
		})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

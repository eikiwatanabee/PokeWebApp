package query

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestGetBooksHandler_Handle(t *testing.T) {
	t.Run("returns books for user", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		userID := uuid.New()
		book1, _ := entity.NewBook(userID, "DDD", "Eric Evans")
		book2, _ := entity.NewBook(userID, "Clean Code", "Robert Martin")
		bookRepo.books[book1.ID] = book1
		bookRepo.books[book2.ID] = book2

		handler := NewGetBooksHandler(bookRepo)
		result, err := handler.Handle(context.Background(), &GetBooksQuery{
			UserID: userID.String(),
			Page:   1,
			Limit:  20,
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if len(result.Books) != 2 {
			t.Errorf("expected 2 books, got %d", len(result.Books))
		}
		if result.TotalCount != 2 {
			t.Errorf("TotalCount = %v, want 2", result.TotalCount)
		}
	})

	t.Run("defaults page and limit", func(t *testing.T) {
		handler := NewGetBooksHandler(newMockBookRepo())
		result, err := handler.Handle(context.Background(), &GetBooksQuery{
			UserID: uuid.New().String(),
			Page:   0,
			Limit:  0,
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.Page != 1 {
			t.Errorf("Page = %v, want 1", result.Page)
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		handler := NewGetBooksHandler(newMockBookRepo())
		_, err := handler.Handle(context.Background(), &GetBooksQuery{UserID: "invalid"})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

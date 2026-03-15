package query

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestGetBookDetailHandler_Handle(t *testing.T) {
	t.Run("returns book with memos", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		memoRepo := newMockMemoRepo()

		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		bookRepo.books[book.ID] = book

		memo, _ := entity.NewMemo(book.ID, "Great book")
		memoRepo.memos[memo.ID] = memo

		handler := NewGetBookDetailHandler(bookRepo, memoRepo)
		result, err := handler.Handle(context.Background(), &GetBookDetailQuery{BookID: book.ID.String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.Title != "DDD" {
			t.Errorf("Title = %v, want DDD", result.Title)
		}
		if result.Status != "unread" {
			t.Errorf("Status = %v, want unread", result.Status)
		}
		if len(result.Memos) != 1 {
			t.Errorf("expected 1 memo, got %d", len(result.Memos))
		}
		if result.FinishedAt != nil {
			t.Error("FinishedAt should be nil for unread book")
		}
	})

	t.Run("book not found", func(t *testing.T) {
		handler := NewGetBookDetailHandler(newMockBookRepo(), newMockMemoRepo())
		_, err := handler.Handle(context.Background(), &GetBookDetailQuery{BookID: uuid.New().String()})
		if err == nil {
			t.Error("Handle() should return error when book not found")
		}
	})
}

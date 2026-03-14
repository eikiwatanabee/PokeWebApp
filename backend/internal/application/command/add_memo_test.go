package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestAddMemoHandler_Handle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		memoRepo := newMockMemoRepo()

		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		bookRepo.books[book.ID] = book

		handler := NewAddMemoHandler(&mockUoW{}, bookRepo, memoRepo)
		result, err := handler.Handle(context.Background(), &AddMemoCommand{
			BookID:  book.ID.String(),
			Content: "Great chapter on aggregates",
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.MemoID == "" {
			t.Error("MemoID should not be empty")
		}
		if len(memoRepo.memos) != 1 {
			t.Errorf("expected 1 memo saved, got %d", len(memoRepo.memos))
		}
	})

	t.Run("empty content", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		bookRepo.books[book.ID] = book

		handler := NewAddMemoHandler(&mockUoW{}, bookRepo, newMockMemoRepo())
		_, err := handler.Handle(context.Background(), &AddMemoCommand{
			BookID:  book.ID.String(),
			Content: "",
		})
		if err != entity.ErrEmptyMemoContent {
			t.Errorf("Handle() error = %v, want %v", err, entity.ErrEmptyMemoContent)
		}
	})

	t.Run("book not found", func(t *testing.T) {
		handler := NewAddMemoHandler(&mockUoW{}, newMockBookRepo(), newMockMemoRepo())
		_, err := handler.Handle(context.Background(), &AddMemoCommand{
			BookID:  uuid.New().String(),
			Content: "content",
		})
		if err == nil {
			t.Error("Handle() should return error when book not found")
		}
	})
}

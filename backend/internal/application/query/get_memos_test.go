package query

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestGetMemosHandler_Handle(t *testing.T) {
	t.Run("returns memos for book", func(t *testing.T) {
		memoRepo := newMockMemoRepo()
		bookID := uuid.New()
		memo1, _ := entity.NewMemo(bookID, "First memo")
		memo2, _ := entity.NewMemo(bookID, "Second memo")
		memoRepo.memos[memo1.ID] = memo1
		memoRepo.memos[memo2.ID] = memo2

		handler := NewGetMemosHandler(memoRepo)
		result, err := handler.Handle(context.Background(), &GetMemosQuery{BookID: bookID.String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if len(result.Memos) != 2 {
			t.Errorf("expected 2 memos, got %d", len(result.Memos))
		}
	})

	t.Run("empty memos", func(t *testing.T) {
		handler := NewGetMemosHandler(newMockMemoRepo())
		result, err := handler.Handle(context.Background(), &GetMemosQuery{BookID: uuid.New().String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if len(result.Memos) != 0 {
			t.Errorf("expected 0 memos, got %d", len(result.Memos))
		}
	})

	t.Run("invalid book ID", func(t *testing.T) {
		handler := NewGetMemosHandler(newMockMemoRepo())
		_, err := handler.Handle(context.Background(), &GetMemosQuery{BookID: "invalid"})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

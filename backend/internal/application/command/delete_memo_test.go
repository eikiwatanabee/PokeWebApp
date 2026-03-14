package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestDeleteMemoHandler_Handle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		memoRepo := newMockMemoRepo()
		memo, _ := entity.NewMemo(uuid.New(), "content")
		memoRepo.memos[memo.ID] = memo

		handler := NewDeleteMemoHandler(&mockUoW{}, memoRepo)
		err := handler.Handle(context.Background(), &DeleteMemoCommand{MemoID: memo.ID.String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if len(memoRepo.memos) != 0 {
			t.Error("memo should be deleted")
		}
	})

	t.Run("invalid memo ID", func(t *testing.T) {
		handler := NewDeleteMemoHandler(&mockUoW{}, newMockMemoRepo())
		err := handler.Handle(context.Background(), &DeleteMemoCommand{MemoID: "invalid"})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

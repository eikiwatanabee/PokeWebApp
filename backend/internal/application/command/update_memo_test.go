package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestUpdateMemoHandler_Handle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		memoRepo := newMockMemoRepo()
		memo, _ := entity.NewMemo(uuid.New(), "original content")
		memoRepo.memos[memo.ID] = memo

		handler := NewUpdateMemoHandler(&mockUoW{}, memoRepo)
		err := handler.Handle(context.Background(), &UpdateMemoCommand{
			MemoID:  memo.ID.String(),
			Content: "updated content",
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if memoRepo.memos[memo.ID].Content != "updated content" {
			t.Errorf("Content = %v, want updated content", memoRepo.memos[memo.ID].Content)
		}
	})

	t.Run("empty content", func(t *testing.T) {
		memoRepo := newMockMemoRepo()
		memo, _ := entity.NewMemo(uuid.New(), "original")
		memoRepo.memos[memo.ID] = memo

		handler := NewUpdateMemoHandler(&mockUoW{}, memoRepo)
		err := handler.Handle(context.Background(), &UpdateMemoCommand{
			MemoID:  memo.ID.String(),
			Content: "",
		})
		if err != entity.ErrEmptyMemoContent {
			t.Errorf("Handle() error = %v, want %v", err, entity.ErrEmptyMemoContent)
		}
	})
}

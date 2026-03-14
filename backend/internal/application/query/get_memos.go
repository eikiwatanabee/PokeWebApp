package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetMemosQuery struct {
	BookID string
}

type GetMemosResult struct {
	Memos []MemoDTO `json:"memos"`
}

type GetMemosHandler struct {
	memoRepo repository.MemoRepository
}

func NewGetMemosHandler(memoRepo repository.MemoRepository) *GetMemosHandler {
	return &GetMemosHandler{memoRepo: memoRepo}
}

func (h *GetMemosHandler) Handle(ctx context.Context, q *GetMemosQuery) (*GetMemosResult, error) {
	bookID, err := uuid.Parse(q.BookID)
	if err != nil {
		return nil, err
	}

	memos, err := h.memoRepo.FindByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	dtos := make([]MemoDTO, len(memos))
	for i, m := range memos {
		dtos[i] = MemoDTO{
			ID:        m.ID.String(),
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &GetMemosResult{Memos: dtos}, nil
}

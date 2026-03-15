package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetBooksQuery struct {
	UserID string
	Status string
	TagID  string
	Page   int
	Limit  int
}

type BookDTO struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Status   string   `json:"status"`
	Tags     []string `json:"tags"`
}

type GetBooksResult struct {
	Books      []BookDTO `json:"books"`
	TotalCount int64     `json:"total_count"`
	Page       int       `json:"page"`
}

type GetBooksHandler struct {
	bookRepo repository.BookRepository
}

func NewGetBooksHandler(bookRepo repository.BookRepository) *GetBooksHandler {
	return &GetBooksHandler{bookRepo: bookRepo}
}

func (h *GetBooksHandler) Handle(ctx context.Context, q *GetBooksQuery) (*GetBooksResult, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return nil, err
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}

	books, total, err := h.bookRepo.FindByUserID(ctx, userID, q.Status, q.TagID, q.Page, q.Limit)
	if err != nil {
		return nil, err
	}

	dtos := make([]BookDTO, len(books))
	for i, b := range books {
		tagNames := make([]string, len(b.Tags))
		for j, t := range b.Tags {
			tagNames[j] = t.Name
		}
		dtos[i] = BookDTO{
			ID:     b.ID.String(),
			Title:  b.Title,
			Author: b.Author,
			Status: b.Status.String(),
			Tags:   tagNames,
		}
	}

	return &GetBooksResult{
		Books:      dtos,
		TotalCount: total,
		Page:       q.Page,
	}, nil
}

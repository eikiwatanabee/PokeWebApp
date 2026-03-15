package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetBookDetailQuery struct {
	BookID string
}

type BookDetailDTO struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Status     string    `json:"status"`
	FinishedAt *string   `json:"finished_at"`
	Tags       []TagDTO  `json:"tags"`
	Memos      []MemoDTO `json:"memos"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

type MemoDTO struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TagDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetBookDetailHandler struct {
	bookRepo repository.BookRepository
	memoRepo repository.MemoRepository
}

func NewGetBookDetailHandler(bookRepo repository.BookRepository, memoRepo repository.MemoRepository) *GetBookDetailHandler {
	return &GetBookDetailHandler{bookRepo: bookRepo, memoRepo: memoRepo}
}

func (h *GetBookDetailHandler) Handle(ctx context.Context, q *GetBookDetailQuery) (*BookDetailDTO, error) {
	bookID, err := uuid.Parse(q.BookID)
	if err != nil {
		return nil, err
	}

	book, err := h.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	memos, err := h.memoRepo.FindByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	tags := make([]TagDTO, len(book.Tags))
	for i, t := range book.Tags {
		tags[i] = TagDTO{ID: t.ID.String(), Name: t.Name}
	}

	memoDTOs := make([]MemoDTO, len(memos))
	for i, m := range memos {
		memoDTOs[i] = MemoDTO{
			ID:        m.ID.String(),
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	var finishedAt *string
	if book.FinishedAt != nil {
		s := book.FinishedAt.Format("2006-01-02T15:04:05Z")
		finishedAt = &s
	}

	return &BookDetailDTO{
		ID:         book.ID.String(),
		Title:      book.Title,
		Author:     book.Author,
		Status:     book.Status.String(),
		FinishedAt: finishedAt,
		Tags:       tags,
		Memos:      memoDTOs,
		CreatedAt:  book.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  book.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

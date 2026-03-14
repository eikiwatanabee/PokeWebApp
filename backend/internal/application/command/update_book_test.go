package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestUpdateBookHandler_Handle(t *testing.T) {
	t.Run("update title and author", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		tagRepo := newMockTagRepo()
		book, _ := entity.NewBook(uuid.New(), "Old Title", "Old Author")
		bookRepo.books[book.ID] = book

		handler := NewUpdateBookHandler(&mockUoW{}, bookRepo, tagRepo)
		err := handler.Handle(context.Background(), &UpdateBookCommand{
			BookID: book.ID.String(),
			Title:  "New Title",
			Author: "New Author",
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if bookRepo.books[book.ID].Title != "New Title" {
			t.Errorf("Title = %v, want New Title", bookRepo.books[book.ID].Title)
		}
		if bookRepo.books[book.ID].Author != "New Author" {
			t.Errorf("Author = %v, want New Author", bookRepo.books[book.ID].Author)
		}
	})

	t.Run("update tags", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		tagRepo := newMockTagRepo()

		tenantID := uuid.New()
		tag, _ := entity.NewTag(tenantID, "技術書")
		tagRepo.tags[tag.ID] = tag

		book, _ := entity.NewBook(uuid.New(), "DDD", "Eric Evans")
		bookRepo.books[book.ID] = book

		handler := NewUpdateBookHandler(&mockUoW{}, bookRepo, tagRepo)
		err := handler.Handle(context.Background(), &UpdateBookCommand{
			BookID: book.ID.String(),
			TagIDs: []string{tag.ID.String()},
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if len(bookRepo.books[book.ID].Tags) != 1 {
			t.Errorf("expected 1 tag, got %d", len(bookRepo.books[book.ID].Tags))
		}
	})
}

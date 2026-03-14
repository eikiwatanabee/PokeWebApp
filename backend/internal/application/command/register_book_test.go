package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestRegisterBookHandler_Handle(t *testing.T) {
	t.Run("success without tags", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		tagRepo := newMockTagRepo()
		handler := NewRegisterBookHandler(&mockUoW{}, bookRepo, tagRepo)

		result, err := handler.Handle(context.Background(), &RegisterBookCommand{
			UserID: uuid.New().String(),
			Title:  "Clean Architecture",
			Author: "Robert C. Martin",
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.BookID == "" {
			t.Error("Handle() BookID should not be empty")
		}
		if len(bookRepo.books) != 1 {
			t.Errorf("expected 1 book saved, got %d", len(bookRepo.books))
		}
	})

	t.Run("success with tags", func(t *testing.T) {
		bookRepo := newMockBookRepo()
		tagRepo := newMockTagRepo()

		tenantID := uuid.New()
		tag, _ := entity.NewTag(tenantID, "技術書")
		tagRepo.tags[tag.ID] = tag

		handler := NewRegisterBookHandler(&mockUoW{}, bookRepo, tagRepo)

		result, err := handler.Handle(context.Background(), &RegisterBookCommand{
			UserID: uuid.New().String(),
			Title:  "DDD",
			Author: "Eric Evans",
			TagIDs: []string{tag.ID.String()},
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.BookID == "" {
			t.Error("Handle() BookID should not be empty")
		}

		bookID, _ := uuid.Parse(result.BookID)
		savedBook := bookRepo.books[bookID]
		if len(savedBook.Tags) != 1 {
			t.Errorf("expected 1 tag, got %d", len(savedBook.Tags))
		}
	})

	t.Run("empty title", func(t *testing.T) {
		handler := NewRegisterBookHandler(&mockUoW{}, newMockBookRepo(), newMockTagRepo())
		_, err := handler.Handle(context.Background(), &RegisterBookCommand{
			UserID: uuid.New().String(),
			Title:  "",
			Author: "Author",
		})
		if err != entity.ErrEmptyTitle {
			t.Errorf("Handle() error = %v, want %v", err, entity.ErrEmptyTitle)
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		handler := NewRegisterBookHandler(&mockUoW{}, newMockBookRepo(), newMockTagRepo())
		_, err := handler.Handle(context.Background(), &RegisterBookCommand{
			UserID: "invalid",
			Title:  "Book",
			Author: "Author",
		})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

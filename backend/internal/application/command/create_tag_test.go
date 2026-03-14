package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestCreateTagHandler_Handle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tagRepo := newMockTagRepo()
		handler := NewCreateTagHandler(&mockUoW{}, tagRepo)

		result, err := handler.Handle(context.Background(), &CreateTagCommand{
			TenantID: uuid.New().String(),
			Name:     "技術書",
		})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if result.TagID == "" {
			t.Error("TagID should not be empty")
		}
		if len(tagRepo.tags) != 1 {
			t.Errorf("expected 1 tag saved, got %d", len(tagRepo.tags))
		}
	})

	t.Run("empty name", func(t *testing.T) {
		handler := NewCreateTagHandler(&mockUoW{}, newMockTagRepo())
		_, err := handler.Handle(context.Background(), &CreateTagCommand{
			TenantID: uuid.New().String(),
			Name:     "",
		})
		if err != entity.ErrEmptyTagName {
			t.Errorf("Handle() error = %v, want %v", err, entity.ErrEmptyTagName)
		}
	})

	t.Run("invalid tenant ID", func(t *testing.T) {
		handler := NewCreateTagHandler(&mockUoW{}, newMockTagRepo())
		_, err := handler.Handle(context.Background(), &CreateTagCommand{
			TenantID: "invalid",
			Name:     "技術書",
		})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

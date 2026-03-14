package command

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestDeleteTagHandler_Handle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tagRepo := newMockTagRepo()
		tag, _ := entity.NewTag(uuid.New(), "技術書")
		tagRepo.tags[tag.ID] = tag

		handler := NewDeleteTagHandler(&mockUoW{}, tagRepo)
		err := handler.Handle(context.Background(), &DeleteTagCommand{TagID: tag.ID.String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if len(tagRepo.tags) != 0 {
			t.Error("tag should be deleted")
		}
	})

	t.Run("invalid tag ID", func(t *testing.T) {
		handler := NewDeleteTagHandler(&mockUoW{}, newMockTagRepo())
		err := handler.Handle(context.Background(), &DeleteTagCommand{TagID: "invalid"})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

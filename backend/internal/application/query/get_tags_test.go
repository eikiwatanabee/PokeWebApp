package query

import (
	"context"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

func TestGetTagsHandler_Handle(t *testing.T) {
	t.Run("returns tags for tenant", func(t *testing.T) {
		tagRepo := newMockTagRepo()
		tenantID := uuid.New()
		tag1, _ := entity.NewTag(tenantID, "技術書")
		tag2, _ := entity.NewTag(tenantID, "ビジネス書")
		tagRepo.tags[tag1.ID] = tag1
		tagRepo.tags[tag2.ID] = tag2

		handler := NewGetTagsHandler(tagRepo)
		result, err := handler.Handle(context.Background(), &GetTagsQuery{TenantID: tenantID.String()})
		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if len(result.Tags) != 2 {
			t.Errorf("expected 2 tags, got %d", len(result.Tags))
		}
	})

	t.Run("invalid tenant ID", func(t *testing.T) {
		handler := NewGetTagsHandler(newMockTagRepo())
		_, err := handler.Handle(context.Background(), &GetTagsQuery{TenantID: "invalid"})
		if err == nil {
			t.Error("Handle() should return error for invalid UUID")
		}
	})
}

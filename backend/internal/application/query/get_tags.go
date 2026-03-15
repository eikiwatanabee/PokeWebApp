package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetTagsQuery struct {
	TenantID string
}

type GetTagsResult struct {
	Tags []TagDTO `json:"tags"`
}

type GetTagsHandler struct {
	tagRepo repository.TagRepository
}

func NewGetTagsHandler(tagRepo repository.TagRepository) *GetTagsHandler {
	return &GetTagsHandler{tagRepo: tagRepo}
}

func (h *GetTagsHandler) Handle(ctx context.Context, q *GetTagsQuery) (*GetTagsResult, error) {
	tenantID, err := uuid.Parse(q.TenantID)
	if err != nil {
		return nil, err
	}

	tags, err := h.tagRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	dtos := make([]TagDTO, len(tags))
	for i, t := range tags {
		dtos[i] = TagDTO{ID: t.ID.String(), Name: t.Name}
	}

	return &GetTagsResult{Tags: dtos}, nil
}

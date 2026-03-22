package query

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetLimitedEventsQuery struct {
	TenantID string
}

type LimitedEventDTO struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	RarityBoost float64 `json:"rarity_boost"`
	StartsAt    string  `json:"starts_at"`
	EndsAt      string  `json:"ends_at"`
	Active      bool    `json:"active"`
	HoursLeft   int     `json:"hours_left"`
}

type GetLimitedEventsResult struct {
	Events []LimitedEventDTO `json:"events"`
}

type GetLimitedEventsHandler struct {
	limitedEventRepo repository.LimitedEventRepository
}

func NewGetLimitedEventsHandler(limitedEventRepo repository.LimitedEventRepository) *GetLimitedEventsHandler {
	return &GetLimitedEventsHandler{limitedEventRepo: limitedEventRepo}
}

func (h *GetLimitedEventsHandler) Handle(ctx context.Context, q *GetLimitedEventsQuery) (*GetLimitedEventsResult, error) {
	tenantID, err := uuid.Parse(q.TenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	events, err := h.limitedEventRepo.FindActiveByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dtos := make([]LimitedEventDTO, len(events))
	for i, e := range events {
		hoursLeft := int(math.Ceil(e.EndsAt.Sub(now).Hours()))
		if hoursLeft < 0 {
			hoursLeft = 0
		}
		dtos[i] = LimitedEventDTO{
			ID:          e.ID.String(),
			Title:       e.Title,
			Description: e.Description,
			Icon:        e.Icon,
			RarityBoost: e.RarityBoost,
			StartsAt:    e.StartsAt.Format(time.RFC3339),
			EndsAt:      e.EndsAt.Format(time.RFC3339),
			Active:      e.IsActive(),
			HoursLeft:   hoursLeft,
		}
	}

	return &GetLimitedEventsResult{Events: dtos}, nil
}

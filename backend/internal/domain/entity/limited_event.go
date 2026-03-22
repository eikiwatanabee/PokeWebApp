package entity

import (
	"time"

	"github.com/google/uuid"
)

type LimitedEvent struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	Title           string
	Description     string
	Icon            string
	RarityBoost     float64 // Multiplier for rarity upgrade chance (e.g., 2.0 = double chance)
	StartsAt        time.Time
	EndsAt          time.Time
	CreatedAt       time.Time
}

func NewLimitedEvent(tenantID uuid.UUID, title, description, icon string, rarityBoost float64, startsAt, endsAt time.Time) *LimitedEvent {
	return &LimitedEvent{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Title:       title,
		Description: description,
		Icon:        icon,
		RarityBoost: rarityBoost,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		CreatedAt:   time.Now(),
	}
}

func (e *LimitedEvent) IsActive() bool {
	now := time.Now()
	return now.After(e.StartsAt) && now.Before(e.EndsAt)
}

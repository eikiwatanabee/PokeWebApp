package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Name      string
	CreatedAt time.Time
}

func NewTag(tenantID uuid.UUID, name string) (*Tag, error) {
	if name == "" {
		return nil, ErrEmptyTagName
	}
	return &Tag{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}

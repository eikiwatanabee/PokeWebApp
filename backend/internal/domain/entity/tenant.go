package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTenant(name string) (*Tenant, error) {
	if name == "" {
		return nil, ErrEmptyTenantName
	}
	now := time.Now()
	return &Tenant{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTeam(tenantID uuid.UUID, name string) (*Team, error) {
	if name == "" {
		return nil, ErrEmptyTeamName
	}
	now := time.Now()
	return &Team{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

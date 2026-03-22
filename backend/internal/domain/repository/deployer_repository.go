package repository

import (
	"context"

	"github.com/google/uuid"
)

type DeployerRepository interface {
	IsDeployer(ctx context.Context, tenantID, userID uuid.UUID) (bool, error)
	AddDeployer(ctx context.Context, tenantID, userID uuid.UUID) error
	RemoveDeployer(ctx context.Context, tenantID, userID uuid.UUID) error
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]uuid.UUID, error)
}

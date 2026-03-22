package repository

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type GitHubActivityRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.GitHubActivity, error)
	FindRecentByTenantID(ctx context.Context, tenantID uuid.UUID, limit int) ([]*entity.GitHubActivity, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CountByUserIDAndType(ctx context.Context, userID uuid.UUID, eventType entity.GitHubEventType) (int64, error)
	TotalXPByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	CountByUserIDSince(ctx context.Context, userID uuid.UUID, since time.Time) (int64, error)
	CountByUserIDAndTypeSince(ctx context.Context, userID uuid.UUID, eventType entity.GitHubEventType, since time.Time) (int64, error)
	TotalXPByUserIDSince(ctx context.Context, userID uuid.UUID, since time.Time) (int, error)
	Save(ctx context.Context, activity *entity.GitHubActivity) error
}

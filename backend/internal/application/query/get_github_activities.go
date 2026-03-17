package query

import (
	"context"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetGitHubActivitiesQuery struct {
	UserID string
	Limit  int
}

type GitHubActivityDTO struct {
	ID        string `json:"id"`
	EventType string `json:"event_type"`
	RepoName  string `json:"repo_name"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	XP        int    `json:"xp"`
	CreatedAt string `json:"created_at"`
}

type GetGitHubActivitiesResult struct {
	Activities []GitHubActivityDTO `json:"activities"`
	TotalCount int64               `json:"total_count"`
}

type GetGitHubActivitiesHandler struct {
	activityRepo repository.GitHubActivityRepository
}

func NewGetGitHubActivitiesHandler(activityRepo repository.GitHubActivityRepository) *GetGitHubActivitiesHandler {
	return &GetGitHubActivitiesHandler{activityRepo: activityRepo}
}

func (h *GetGitHubActivitiesHandler) Handle(ctx context.Context, q *GetGitHubActivitiesQuery) (*GetGitHubActivitiesResult, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return nil, err
	}

	limit := q.Limit
	if limit <= 0 {
		limit = 50
	}

	activities, err := h.activityRepo.FindByUserID(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	totalCount, err := h.activityRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]GitHubActivityDTO, len(activities))
	for i, a := range activities {
		dtos[i] = GitHubActivityDTO{
			ID:        a.ID.String(),
			EventType: string(a.EventType),
			RepoName:  a.RepoName,
			Title:     a.Title,
			URL:       a.URL,
			XP:        a.XP,
			CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &GetGitHubActivitiesResult{
		Activities: dtos,
		TotalCount: totalCount,
	}, nil
}

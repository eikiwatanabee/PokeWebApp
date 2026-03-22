package query

import (
	"context"
	"fmt"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetTeamFeedQuery struct {
	TenantID string
	Limit    int
}

type FeedItemDTO struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	UserAvatarURL  string `json:"user_avatar_url"`
	GitHubUsername string `json:"github_username"`
	EventType      string `json:"event_type"`
	RepoName       string `json:"repo_name"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	XP             int    `json:"xp"`
	CreatedAt      string `json:"created_at"`
}

type GetTeamFeedResult struct {
	Items      []FeedItemDTO `json:"items"`
	TotalCount int           `json:"total_count"`
}

type GetTeamFeedHandler struct {
	activityRepo repository.GitHubActivityRepository
	userRepo     repository.UserRepository
}

func NewGetTeamFeedHandler(
	activityRepo repository.GitHubActivityRepository,
	userRepo repository.UserRepository,
) *GetTeamFeedHandler {
	return &GetTeamFeedHandler{activityRepo: activityRepo, userRepo: userRepo}
}

func (h *GetTeamFeedHandler) Handle(ctx context.Context, q *GetTeamFeedQuery) (*GetTeamFeedResult, error) {
	tenantID, err := uuid.Parse(q.TenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	limit := q.Limit
	if limit <= 0 {
		limit = 50
	}

	activities, err := h.activityRepo.FindRecentByTenantID(ctx, tenantID, limit)
	if err != nil {
		return nil, err
	}

	// Collect unique user IDs and fetch users
	userIDs := make(map[uuid.UUID]bool)
	for _, a := range activities {
		userIDs[a.UserID] = true
	}

	users, err := h.userRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	userMap := make(map[uuid.UUID]struct {
		Name           string
		AvatarURL      string
		GitHubUsername string
	})
	for _, u := range users {
		userMap[u.ID] = struct {
			Name           string
			AvatarURL      string
			GitHubUsername string
		}{
			Name:           u.Name,
			AvatarURL:      u.AvatarURL,
			GitHubUsername: u.GitHubUsername,
		}
	}

	items := make([]FeedItemDTO, len(activities))
	for i, a := range activities {
		u := userMap[a.UserID]
		items[i] = FeedItemDTO{
			ID:             a.ID.String(),
			UserID:         a.UserID.String(),
			UserName:       u.Name,
			UserAvatarURL:  u.AvatarURL,
			GitHubUsername: u.GitHubUsername,
			EventType:      string(a.EventType),
			RepoName:       a.RepoName,
			Title:          a.Title,
			URL:            a.URL,
			XP:             a.XP,
			CreatedAt:      a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &GetTeamFeedResult{
		Items:      items,
		TotalCount: len(items),
	}, nil
}

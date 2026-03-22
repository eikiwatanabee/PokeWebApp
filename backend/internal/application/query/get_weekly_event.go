package query

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type GetWeeklyEventQuery struct {
	TenantID string
}

type WeeklyEventTeamScoreDTO struct {
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
	Score    int64  `json:"score"`
	Rank     int    `json:"rank"`
}

type GetWeeklyEventResult struct {
	EventType   string                    `json:"event_type"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Icon        string                    `json:"icon"`
	WeekStart   string                    `json:"week_start"`
	WeekEnd     string                    `json:"week_end"`
	TeamScores  []WeeklyEventTeamScoreDTO `json:"team_scores"`
}

type GetWeeklyEventHandler struct {
	weeklyEventRepo repository.WeeklyEventRepository
	teamRepo        repository.TeamRepository
	userRepo        repository.UserRepository
	activityRepo    repository.GitHubActivityRepository
}

func NewGetWeeklyEventHandler(
	weeklyEventRepo repository.WeeklyEventRepository,
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
	activityRepo repository.GitHubActivityRepository,
) *GetWeeklyEventHandler {
	return &GetWeeklyEventHandler{
		weeklyEventRepo: weeklyEventRepo,
		teamRepo:        teamRepo,
		userRepo:        userRepo,
		activityRepo:    activityRepo,
	}
}

func (h *GetWeeklyEventHandler) Handle(ctx context.Context, q *GetWeeklyEventQuery) (*GetWeeklyEventResult, error) {
	tenantID, err := uuid.Parse(q.TenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	weekStart := entity.CurrentWeekStart()

	// Get or create weekly event
	event, err := h.weeklyEventRepo.FindByTenantIDAndWeek(ctx, tenantID, weekStart)
	if err != nil {
		return nil, err
	}
	if event == nil {
		event = entity.GenerateWeeklyEvent(tenantID, weekStart)
		if err := h.weeklyEventRepo.Save(ctx, event); err != nil {
			return nil, err
		}
	}

	// Find event definition
	var def entity.WeeklyEventDefinition
	for _, d := range entity.WeeklyEventDefinitions {
		if d.Type == event.EventType {
			def = d
			break
		}
	}

	// Get teams and users
	teams, err := h.teamRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	users, err := h.userRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Build team member map
	teamMembers := make(map[uuid.UUID][]uuid.UUID)
	for _, u := range users {
		if u.TeamID != nil {
			teamMembers[*u.TeamID] = append(teamMembers[*u.TeamID], u.ID)
		}
	}

	// Calculate scores per team
	teamScores := make([]WeeklyEventTeamScoreDTO, 0, len(teams))
	for _, team := range teams {
		memberIDs := teamMembers[team.ID]
		var totalScore int64

		for _, memberID := range memberIDs {
			score, err := h.calculateScore(ctx, memberID, def, weekStart)
			if err != nil {
				return nil, err
			}
			totalScore += score
		}

		teamScores = append(teamScores, WeeklyEventTeamScoreDTO{
			TeamID:   team.ID.String(),
			TeamName: team.Name,
			Score:    totalScore,
		})
	}

	// Sort by score descending
	sort.Slice(teamScores, func(i, j int) bool {
		return teamScores[i].Score > teamScores[j].Score
	})

	for i := range teamScores {
		teamScores[i].Rank = i + 1
	}

	return &GetWeeklyEventResult{
		EventType:   string(event.EventType),
		Title:       def.Title,
		Description: def.Description,
		Icon:        def.Icon,
		WeekStart:   weekStart.Format("2006-01-02"),
		WeekEnd:     event.WeekEnd.Format("2006-01-02"),
		TeamScores:  teamScores,
	}, nil
}

func (h *GetWeeklyEventHandler) calculateScore(ctx context.Context, userID uuid.UUID, def entity.WeeklyEventDefinition, since time.Time) (int64, error) {
	if def.ScoreType == "" {
		// All-rounder: total XP since week start
		xp, err := h.activityRepo.TotalXPByUserIDSince(ctx, userID, since)
		return int64(xp), err
	}
	// Specific event type count
	return h.activityRepo.CountByUserIDAndTypeSince(ctx, userID, def.ScoreType, since)
}

package service

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type DailyMissionService struct {
	missionRepo  repository.DailyMissionRepository
	activityRepo repository.GitHubActivityRepository
}

func NewDailyMissionService(
	missionRepo repository.DailyMissionRepository,
	activityRepo repository.GitHubActivityRepository,
) *DailyMissionService {
	return &DailyMissionService{
		missionRepo:  missionRepo,
		activityRepo: activityRepo,
	}
}

// GetOrCreateToday returns today's missions, creating them if they don't exist yet.
func (s *DailyMissionService) GetOrCreateToday(ctx context.Context, userID uuid.UUID) ([]*entity.DailyMission, error) {
	today := time.Now().Truncate(24 * time.Hour)

	missions, err := s.missionRepo.FindByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return nil, err
	}

	if len(missions) > 0 {
		return missions, nil
	}

	// Generate new missions for today
	missions = entity.GenerateDailyMissions(userID, today)
	if err := s.missionRepo.SaveAll(ctx, missions); err != nil {
		return nil, err
	}

	return missions, nil
}

// ProcessActivity updates mission progress based on a new GitHub activity event.
// Returns total bonus XP earned (individual mission bonuses + all-complete bonus).
func (s *DailyMissionService) ProcessActivity(ctx context.Context, userID uuid.UUID, eventType entity.GitHubEventType) (int, error) {
	missions, err := s.GetOrCreateToday(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Build template lookup
	tmplMap := make(map[string]entity.MissionTemplate)
	for _, t := range entity.MissionTemplates {
		tmplMap[t.ID] = t
	}

	totalBonus := 0
	allDone := true

	for _, m := range missions {
		tmpl, ok := tmplMap[m.TemplateID]
		if !ok {
			continue
		}

		// Check if this event matches the mission
		matches := false
		if tmpl.EventType == "" {
			// "any" mission - matches all event types
			matches = true
		} else if tmpl.EventType == eventType {
			matches = true
		}

		if matches {
			if m.IncrementProgress(1) {
				totalBonus += m.BonusXP
			}
			if err := s.missionRepo.Save(ctx, m); err != nil {
				return 0, err
			}
		}

		if m.Status != entity.MissionCompleted {
			allDone = false
		}
	}

	// All-complete bonus
	if allDone && len(missions) == entity.DailyMissionSlotCount {
		totalBonus += entity.CompletionBonusXP
	}

	return totalBonus, nil
}

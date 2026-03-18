package query

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/service"
	"github.com/google/uuid"
)

type GetDailyMissionsQuery struct {
	UserID string
}

type DailyMissionDTO struct {
	ID          string `json:"id"`
	TemplateID  string `json:"template_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Progress    int    `json:"progress"`
	Required    int    `json:"required"`
	BonusXP     int    `json:"bonus_xp"`
	Status      string `json:"status"`
}

type GetDailyMissionsResult struct {
	Missions       []DailyMissionDTO `json:"missions"`
	AllCompleted   bool              `json:"all_completed"`
	CompletionBonus int             `json:"completion_bonus"`
	Date           string           `json:"date"`
}

type GetDailyMissionsHandler struct {
	missionSvc *service.DailyMissionService
}

func NewGetDailyMissionsHandler(missionSvc *service.DailyMissionService) *GetDailyMissionsHandler {
	return &GetDailyMissionsHandler{missionSvc: missionSvc}
}

func (h *GetDailyMissionsHandler) Handle(ctx context.Context, q *GetDailyMissionsQuery) (*GetDailyMissionsResult, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return nil, err
	}

	missions, err := h.missionSvc.GetOrCreateToday(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build template lookup
	tmplMap := make(map[string]entity.MissionTemplate)
	for _, t := range entity.MissionTemplates {
		tmplMap[t.ID] = t
	}

	dtos := make([]DailyMissionDTO, len(missions))
	allCompleted := true
	for i, m := range missions {
		tmpl := tmplMap[m.TemplateID]
		dtos[i] = DailyMissionDTO{
			ID:          m.ID.String(),
			TemplateID:  m.TemplateID,
			Title:       tmpl.Title,
			Description: tmpl.Description,
			Icon:        tmpl.Icon,
			Progress:    m.Progress,
			Required:    m.Required,
			BonusXP:     m.BonusXP,
			Status:      string(m.Status),
		}
		if m.Status != entity.MissionCompleted {
			allCompleted = false
		}
	}

	return &GetDailyMissionsResult{
		Missions:        dtos,
		AllCompleted:    allCompleted,
		CompletionBonus: entity.CompletionBonusXP,
		Date:            time.Now().Truncate(24 * time.Hour).Format("2006-01-02"),
	}, nil
}

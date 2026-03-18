package entity

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// MissionStatus represents the completion state of a daily mission.
type MissionStatus string

const (
	MissionPending   MissionStatus = "pending"
	MissionCompleted MissionStatus = "completed"
)

// MissionTemplate defines a type of mission that can be generated.
type MissionTemplate struct {
	ID          string
	Title       string
	Description string
	Icon        string
	EventType   GitHubEventType
	Required    int
	BonusXP     int
}

// MissionTemplates is the pool of possible daily missions.
var MissionTemplates = []MissionTemplate{
	// コミット系
	{"commit_1", "はじめのいっぽ", "今日コミットを1回しよう", "💻", EventCommit, 1, 20},
	{"commit_3", "コードラッシュ", "今日コミットを3回しよう", "💻", EventCommit, 3, 50},
	{"commit_5", "コードマラソン", "今日コミットを5回しよう", "💻", EventCommit, 5, 100},

	// PR系
	{"pr_open_1", "PRチャレンジ", "今日PRを1つ作ろう", "📝", EventPROpen, 1, 30},
	{"pr_merge_1", "マージタイム", "今日PRを1つマージしよう", "🎉", EventPRMerge, 1, 60},

	// レビュー系
	{"review_1", "チームワーク", "今日レビューを1回しよう", "👀", EventReview, 1, 30},
	{"review_2", "レビューマスター", "今日レビューを2回しよう", "👀", EventReview, 2, 60},
	{"review_3", "コードガーディアン", "今日レビューを3回しよう", "🛡️", EventReview, 3, 100},

	// Issue系
	{"issue_1", "バグバスター", "今日Issueを1つクローズしよう", "🐛", EventIssueClose, 1, 25},
	{"issue_2", "イシュースレイヤー", "今日Issueを2つクローズしよう", "⚔️", EventIssueClose, 2, 60},

	// ミックス系 (any activity)
	{"any_3", "アクティブトレーナー", "今日なんでも3アクション", "🔥", "", 3, 40},
	{"any_5", "スーパーアクティブ", "今日なんでも5アクション", "⚡", "", 5, 80},
}

// DailyMissionSlotCount is the number of missions generated per day.
const DailyMissionSlotCount = 3

// CompletionBonusXP is the bonus for completing all daily missions.
const CompletionBonusXP = 100

// DailyMission is a single mission instance assigned to a user for a specific date.
type DailyMission struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Date       time.Time // date only (truncated)
	TemplateID string
	Progress   int
	Required   int
	BonusXP    int
	Status     MissionStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// DailyMissionSet represents a user's missions for a specific day.
type DailyMissionSet struct {
	UserID         uuid.UUID
	Date           time.Time
	Missions       []*DailyMission
	AllCompleted   bool
	BonusClaimed   bool
}

// GenerateDailyMissions creates a set of random missions for a user on a given date.
func GenerateDailyMissions(userID uuid.UUID, date time.Time) []*DailyMission {
	truncated := date.Truncate(24 * time.Hour)
	now := time.Now()

	// Deterministic seed based on user+date for consistency
	seed := int64(userID[0])<<40 | int64(userID[1])<<32 |
		int64(truncated.Year())<<16 | int64(truncated.YearDay())
	r := rand.New(rand.NewSource(seed))

	// Shuffle and pick DailyMissionSlotCount templates
	indices := r.Perm(len(MissionTemplates))
	missions := make([]*DailyMission, 0, DailyMissionSlotCount)

	for i := 0; i < DailyMissionSlotCount && i < len(indices); i++ {
		tmpl := MissionTemplates[indices[i]]
		missions = append(missions, &DailyMission{
			ID:         uuid.New(),
			UserID:     userID,
			Date:       truncated,
			TemplateID: tmpl.ID,
			Progress:   0,
			Required:   tmpl.Required,
			BonusXP:    tmpl.BonusXP,
			Status:     MissionPending,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	return missions
}

// IncrementProgress adds progress to a mission and marks complete if threshold met.
// Returns true if the mission was newly completed.
func (m *DailyMission) IncrementProgress(amount int) bool {
	if m.Status == MissionCompleted {
		return false
	}
	m.Progress += amount
	m.UpdatedAt = time.Now()
	if m.Progress >= m.Required {
		m.Progress = m.Required
		m.Status = MissionCompleted
		return true
	}
	return false
}

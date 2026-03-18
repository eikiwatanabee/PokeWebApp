package entity

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type WeeklyEventType string

const (
	WeeklyEventCommitRush   WeeklyEventType = "commit_rush"
	WeeklyEventPRMarathon   WeeklyEventType = "pr_marathon"
	WeeklyEventReviewStar   WeeklyEventType = "review_star"
	WeeklyEventIssueCrusher WeeklyEventType = "issue_crusher"
	WeeklyEventAllRounder   WeeklyEventType = "all_rounder"
)

type WeeklyEventDefinition struct {
	Type        WeeklyEventType
	Title       string
	Description string
	Icon        string
	ScoreType   GitHubEventType // which event type counts; empty means all
}

var WeeklyEventDefinitions = []WeeklyEventDefinition{
	{WeeklyEventCommitRush, "コミットラッシュ", "今週一番コミットしたチームが勝利！", "💻", EventCommit},
	{WeeklyEventPRMarathon, "PRマラソン", "マージされたPR数で競おう！", "🎉", EventPRMerge},
	{WeeklyEventReviewStar, "レビュースター", "レビュー数で勝負！コードレビューで貢献しよう", "👀", EventReview},
	{WeeklyEventIssueCrusher, "Issue クラッシャー", "Issue解決数で競おう！", "✅", EventIssueClose},
	{WeeklyEventAllRounder, "オールラウンダー", "全アクティビティのXP合計で勝負！", "⭐", ""},
}

type WeeklyEvent struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	EventType WeeklyEventType
	WeekStart time.Time // Monday of the week
	WeekEnd   time.Time // Sunday end
	CreatedAt time.Time
}

// CurrentWeekStart returns the Monday 00:00 of the current week.
func CurrentWeekStart() time.Time {
	now := time.Now().UTC()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -int(weekday-time.Monday))
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
}

// GenerateWeeklyEvent deterministically generates a weekly event for a given week.
func GenerateWeeklyEvent(tenantID uuid.UUID, weekStart time.Time) *WeeklyEvent {
	// Seed based on tenant + week for determinism
	seed := int64(weekStart.Unix()) + int64(tenantID[0])<<8 + int64(tenantID[1])
	rng := rand.New(rand.NewSource(seed))
	def := WeeklyEventDefinitions[rng.Intn(len(WeeklyEventDefinitions))]

	weekEnd := weekStart.AddDate(0, 0, 7).Add(-time.Second)

	return &WeeklyEvent{
		ID:        uuid.New(),
		TenantID:  tenantID,
		EventType: def.Type,
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
		CreatedAt: time.Now(),
	}
}

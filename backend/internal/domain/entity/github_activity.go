package entity

import (
	"time"

	"github.com/google/uuid"
)

type GitHubEventType string

const (
	EventCommit      GitHubEventType = "commit"
	EventPRMerge     GitHubEventType = "pr_merge"
	EventPROpen      GitHubEventType = "pr_open"
	EventIssueClose  GitHubEventType = "issue_close"
	EventReview      GitHubEventType = "review"
)

type GitHubActivity struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	EventType   GitHubEventType
	RepoName    string
	Title       string
	URL         string
	XP          int
	CreatedAt   time.Time
}

var xpTable = map[GitHubEventType]int{
	EventCommit:     10,
	EventPRMerge:    50,
	EventPROpen:     20,
	EventIssueClose: 15,
	EventReview:     25,
}

func NewGitHubActivity(userID uuid.UUID, eventType GitHubEventType, repoName, title, url string) *GitHubActivity {
	xp := xpTable[eventType]
	if xp == 0 {
		xp = 5
	}
	return &GitHubActivity{
		ID:        uuid.New(),
		UserID:    userID,
		EventType: eventType,
		RepoName:  repoName,
		Title:     title,
		URL:       url,
		XP:        xp,
		CreatedAt: time.Now(),
	}
}

func XPForEvent(eventType GitHubEventType) int {
	xp := xpTable[eventType]
	if xp == 0 {
		return 5
	}
	return xp
}

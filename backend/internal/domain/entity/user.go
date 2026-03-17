package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
)

type User struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	TeamID         *uuid.UUID
	GoogleID       string
	GitHubID       string
	GitHubUsername  string
	Email          string
	Name           string
	AvatarURL      string
	TotalXP        int
	Level          int
	Role           UserRole
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewUser(tenantID uuid.UUID, googleID, email, name string) (*User, error) {
	if googleID == "" {
		return nil, ErrEmptyGoogleID
	}
	if email == "" {
		return nil, ErrEmptyEmail
	}
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		GoogleID:  googleID,
		Email:     email,
		Name:      name,
		TotalXP:   0,
		Level:     1,
		Role:      RoleMember,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func NewGitHubUser(tenantID uuid.UUID, githubID, username, email, name, avatarURL string) (*User, error) {
	if githubID == "" {
		return nil, ErrEmptyGitHubID
	}
	if email == "" {
		return nil, ErrEmptyEmail
	}
	now := time.Now()
	return &User{
		ID:            uuid.New(),
		TenantID:      tenantID,
		GoogleID:      "github:" + githubID,
		GitHubID:      githubID,
		GitHubUsername: username,
		Email:         email,
		Name:          name,
		AvatarURL:     avatarURL,
		TotalXP:       0,
		Level:         1,
		Role:          RoleMember,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) PromoteToAdmin() {
	u.Role = RoleAdmin
	u.UpdatedAt = time.Now()
}

func (u *User) JoinTeam(teamID uuid.UUID) {
	u.TeamID = &teamID
	u.UpdatedAt = time.Now()
}

func (u *User) AddXP(xp int) {
	u.TotalXP += xp
	u.Level = calculateLevel(u.TotalXP)
	u.UpdatedAt = time.Now()
}

func calculateLevel(totalXP int) int {
	// Each level requires progressively more XP
	// Level 1: 0, Level 2: 100, Level 3: 250, Level 4: 450, ...
	level := 1
	required := 100
	remaining := totalXP
	for remaining >= required {
		remaining -= required
		level++
		required = 100 + (level-1)*50
	}
	return level
}

func (u *User) XPToNextLevel() int {
	required := 100 + (u.Level-1)*50
	// Calculate XP spent on previous levels
	spent := 0
	for l := 1; l < u.Level; l++ {
		spent += 100 + (l-1)*50
	}
	return required - (u.TotalXP - spent)
}

package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Infrastructure models — separate from domain entities.
// These have GORM tags and handle DB mapping.

type TenantModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (TenantModel) TableName() string { return "tenants" }

type UserModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	TeamID         *uuid.UUID `gorm:"type:uuid;index"`
	GoogleID       string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	GitHubID       string     `gorm:"type:varchar(255);index"`
	GitHubUsername string     `gorm:"type:varchar(255);index"`
	Email          string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	Name           string     `gorm:"type:varchar(255);not null"`
	AvatarURL      string     `gorm:"type:varchar(500)"`
	TotalXP        int        `gorm:"type:int;not null;default:0"`
	Level          int        `gorm:"type:int;not null;default:1"`
	Role           string     `gorm:"type:varchar(50);not null;default:'member'"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}

type TeamModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (TeamModel) TableName() string { return "teams" }

func (UserModel) TableName() string { return "users" }

type BookModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	Title      string     `gorm:"type:varchar(255);not null"`
	Author     string     `gorm:"type:varchar(255);not null"`
	Status     string     `gorm:"type:varchar(50);not null;default:'unread'"`
	FinishedAt *time.Time `gorm:"type:timestamp"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
	Tags       []TagModel `gorm:"many2many:book_tags;"`
}

func (BookModel) TableName() string { return "books" }

type MemoModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BookID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (MemoModel) TableName() string { return "memos" }

type TagModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (TagModel) TableName() string { return "tags" }

type UserPokemonModel struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	BookID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	PokemonID     int            `gorm:"type:int;not null"`
	PokemonName   string         `gorm:"type:varchar(100);not null"`
	PokemonSprite string         `gorm:"type:varchar(500);not null"`
	PokemonTypes  datatypes.JSON `gorm:"type:jsonb"`
	CaughtAt      time.Time      `gorm:"autoCreateTime"`
}

func (UserPokemonModel) TableName() string { return "user_pokemon" }

type GitHubActivityModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	EventType string    `gorm:"type:varchar(50);not null"`
	RepoName  string    `gorm:"type:varchar(255);not null"`
	Title     string    `gorm:"type:varchar(500);not null"`
	URL       string    `gorm:"type:varchar(500)"`
	XP        int       `gorm:"type:int;not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (GitHubActivityModel) TableName() string { return "github_activities" }

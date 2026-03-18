package entity

import (
	"time"

	"github.com/google/uuid"
)

// LoginBonus records a daily login bonus claim.
type LoginBonus struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Date          time.Time
	BonusXP       int
	ConsecutiveDays int
	CreatedAt     time.Time
}

// LoginBonusXP returns the XP reward based on consecutive login days.
func LoginBonusXP(consecutiveDays int) int {
	switch {
	case consecutiveDays >= 30:
		return 100
	case consecutiveDays >= 14:
		return 75
	case consecutiveDays >= 7:
		return 50
	case consecutiveDays >= 3:
		return 30
	default:
		return 10
	}
}

// NewLoginBonus creates a login bonus record.
func NewLoginBonus(userID uuid.UUID, date time.Time, consecutiveDays int) *LoginBonus {
	return &LoginBonus{
		ID:              uuid.New(),
		UserID:          userID,
		Date:            date.Truncate(24 * time.Hour),
		BonusXP:         LoginBonusXP(consecutiveDays),
		ConsecutiveDays: consecutiveDays,
		CreatedAt:       time.Now(),
	}
}

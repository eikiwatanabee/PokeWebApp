package command

import (
	"context"
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type ClaimLoginBonusCommand struct {
	UserID string
}

type ClaimLoginBonusResult struct {
	Claimed         bool `json:"claimed"`
	BonusXP         int  `json:"bonus_xp"`
	ConsecutiveDays int  `json:"consecutive_days"`
	AlreadyClaimed  bool `json:"already_claimed"`
}

type ClaimLoginBonusHandler struct {
	userRepo       repository.UserRepository
	loginBonusRepo repository.LoginBonusRepository
}

func NewClaimLoginBonusHandler(
	userRepo repository.UserRepository,
	loginBonusRepo repository.LoginBonusRepository,
) *ClaimLoginBonusHandler {
	return &ClaimLoginBonusHandler{
		userRepo:       userRepo,
		loginBonusRepo: loginBonusRepo,
	}
}

func (h *ClaimLoginBonusHandler) Handle(ctx context.Context, cmd *ClaimLoginBonusCommand) (*ClaimLoginBonusResult, error) {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)

	// Check if already claimed today
	existing, err := h.loginBonusRepo.FindByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return &ClaimLoginBonusResult{
			Claimed:         false,
			BonusXP:         existing.BonusXP,
			ConsecutiveDays: existing.ConsecutiveDays,
			AlreadyClaimed:  true,
		}, nil
	}

	// Calculate consecutive days
	consecutiveDays := 1
	latest, err := h.loginBonusRepo.FindLatestByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if latest != nil {
		yesterday := today.Add(-24 * time.Hour)
		if latest.Date.Equal(yesterday) {
			consecutiveDays = latest.ConsecutiveDays + 1
		}
	}

	// Create bonus and add XP
	bonus := entity.NewLoginBonus(userID, today, consecutiveDays)
	if err := h.loginBonusRepo.Save(ctx, bonus); err != nil {
		return nil, err
	}

	// Add XP to user
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		user.AddXP(bonus.BonusXP)
		if err := h.userRepo.Save(ctx, user); err != nil {
			return nil, err
		}
	}

	return &ClaimLoginBonusResult{
		Claimed:         true,
		BonusXP:         bonus.BonusXP,
		ConsecutiveDays: consecutiveDays,
	}, nil
}

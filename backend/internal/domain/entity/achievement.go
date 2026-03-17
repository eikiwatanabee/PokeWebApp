package entity

import (
	"time"

	"github.com/google/uuid"
)

type AchievementType string

const (
	AchievementFirstCommit   AchievementType = "first_commit"
	AchievementStreak7       AchievementType = "streak_7"
	AchievementStreak30      AchievementType = "streak_30"
	AchievementStreak100     AchievementType = "streak_100"
	AchievementLevel10       AchievementType = "level_10"
	AchievementLevel25       AchievementType = "level_25"
	AchievementLevel50       AchievementType = "level_50"
	AchievementPokemon10     AchievementType = "pokemon_10"
	AchievementPokemon50     AchievementType = "pokemon_50"
	AchievementPokemon100    AchievementType = "pokemon_100"
	AchievementReview10      AchievementType = "review_10"
	AchievementLegendary     AchievementType = "legendary_catch"
)

type AchievementDefinition struct {
	Type        AchievementType
	Name        string
	Description string
	Icon        string
}

var AchievementDefinitions = []AchievementDefinition{
	{AchievementFirstCommit, "はじめのいっぽ", "最初のコミット", "🐣"},
	{AchievementStreak7, "ほのおのトレーナー", "7日連続コミット", "🔥"},
	{AchievementStreak30, "てつじんトレーナー", "30日連続コミット", "🏔️"},
	{AchievementStreak100, "でんせつトレーナー", "100日連続コミット", "⭐"},
	{AchievementLevel10, "やまのぬし", "レベル10到達", "🏔️"},
	{AchievementLevel25, "エリートトレーナー", "レベル25到達", "🎖️"},
	{AchievementLevel50, "チャンピオン", "レベル50到達", "👑"},
	{AchievementPokemon10, "コレクター", "ポケモン10匹ゲット", "📦"},
	{AchievementPokemon50, "ずかんコレクター", "ポケモン50匹ゲット", "📚"},
	{AchievementPokemon100, "ポケモンマスター", "ポケモン100匹ゲット", "🏆"},
	{AchievementReview10, "なかまおもい", "10回レビュー", "🤝"},
	{AchievementLegendary, "でんせつ", "Legendaryポケモンゲット", "🌟"},
}

type UserAchievement struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	AchievementType AchievementType
	UnlockedAt      time.Time
}

func NewUserAchievement(userID uuid.UUID, achievementType AchievementType) *UserAchievement {
	return &UserAchievement{
		ID:              uuid.New(),
		UserID:          userID,
		AchievementType: achievementType,
		UnlockedAt:      time.Now(),
	}
}

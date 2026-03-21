package entity

import (
	"time"

	"github.com/google/uuid"
)

type AchievementType string

const (
	// --- コミット系 ---
	AchievementFirstCommit  AchievementType = "first_commit"
	AchievementCommit50     AchievementType = "commit_50"
	AchievementCommit100    AchievementType = "commit_100"
	AchievementCommit500    AchievementType = "commit_500"
	AchievementCommit1000   AchievementType = "commit_1000"

	// --- ストリーク系 ---
	AchievementStreak3      AchievementType = "streak_3"
	AchievementStreak7      AchievementType = "streak_7"
	AchievementStreak14     AchievementType = "streak_14"
	AchievementStreak30     AchievementType = "streak_30"
	AchievementStreak60     AchievementType = "streak_60"
	AchievementStreak100    AchievementType = "streak_100"
	AchievementStreak365    AchievementType = "streak_365"

	// --- レベル系 ---
	AchievementLevel5       AchievementType = "level_5"
	AchievementLevel10      AchievementType = "level_10"
	AchievementLevel25      AchievementType = "level_25"
	AchievementLevel50      AchievementType = "level_50"
	AchievementLevel100     AchievementType = "level_100"

	// --- XP系 ---
	AchievementXP1000       AchievementType = "xp_1000"
	AchievementXP5000       AchievementType = "xp_5000"
	AchievementXP10000      AchievementType = "xp_10000"
	AchievementXP50000      AchievementType = "xp_50000"

	// --- ポケモン系 ---
	AchievementPokemon1     AchievementType = "pokemon_1"
	AchievementPokemon10    AchievementType = "pokemon_10"
	AchievementPokemon25    AchievementType = "pokemon_25"
	AchievementPokemon50    AchievementType = "pokemon_50"
	AchievementPokemon100   AchievementType = "pokemon_100"
	AchievementPokemon200   AchievementType = "pokemon_200"
	AchievementPokemon500   AchievementType = "pokemon_500"

	// --- レアリティ系 ---
	AchievementRareCatch    AchievementType = "rare_catch"
	AchievementEpicCatch    AchievementType = "epic_catch"
	AchievementLegendary    AchievementType = "legendary_catch"

	// --- PR系 ---
	AchievementFirstMerge   AchievementType = "first_merge"
	AchievementMerge10      AchievementType = "merge_10"
	AchievementMerge25      AchievementType = "merge_25"
	AchievementMerge50      AchievementType = "merge_50"
	AchievementMerge100     AchievementType = "merge_100"

	// --- レビュー系 ---
	AchievementFirstReview  AchievementType = "first_review"
	AchievementReview10     AchievementType = "review_10"
	AchievementReview25     AchievementType = "review_25"
	AchievementReview50     AchievementType = "review_50"
	AchievementReview100    AchievementType = "review_100"

	// --- Issue系 ---
	AchievementFirstIssue   AchievementType = "first_issue"
	AchievementIssue10      AchievementType = "issue_10"
	AchievementIssue50      AchievementType = "issue_50"

	// --- デプロイ系 ---
	AchievementFirstDeploy  AchievementType = "first_deploy"
	AchievementDeploy10     AchievementType = "deploy_10"
	AchievementDeploy50     AchievementType = "deploy_50"
	AchievementDeploy100    AchievementType = "deploy_100"

	// --- 特殊系 ---
	AchievementAllRounder   AchievementType = "all_rounder"
)

type AchievementDefinition struct {
	Type        AchievementType
	Name        string
	Description string
	Icon        string
	Category    string
}

var AchievementDefinitions = []AchievementDefinition{
	// コミット
	{AchievementFirstCommit, "はじめのいっぽ", "最初のコミット", "🐣", "commit"},
	{AchievementCommit50, "コードのたまご", "50回コミット", "🥚", "commit"},
	{AchievementCommit100, "キーボードウォリアー", "100回コミット", "⌨️", "commit"},
	{AchievementCommit500, "コードマシン", "500回コミット", "🤖", "commit"},
	{AchievementCommit1000, "リビングレジェンド", "1000回コミット", "🏛️", "commit"},

	// ストリーク
	{AchievementStreak3, "みならいトレーナー", "3日連続コミット", "🌱", "streak"},
	{AchievementStreak7, "ほのおのトレーナー", "7日連続コミット", "🔥", "streak"},
	{AchievementStreak14, "きあいのハチマキ", "14日連続コミット", "💪", "streak"},
	{AchievementStreak30, "てつじんトレーナー", "30日連続コミット", "🏔️", "streak"},
	{AchievementStreak60, "ふくつのこころ", "60日連続コミット", "💎", "streak"},
	{AchievementStreak100, "でんせつトレーナー", "100日連続コミット", "⭐", "streak"},
	{AchievementStreak365, "いちねんせい", "365日連続コミット", "🎓", "streak"},

	// レベル
	{AchievementLevel5, "かけだしトレーナー", "レベル5到達", "🎒", "level"},
	{AchievementLevel10, "やまのぬし", "レベル10到達", "🏔️", "level"},
	{AchievementLevel25, "エリートトレーナー", "レベル25到達", "🎖️", "level"},
	{AchievementLevel50, "チャンピオン", "レベル50到達", "👑", "level"},
	{AchievementLevel100, "ポケモンマスター", "レベル100到達", "🌈", "level"},

	// XP
	{AchievementXP1000, "けいけんち1000", "累計1,000 XP獲得", "✨", "xp"},
	{AchievementXP5000, "けいけんち5000", "累計5,000 XP獲得", "💫", "xp"},
	{AchievementXP10000, "けいけんちまんてん", "累計10,000 XP獲得", "🌟", "xp"},
	{AchievementXP50000, "きょうがくのXP", "累計50,000 XP獲得", "💥", "xp"},

	// ポケモン
	{AchievementPokemon1, "はじめてのゲット", "最初のポケモンゲット", "⚡", "pokemon"},
	{AchievementPokemon10, "コレクター", "ポケモン10匹ゲット", "📦", "pokemon"},
	{AchievementPokemon25, "ポケモンずき", "ポケモン25匹ゲット", "💕", "pokemon"},
	{AchievementPokemon50, "ずかんコレクター", "ポケモン50匹ゲット", "📚", "pokemon"},
	{AchievementPokemon100, "ずかんかんせい", "ポケモン100匹ゲット", "🏆", "pokemon"},
	{AchievementPokemon200, "ポケモンはかせ", "ポケモン200匹ゲット", "🔬", "pokemon"},
	{AchievementPokemon500, "ポケモンマニア", "ポケモン500匹ゲット", "🗾", "pokemon"},

	// レアリティ
	{AchievementRareCatch, "レアハンター", "Rareポケモンゲット", "🟦", "rarity"},
	{AchievementEpicCatch, "エピックハンター", "Epicポケモンゲット", "🟪", "rarity"},
	{AchievementLegendary, "でんせつとのであい", "Legendaryポケモンゲット", "🟨", "rarity"},

	// PR
	{AchievementFirstMerge, "はじめてのマージ", "最初のPRマージ", "🎉", "pr"},
	{AchievementMerge10, "マージャー", "PRを10回マージ", "🔀", "pr"},
	{AchievementMerge25, "マージマスター", "PRを25回マージ", "🔄", "pr"},
	{AchievementMerge50, "プルリクエスト職人", "PRを50回マージ", "🛠️", "pr"},
	{AchievementMerge100, "マージキング", "PRを100回マージ", "🫅", "pr"},

	// レビュー
	{AchievementFirstReview, "はじめてのレビュー", "最初のコードレビュー", "👀", "review"},
	{AchievementReview10, "なかまおもい", "10回レビュー", "🤝", "review"},
	{AchievementReview25, "コードリーダー", "25回レビュー", "📖", "review"},
	{AchievementReview50, "レビューマスター", "50回レビュー", "🧐", "review"},
	{AchievementReview100, "コードガーディアン", "100回レビュー", "🛡️", "review"},

	// Issue
	{AchievementFirstIssue, "バグハンター", "最初のIssue完了", "🐛", "issue"},
	{AchievementIssue10, "イシュークラッシャー", "10個のIssue完了", "🔨", "issue"},
	{AchievementIssue50, "イシューマスター", "50個のIssue完了", "⚔️", "issue"},

	// デプロイ
	{AchievementFirstDeploy, "はじめてのデプロイ", "最初のデプロイ成功", "🚀", "deploy"},
	{AchievementDeploy10, "デプロイマスター", "10回デプロイ", "🛸", "deploy"},
	{AchievementDeploy50, "デプロイの鬼", "50回デプロイ", "🌍", "deploy"},
	{AchievementDeploy100, "インフラの神", "100回デプロイ", "🏗️", "deploy"},

	// 特殊
	{AchievementAllRounder, "オールラウンダー", "コミット・PR・レビュー・Issue全種達成", "🎯", "special"},
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

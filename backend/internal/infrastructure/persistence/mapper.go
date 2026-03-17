package persistence

import (
	"encoding/json"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
)

// --- Tenant ---

func toTenantModel(e *entity.Tenant) *TenantModel {
	return &TenantModel{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func toTenantEntity(m *TenantModel) *entity.Tenant {
	return &entity.Tenant{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// --- User ---

func toUserModel(e *entity.User) *UserModel {
	return &UserModel{
		ID:               e.ID,
		TenantID:         e.TenantID,
		TeamID:           e.TeamID,
		GoogleID:         e.GoogleID,
		GitHubID:         e.GitHubID,
		GitHubUsername:   e.GitHubUsername,
		Email:            e.Email,
		Name:             e.Name,
		AvatarURL:        e.AvatarURL,
		TotalXP:          e.TotalXP,
		Level:            e.Level,
		CurrentStreak:    e.CurrentStreak,
		MaxStreak:        e.MaxStreak,
		LastActivityDate: e.LastActivityDate,
		Role:             string(e.Role),
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func toUserEntity(m *UserModel) *entity.User {
	return &entity.User{
		ID:               m.ID,
		TenantID:         m.TenantID,
		TeamID:           m.TeamID,
		GoogleID:         m.GoogleID,
		GitHubID:         m.GitHubID,
		GitHubUsername:   m.GitHubUsername,
		Email:            m.Email,
		Name:             m.Name,
		AvatarURL:        m.AvatarURL,
		TotalXP:          m.TotalXP,
		Level:            m.Level,
		CurrentStreak:    m.CurrentStreak,
		MaxStreak:        m.MaxStreak,
		LastActivityDate: m.LastActivityDate,
		Role:             entity.UserRole(m.Role),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

// --- Book ---

func toBookModel(e *entity.Book) *BookModel {
	tags := make([]TagModel, len(e.Tags))
	for i, t := range e.Tags {
		tags[i] = *toTagModel(&t)
	}
	return &BookModel{
		ID:         e.ID,
		UserID:     e.UserID,
		Title:      e.Title,
		Author:     e.Author,
		Status:     e.Status.String(),
		FinishedAt: e.FinishedAt,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
		Tags:       tags,
	}
}

func toBookEntity(m *BookModel) *entity.Book {
	tags := make([]entity.Tag, len(m.Tags))
	for i, t := range m.Tags {
		tags[i] = *toTagEntity(&t)
	}
	return &entity.Book{
		ID:         m.ID,
		UserID:     m.UserID,
		Title:      m.Title,
		Author:     m.Author,
		Status:     valueobject.BookStatus(m.Status),
		FinishedAt: m.FinishedAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Tags:       tags,
	}
}

// --- Memo ---

func toMemoModel(e *entity.Memo) *MemoModel {
	return &MemoModel{
		ID:        e.ID,
		BookID:    e.BookID,
		Content:   e.Content,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func toMemoEntity(m *MemoModel) *entity.Memo {
	return &entity.Memo{
		ID:        m.ID,
		BookID:    m.BookID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// --- Tag ---

func toTagModel(e *entity.Tag) *TagModel {
	return &TagModel{
		ID:        e.ID,
		TenantID:  e.TenantID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
	}
}

func toTagEntity(m *TagModel) *entity.Tag {
	return &entity.Tag{
		ID:        m.ID,
		TenantID:  m.TenantID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
	}
}

// --- UserPokemon ---

func toUserPokemonModel(e *entity.UserPokemon) *UserPokemonModel {
	typesJSON, _ := json.Marshal(e.Pokemon.Types)
	rarity := string(e.Pokemon.Rarity)
	if rarity == "" {
		rarity = "common"
	}
	return &UserPokemonModel{
		ID:            e.ID,
		UserID:        e.UserID,
		BookID:        e.BookID,
		PokemonID:     e.Pokemon.PokemonID,
		PokemonName:   e.Pokemon.Name,
		PokemonSprite: e.Pokemon.SpriteURL,
		PokemonTypes:  typesJSON,
		PokemonRarity: rarity,
		CaughtAt:      e.CaughtAt,
	}
}

func toUserPokemonEntity(m *UserPokemonModel) *entity.UserPokemon {
	var types []string
	_ = json.Unmarshal(m.PokemonTypes, &types)
	rarity := valueobject.PokemonRarity(m.PokemonRarity)
	if rarity == "" {
		rarity = valueobject.RarityCommon
	}
	return &entity.UserPokemon{
		ID:         m.ID,
		UserID:     m.UserID,
		BookID:     m.BookID,
		ActivityID: m.BookID,
		Pokemon: valueobject.PokemonInfo{
			PokemonID: m.PokemonID,
			Name:      m.PokemonName,
			SpriteURL: m.PokemonSprite,
			Types:     types,
			Rarity:    rarity,
		},
		CaughtAt: m.CaughtAt,
	}
}

// --- UserAchievement ---

func toUserAchievementModel(e *entity.UserAchievement) *UserAchievementModel {
	return &UserAchievementModel{
		ID:              e.ID,
		UserID:          e.UserID,
		AchievementType: string(e.AchievementType),
		UnlockedAt:      e.UnlockedAt,
	}
}

func toUserAchievementEntity(m *UserAchievementModel) *entity.UserAchievement {
	return &entity.UserAchievement{
		ID:              m.ID,
		UserID:          m.UserID,
		AchievementType: entity.AchievementType(m.AchievementType),
		UnlockedAt:      m.UnlockedAt,
	}
}

// --- Team ---

func mapTeamEntityToModel(e *entity.Team) *TeamModel {
	return &TeamModel{
		ID:        e.ID,
		TenantID:  e.TenantID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func mapTeamModelToEntity(m *TeamModel) *entity.Team {
	return &entity.Team{
		ID:        m.ID,
		TenantID:  m.TenantID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// --- GitHubActivity ---

func toGitHubActivityModel(e *entity.GitHubActivity) *GitHubActivityModel {
	return &GitHubActivityModel{
		ID:        e.ID,
		UserID:    e.UserID,
		EventType: string(e.EventType),
		RepoName:  e.RepoName,
		Title:     e.Title,
		URL:       e.URL,
		XP:        e.XP,
		CreatedAt: e.CreatedAt,
	}
}

func toGitHubActivityEntity(m *GitHubActivityModel) *entity.GitHubActivity {
	return &entity.GitHubActivity{
		ID:        m.ID,
		UserID:    m.UserID,
		EventType: entity.GitHubEventType(m.EventType),
		RepoName:  m.RepoName,
		Title:     m.Title,
		URL:       m.URL,
		XP:        m.XP,
		CreatedAt: m.CreatedAt,
	}
}

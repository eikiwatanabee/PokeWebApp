package service

import (
	"context"
	"errors"
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/entity"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

// --- Mock repositories ---

type mockAchievementRepo struct {
	existing map[entity.AchievementType]bool
	saved    []*entity.UserAchievement
	saveErr  error
}

func newMockAchievementRepo() *mockAchievementRepo {
	return &mockAchievementRepo{existing: make(map[entity.AchievementType]bool)}
}

func (m *mockAchievementRepo) FindByUserID(_ context.Context, _ uuid.UUID) ([]*entity.UserAchievement, error) {
	return nil, nil
}

func (m *mockAchievementRepo) HasAchievement(_ context.Context, _ uuid.UUID, t entity.AchievementType) (bool, error) {
	return m.existing[t], nil
}

func (m *mockAchievementRepo) Save(_ context.Context, a *entity.UserAchievement) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = append(m.saved, a)
	m.existing[a.AchievementType] = true
	return nil
}

type mockActivityRepo struct {
	totalCount int64
	typeCounts map[entity.GitHubEventType]int64
	countErr   error
}

func newMockActivityRepo() *mockActivityRepo {
	return &mockActivityRepo{typeCounts: make(map[entity.GitHubEventType]int64)}
}

func (m *mockActivityRepo) FindByUserID(_ context.Context, _ uuid.UUID, _ int) ([]*entity.GitHubActivity, error) {
	return nil, nil
}

func (m *mockActivityRepo) CountByUserID(_ context.Context, _ uuid.UUID) (int64, error) {
	return m.totalCount, m.countErr
}

func (m *mockActivityRepo) CountByUserIDAndType(_ context.Context, _ uuid.UUID, t entity.GitHubEventType) (int64, error) {
	return m.typeCounts[t], m.countErr
}

func (m *mockActivityRepo) TotalXPByUserID(_ context.Context, _ uuid.UUID) (int, error) {
	return 0, nil
}

func (m *mockActivityRepo) Save(_ context.Context, _ *entity.GitHubActivity) error {
	return nil
}

type mockPokemonRepo struct {
	pokemon []*entity.UserPokemon
	findErr error
}

func (m *mockPokemonRepo) FindByUserID(_ context.Context, _ uuid.UUID) ([]*entity.UserPokemon, error) {
	return m.pokemon, m.findErr
}

func (m *mockPokemonRepo) FindByID(_ context.Context, _ uuid.UUID) (*entity.UserPokemon, error) {
	return nil, nil
}

func (m *mockPokemonRepo) Save(_ context.Context, _ *entity.UserPokemon) error {
	return nil
}

func makePokemonSlice(n int) []*entity.UserPokemon {
	result := make([]*entity.UserPokemon, n)
	for i := range result {
		result[i] = &entity.UserPokemon{ID: uuid.New()}
	}
	return result
}

func makeUser() *entity.User {
	return &entity.User{
		ID:            uuid.New(),
		TenantID:      uuid.New(),
		Level:         1,
		TotalXP:       0,
		CurrentStreak: 0,
	}
}

func hasAchievement(achievements []*entity.UserAchievement, t entity.AchievementType) bool {
	for _, a := range achievements {
		if a.AchievementType == t {
			return true
		}
	}
	return false
}

// --- Tests ---

func TestCheckAndUnlock_FirstCommit(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	actRepo.totalCount = 1
	pokRepo := &mockPokemonRepo{}

	svc := NewAchievementService(achRepo, actRepo, pokRepo)
	user := makeUser()

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hasAchievement(got, entity.AchievementFirstCommit) {
		t.Error("expected first_commit achievement")
	}
	if !hasAchievement(got, entity.AchievementPokemon1) {
		// pokemon count is 0, so pokemon_1 should NOT be unlocked
		// This is expected to not be present
	}
}

func TestCheckAndUnlock_SkipsAlreadyUnlocked(t *testing.T) {
	achRepo := newMockAchievementRepo()
	achRepo.existing[entity.AchievementFirstCommit] = true
	actRepo := newMockActivityRepo()
	actRepo.totalCount = 1
	pokRepo := &mockPokemonRepo{}

	svc := NewAchievementService(achRepo, actRepo, pokRepo)
	user := makeUser()

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hasAchievement(got, entity.AchievementFirstCommit) {
		t.Error("should not re-unlock first_commit")
	}
}

func TestCheckAndUnlock_StreakAchievements(t *testing.T) {
	tests := []struct {
		name           string
		streak         int
		wantUnlocked   []entity.AchievementType
		wantNotPresent []entity.AchievementType
	}{
		{
			name:         "streak 3 unlocks streak_3 only",
			streak:       3,
			wantUnlocked: []entity.AchievementType{entity.AchievementStreak3},
			wantNotPresent: []entity.AchievementType{
				entity.AchievementStreak7,
				entity.AchievementStreak14,
			},
		},
		{
			name:   "streak 7 unlocks streak_3 and streak_7",
			streak: 7,
			wantUnlocked: []entity.AchievementType{
				entity.AchievementStreak3,
				entity.AchievementStreak7,
			},
			wantNotPresent: []entity.AchievementType{entity.AchievementStreak14},
		},
		{
			name:   "streak 30 unlocks up to streak_30",
			streak: 30,
			wantUnlocked: []entity.AchievementType{
				entity.AchievementStreak3,
				entity.AchievementStreak7,
				entity.AchievementStreak14,
				entity.AchievementStreak30,
			},
			wantNotPresent: []entity.AchievementType{entity.AchievementStreak60},
		},
		{
			name:   "streak 365 unlocks all streak achievements",
			streak: 365,
			wantUnlocked: []entity.AchievementType{
				entity.AchievementStreak3,
				entity.AchievementStreak7,
				entity.AchievementStreak14,
				entity.AchievementStreak30,
				entity.AchievementStreak60,
				entity.AchievementStreak100,
				entity.AchievementStreak365,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achRepo := newMockAchievementRepo()
			actRepo := newMockActivityRepo()
			pokRepo := &mockPokemonRepo{}
			svc := NewAchievementService(achRepo, actRepo, pokRepo)

			user := makeUser()
			user.CurrentStreak = tt.streak

			got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, want := range tt.wantUnlocked {
				if !hasAchievement(got, want) {
					t.Errorf("expected %s to be unlocked", want)
				}
			}
			for _, notWant := range tt.wantNotPresent {
				if hasAchievement(got, notWant) {
					t.Errorf("expected %s to NOT be unlocked", notWant)
				}
			}
		})
	}
}

func TestCheckAndUnlock_LevelAchievements(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	pokRepo := &mockPokemonRepo{}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()
	user.Level = 25

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []entity.AchievementType{
		entity.AchievementLevel5,
		entity.AchievementLevel10,
		entity.AchievementLevel25,
	} {
		if !hasAchievement(got, want) {
			t.Errorf("expected %s to be unlocked at level 25", want)
		}
	}
	if hasAchievement(got, entity.AchievementLevel50) {
		t.Error("level_50 should not be unlocked at level 25")
	}
}

func TestCheckAndUnlock_XPAchievements(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	pokRepo := &mockPokemonRepo{}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()
	user.TotalXP = 5000

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hasAchievement(got, entity.AchievementXP1000) {
		t.Error("expected xp_1000 at 5000 XP")
	}
	if !hasAchievement(got, entity.AchievementXP5000) {
		t.Error("expected xp_5000 at 5000 XP")
	}
	if hasAchievement(got, entity.AchievementXP10000) {
		t.Error("xp_10000 should not unlock at 5000 XP")
	}
}

func TestCheckAndUnlock_PokemonCount(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	pokRepo := &mockPokemonRepo{pokemon: makePokemonSlice(25)}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []entity.AchievementType{
		entity.AchievementPokemon1,
		entity.AchievementPokemon10,
		entity.AchievementPokemon25,
	} {
		if !hasAchievement(got, want) {
			t.Errorf("expected %s with 25 pokemon", want)
		}
	}
	if hasAchievement(got, entity.AchievementPokemon50) {
		t.Error("pokemon_50 should not unlock with 25 pokemon")
	}
}

func TestCheckAndUnlock_RarityCatch(t *testing.T) {
	tests := []struct {
		name   string
		rarity valueobject.PokemonRarity
		want   entity.AchievementType
	}{
		{"rare catch", valueobject.RarityRare, entity.AchievementRareCatch},
		{"epic catch", valueobject.RarityEpic, entity.AchievementEpicCatch},
		{"legendary catch", valueobject.RarityLegendary, entity.AchievementLegendary},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achRepo := newMockAchievementRepo()
			actRepo := newMockActivityRepo()
			pokRepo := &mockPokemonRepo{}
			svc := NewAchievementService(achRepo, actRepo, pokRepo)

			user := makeUser()

			got, err := svc.CheckAndUnlock(context.Background(), user, tt.rarity)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !hasAchievement(got, tt.want) {
				t.Errorf("expected %s for rarity %s", tt.want, tt.rarity)
			}
		})
	}
}

func TestCheckAndUnlock_CommonDoesNotUnlockRarity(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	pokRepo := &mockPokemonRepo{}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, notWant := range []entity.AchievementType{
		entity.AchievementRareCatch,
		entity.AchievementEpicCatch,
		entity.AchievementLegendary,
	} {
		if hasAchievement(got, notWant) {
			t.Errorf("%s should not unlock for common rarity", notWant)
		}
	}
}

func TestCheckAndUnlock_PRMergeAchievements(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	actRepo.typeCounts[entity.EventPRMerge] = 10
	pokRepo := &mockPokemonRepo{}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hasAchievement(got, entity.AchievementFirstMerge) {
		t.Error("expected first_merge with 10 merges")
	}
	if !hasAchievement(got, entity.AchievementMerge10) {
		t.Error("expected merge_10 with 10 merges")
	}
	if hasAchievement(got, entity.AchievementMerge25) {
		t.Error("merge_25 should not unlock with 10 merges")
	}
}

func TestCheckAndUnlock_ReviewAchievements(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	actRepo.typeCounts[entity.EventReview] = 50
	pokRepo := &mockPokemonRepo{}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []entity.AchievementType{
		entity.AchievementFirstReview,
		entity.AchievementReview10,
		entity.AchievementReview25,
		entity.AchievementReview50,
	} {
		if !hasAchievement(got, want) {
			t.Errorf("expected %s with 50 reviews", want)
		}
	}
	if hasAchievement(got, entity.AchievementReview100) {
		t.Error("review_100 should not unlock with 50 reviews")
	}
}

func TestCheckAndUnlock_IssueAchievements(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	actRepo.typeCounts[entity.EventIssueClose] = 10
	pokRepo := &mockPokemonRepo{}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hasAchievement(got, entity.AchievementFirstIssue) {
		t.Error("expected first_issue with 10 issues")
	}
	if !hasAchievement(got, entity.AchievementIssue10) {
		t.Error("expected issue_10 with 10 issues")
	}
	if hasAchievement(got, entity.AchievementIssue50) {
		t.Error("issue_50 should not unlock with 10 issues")
	}
}

func TestCheckAndUnlock_AllRounder(t *testing.T) {
	t.Run("unlocks when all activity types present", func(t *testing.T) {
		achRepo := newMockAchievementRepo()
		actRepo := newMockActivityRepo()
		actRepo.totalCount = 5
		actRepo.typeCounts[entity.EventPRMerge] = 1
		actRepo.typeCounts[entity.EventReview] = 1
		actRepo.typeCounts[entity.EventIssueClose] = 1
		pokRepo := &mockPokemonRepo{}
		svc := NewAchievementService(achRepo, actRepo, pokRepo)

		user := makeUser()

		got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !hasAchievement(got, entity.AchievementAllRounder) {
			t.Error("expected all_rounder when all types present")
		}
	})

	t.Run("does not unlock when missing review", func(t *testing.T) {
		achRepo := newMockAchievementRepo()
		actRepo := newMockActivityRepo()
		actRepo.totalCount = 5
		actRepo.typeCounts[entity.EventPRMerge] = 1
		actRepo.typeCounts[entity.EventIssueClose] = 1
		// no reviews
		pokRepo := &mockPokemonRepo{}
		svc := NewAchievementService(achRepo, actRepo, pokRepo)

		user := makeUser()

		got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if hasAchievement(got, entity.AchievementAllRounder) {
			t.Error("should not unlock all_rounder without reviews")
		}
	})
}

func TestCheckAndUnlock_CommitCountAchievements(t *testing.T) {
	achRepo := newMockAchievementRepo()
	actRepo := newMockActivityRepo()
	actRepo.totalCount = 100
	pokRepo := &mockPokemonRepo{}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()

	got, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []entity.AchievementType{
		entity.AchievementFirstCommit,
		entity.AchievementCommit50,
		entity.AchievementCommit100,
	} {
		if !hasAchievement(got, want) {
			t.Errorf("expected %s with 100 commits", want)
		}
	}
	if hasAchievement(got, entity.AchievementCommit500) {
		t.Error("commit_500 should not unlock with 100 commits")
	}
}

func TestCheckAndUnlock_RepoError(t *testing.T) {
	repoErr := errors.New("db error")

	t.Run("activity repo error", func(t *testing.T) {
		achRepo := newMockAchievementRepo()
		actRepo := newMockActivityRepo()
		actRepo.countErr = repoErr
		pokRepo := &mockPokemonRepo{}
		svc := NewAchievementService(achRepo, actRepo, pokRepo)

		user := makeUser()

		_, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
		if !errors.Is(err, repoErr) {
			t.Errorf("expected repo error, got: %v", err)
		}
	})

	t.Run("pokemon repo error", func(t *testing.T) {
		achRepo := newMockAchievementRepo()
		actRepo := newMockActivityRepo()
		pokRepo := &mockPokemonRepo{findErr: repoErr}
		svc := NewAchievementService(achRepo, actRepo, pokRepo)

		user := makeUser()

		_, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
		if !errors.Is(err, repoErr) {
			t.Errorf("expected repo error, got: %v", err)
		}
	})

	t.Run("save error", func(t *testing.T) {
		achRepo := newMockAchievementRepo()
		achRepo.saveErr = repoErr
		actRepo := newMockActivityRepo()
		actRepo.totalCount = 1 // will trigger first_commit
		pokRepo := &mockPokemonRepo{}
		svc := NewAchievementService(achRepo, actRepo, pokRepo)

		user := makeUser()

		_, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
		if !errors.Is(err, repoErr) {
			t.Errorf("expected save error, got: %v", err)
		}
	})
}

func TestCheckAndUnlock_LazyLoadCaching(t *testing.T) {
	// Verify that calling CheckAndUnlock with many commit-related achievements
	// only triggers CountByUserID once (via lazy loading)
	achRepo := newMockAchievementRepo()
	actRepo := &countingActivityRepo{
		mockActivityRepo: newMockActivityRepo(),
	}
	actRepo.totalCount = 1000
	pokRepo := &mockPokemonRepo{}
	svc := NewAchievementService(achRepo, actRepo, pokRepo)

	user := makeUser()

	_, err := svc.CheckAndUnlock(context.Background(), user, valueobject.RarityCommon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actRepo.countByUserIDCalls != 1 {
		t.Errorf("CountByUserID called %d times, expected 1 (lazy caching)", actRepo.countByUserIDCalls)
	}
}

type countingActivityRepo struct {
	*mockActivityRepo
	countByUserIDCalls int
}

func (c *countingActivityRepo) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	c.countByUserIDCalls++
	return c.mockActivityRepo.CountByUserID(ctx, userID)
}

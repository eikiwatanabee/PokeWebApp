package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/service"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/auth"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/persistence"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/pokeapi"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/presentation/handler"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/presentation/router"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// --- Database ---
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "pokebookmanager"),
		getEnv("DB_PORT", "5432"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(
		&persistence.TenantModel{},
		&persistence.TeamModel{},
		&persistence.UserModel{},
		&persistence.BookModel{},
		&persistence.MemoModel{},
		&persistence.TagModel{},
		&persistence.UserPokemonModel{},
		&persistence.GitHubActivityModel{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// --- Infrastructure ---
	uow := persistence.NewGormUnitOfWork(db)
	bookRepo := persistence.NewGormBookRepository(db)
	userRepo := persistence.NewGormUserRepository(db)
	tagRepo := persistence.NewGormTagRepository(db)
	memoRepo := persistence.NewGormMemoRepository(db)
	pokemonRepo := persistence.NewGormPokemonRepository(db)
	tenantRepo := persistence.NewGormTenantRepository(db)
	teamRepo := persistence.NewGormTeamRepository(db)
	activityRepo := persistence.NewGormGitHubActivityRepository(db)
	pokeAPIClient := pokeapi.NewClient()

	// --- Auth ---
	jwtManager := auth.NewJWTManager(getEnv("JWT_SECRET", "dev-secret-change-in-production"))
	googleOAuth := auth.NewGoogleOAuth(
		getEnv("GOOGLE_CLIENT_ID", ""),
		getEnv("GOOGLE_CLIENT_SECRET", ""),
		getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
	)
	githubOAuth := auth.NewGitHubOAuth(
		getEnv("GITHUB_CLIENT_ID", ""),
		getEnv("GITHUB_CLIENT_SECRET", ""),
		getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/api/auth/github/callback"),
	)
	webhookVerifier := auth.NewWebhookVerifier(getEnv("GITHUB_WEBHOOK_SECRET", ""))

	// --- Domain Services ---
	gachaSvc := service.NewPokemonGachaService(pokeAPIClient)

	// --- Command Handlers ---
	registerBookHandler := command.NewRegisterBookHandler(uow, bookRepo, tagRepo)
	startReadingHandler := command.NewStartReadingHandler(uow, bookRepo)
	finishReadingHandler := command.NewFinishReadingHandler(uow, bookRepo, pokemonRepo, gachaSvc)
	updateBookHandler := command.NewUpdateBookHandler(uow, bookRepo, tagRepo)
	deleteBookHandler := command.NewDeleteBookHandler(uow, bookRepo)
	addMemoHandler := command.NewAddMemoHandler(uow, bookRepo, memoRepo)
	updateMemoHandler := command.NewUpdateMemoHandler(uow, memoRepo)
	deleteMemoHandler := command.NewDeleteMemoHandler(uow, memoRepo)
	createTagHandler := command.NewCreateTagHandler(uow, tagRepo)
	deleteTagHandler := command.NewDeleteTagHandler(uow, tagRepo)
	processGitHubEventHandler := command.NewProcessGitHubEventHandler(uow, userRepo, activityRepo, pokemonRepo, gachaSvc)

	// --- Query Handlers ---
	getBooksHandler := query.NewGetBooksHandler(bookRepo)
	getBookDetailHandler := query.NewGetBookDetailHandler(bookRepo, memoRepo)
	getMemosHandler := query.NewGetMemosHandler(memoRepo)
	getTagsHandler := query.NewGetTagsHandler(tagRepo)
	getPokedexHandler := query.NewGetPokedexHandler(pokemonRepo)
	getTeamRankingHandler := query.NewGetTeamRankingHandler(teamRepo, userRepo, pokemonRepo, activityRepo)
	getActivitiesHandler := query.NewGetGitHubActivitiesHandler(activityRepo)
	getUserStatsHandler := query.NewGetUserStatsHandler(userRepo, activityRepo, pokemonRepo)

	// --- Command Handlers (cont.) ---
	chooseStarterHandler := command.NewChooseStarterHandler(uow, pokemonRepo, gachaSvc)
	createTeamHandler := command.NewCreateTeamHandler(uow, teamRepo)
	joinTeamHandler := command.NewJoinTeamHandler(uow, userRepo, teamRepo)

	// --- Presentation Handlers ---
	authHandler := handler.NewAuthHandler(googleOAuth, githubOAuth, jwtManager, userRepo, tenantRepo)
	bookHandler := handler.NewBookHandler(
		registerBookHandler, startReadingHandler, finishReadingHandler,
		updateBookHandler, deleteBookHandler, getBooksHandler, getBookDetailHandler,
	)
	memoHandler := handler.NewMemoHandler(addMemoHandler, updateMemoHandler, deleteMemoHandler, getMemosHandler)
	tagHandler := handler.NewTagHandler(createTagHandler, deleteTagHandler, getTagsHandler)
	pokedexHandler := handler.NewPokedexHandler(getPokedexHandler)
	starterHandler := handler.NewStarterHandler(chooseStarterHandler, getPokedexHandler)
	teamHandler := handler.NewTeamHandler(createTeamHandler, joinTeamHandler, getTeamRankingHandler)
	webhookHandler := handler.NewWebhookHandler(webhookVerifier, processGitHubEventHandler)
	activityHandler := handler.NewActivityHandler(getActivitiesHandler, getUserStatsHandler)

	// --- Router ---
	r := gin.Default()
	router.Setup(r, jwtManager, authHandler, bookHandler, memoHandler, tagHandler, pokedexHandler, starterHandler, teamHandler, webhookHandler, activityHandler)

	// --- Start ---
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

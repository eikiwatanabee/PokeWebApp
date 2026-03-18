package router

import (
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/infrastructure/auth"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/presentation/handler"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/presentation/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(
	r *gin.Engine,
	jwtManager *auth.JWTManager,
	authHandler *handler.AuthHandler,
	bookHandler *handler.BookHandler,
	memoHandler *handler.MemoHandler,
	tagHandler *handler.TagHandler,
	pokedexHandler *handler.PokedexHandler,
	starterHandler *handler.StarterHandler,
	teamHandler *handler.TeamHandler,
	webhookHandler *handler.WebhookHandler,
	activityHandler *handler.ActivityHandler,
	dailyMissionHandler *handler.DailyMissionHandler,
	rankingHandler *handler.RankingHandler,
	feedHandler *handler.FeedHandler,
	weeklyEventHandler *handler.WeeklyEventHandler,
	limitedEventHandler *handler.LimitedEventHandler,
	tradeHandler *handler.TradeHandler,
) {
	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth (public)
	authGroup := api.Group("/auth")
	{
		authGroup.GET("/google", authHandler.GoogleLogin)
		authGroup.GET("/google/callback", authHandler.GoogleCallback)
		authGroup.GET("/github", authHandler.GitHubLogin)
		authGroup.GET("/github/callback", authHandler.GitHubCallback)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/dev-login", authHandler.DevLogin)
	}

	// GitHub Webhook (public, verified by signature)
	api.POST("/webhook/github", webhookHandler.HandleGitHubWebhook)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(jwtManager))
	{
		// GitHub Activities
		activities := protected.Group("/activities")
		{
			activities.GET("", activityHandler.GetActivities)
		}

		// User Stats
		protected.GET("/stats", activityHandler.GetStats)

		// Books (kept for backwards compatibility)
		books := protected.Group("/books")
		{
			books.POST("", bookHandler.Register)
			books.GET("", bookHandler.GetBooks)
			books.GET("/:id", bookHandler.GetBookDetail)
			books.PUT("/:id", bookHandler.UpdateBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
			books.POST("/:id/start", bookHandler.StartReading)
			books.POST("/:id/finish", bookHandler.FinishReading)

			// Memos (nested under books)
			books.POST("/:id/memos", memoHandler.AddMemo)
			books.GET("/:id/memos", memoHandler.GetMemos)
		}

		// Memos (standalone for update/delete)
		memos := protected.Group("/memos")
		{
			memos.PUT("/:id", memoHandler.UpdateMemo)
			memos.DELETE("/:id", memoHandler.DeleteMemo)
		}

		// Tags
		tags := protected.Group("/tags")
		{
			tags.POST("", tagHandler.CreateTag)
			tags.GET("", tagHandler.GetTags)
			tags.DELETE("/:id", tagHandler.DeleteTag)
		}

		// Pokedex
		pokedex := protected.Group("/pokedex")
		{
			pokedex.GET("", pokedexHandler.GetPokedex)
		}

		// Starter Pokemon
		starter := protected.Group("/starter")
		{
			starter.GET("/check", starterHandler.NeedsStarter)
			starter.POST("/choose", starterHandler.ChooseStarter)
		}

		// Teams
		teams := protected.Group("/teams")
		{
			teams.POST("", teamHandler.CreateTeam)
			teams.POST("/:id/join", teamHandler.JoinTeam)
			teams.GET("/ranking", teamHandler.GetRanking)
		}

		// Daily Missions & Login Bonus
		daily := protected.Group("/daily")
		{
			daily.GET("/missions", dailyMissionHandler.GetMissions)
			daily.POST("/login-bonus", dailyMissionHandler.ClaimLoginBonus)
		}

		// User Ranking & Trainer Cards
		protected.GET("/ranking/users", rankingHandler.GetUserRanking)
		protected.GET("/trainers/:id", rankingHandler.GetTrainerCard)
		protected.GET("/trainers/me", rankingHandler.GetTrainerCard)

		// Team Feed
		protected.GET("/feed", feedHandler.GetTeamFeed)

		// Weekly Event
		protected.GET("/events/weekly", weeklyEventHandler.GetCurrentEvent)

		// Limited Events
		protected.GET("/events/limited", limitedEventHandler.GetActiveEvents)
		protected.POST("/events/limited", limitedEventHandler.CreateEvent) // TODO: restrict to admin

		// Trades
		trades := protected.Group("/trades")
		{
			trades.GET("", tradeHandler.GetTrades)
			trades.POST("", tradeHandler.CreateTrade)
			trades.POST("/:id/accept", tradeHandler.AcceptTrade)
			trades.POST("/:id/cancel", tradeHandler.CancelTrade)
		}
	}
}

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
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(jwtManager))
	{
		// Books
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
	}
}

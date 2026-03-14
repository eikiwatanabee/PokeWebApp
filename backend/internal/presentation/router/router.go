package router

import (
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/presentation/handler"
	"github.com/gin-gonic/gin"
)

func Setup(
	r *gin.Engine,
	bookHandler *handler.BookHandler,
	pokedexHandler *handler.PokedexHandler,
) {
	api := r.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// TODO: Add auth middleware
	// api.Use(middleware.Auth())

	books := api.Group("/books")
	{
		books.POST("", bookHandler.Register)
		books.GET("", bookHandler.GetBooks)
		books.POST("/:id/finish", bookHandler.FinishReading)
	}

	pokedex := api.Group("/pokedex")
	{
		pokedex.GET("", pokedexHandler.GetPokedex)
	}
}

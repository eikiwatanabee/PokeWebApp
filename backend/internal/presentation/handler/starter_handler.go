package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type StarterHandler struct {
	chooseStarter *command.ChooseStarterHandler
	getPokedex    *query.GetPokedexHandler
}

func NewStarterHandler(cs *command.ChooseStarterHandler, gp *query.GetPokedexHandler) *StarterHandler {
	return &StarterHandler{chooseStarter: cs, getPokedex: gp}
}

func (h *StarterHandler) ChooseStarter(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		PokemonID int `json:"pokemon_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.chooseStarter.Handle(c.Request.Context(), &command.ChooseStarterCommand{
		UserID:    userID,
		PokemonID: req.PokemonID,
	})
	if err != nil {
		switch err {
		case command.ErrInvalidStarter:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case command.ErrAlreadyHasStarter:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *StarterHandler) NeedsStarter(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.getPokedex.Handle(c.Request.Context(), &query.GetPokedexQuery{UserID: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"needs_starter": result.Total == 0})
}

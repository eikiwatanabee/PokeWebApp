package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type PokedexHandler struct {
	getPokedex *query.GetPokedexHandler
}

func NewPokedexHandler(gp *query.GetPokedexHandler) *PokedexHandler {
	return &PokedexHandler{getPokedex: gp}
}

func (h *PokedexHandler) GetPokedex(c *gin.Context) {
	userID := c.GetString("user_id")

	q := &query.GetPokedexQuery{UserID: userID}
	result, err := h.getPokedex.Handle(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

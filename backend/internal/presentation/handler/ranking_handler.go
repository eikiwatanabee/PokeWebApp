package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	getUserRanking *query.GetUserRankingHandler
	getTrainerCard *query.GetTrainerCardHandler
}

func NewRankingHandler(
	gur *query.GetUserRankingHandler,
	gtc *query.GetTrainerCardHandler,
) *RankingHandler {
	return &RankingHandler{
		getUserRanking: gur,
		getTrainerCard: gtc,
	}
}

func (h *RankingHandler) GetUserRanking(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	sortBy := c.DefaultQuery("sort_by", "xp")

	result, err := h.getUserRanking.Handle(c.Request.Context(), &query.GetUserRankingQuery{
		TenantID: tenantID,
		SortBy:   sortBy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *RankingHandler) GetTrainerCard(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" || userID == "me" {
		userID = c.GetString("user_id")
	}

	result, err := h.getTrainerCard.Handle(c.Request.Context(), &query.GetTrainerCardQuery{
		UserID: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

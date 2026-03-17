package handler

import (
	"net/http"
	"strconv"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	getActivities *query.GetGitHubActivitiesHandler
	getUserStats  *query.GetUserStatsHandler
}

func NewActivityHandler(
	ga *query.GetGitHubActivitiesHandler,
	gus *query.GetUserStatsHandler,
) *ActivityHandler {
	return &ActivityHandler{getActivities: ga, getUserStats: gus}
}

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	userID := c.GetString("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	result, err := h.getActivities.Handle(c.Request.Context(), &query.GetGitHubActivitiesQuery{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ActivityHandler) GetStats(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.getUserStats.Handle(c.Request.Context(), &query.GetUserStatsQuery{
		UserID: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

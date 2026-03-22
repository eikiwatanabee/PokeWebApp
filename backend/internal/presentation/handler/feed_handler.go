package handler

import (
	"net/http"
	"strconv"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type FeedHandler struct {
	getTeamFeed *query.GetTeamFeedHandler
}

func NewFeedHandler(getTeamFeed *query.GetTeamFeedHandler) *FeedHandler {
	return &FeedHandler{getTeamFeed: getTeamFeed}
}

func (h *FeedHandler) GetTeamFeed(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	result, err := h.getTeamFeed.Handle(c.Request.Context(), &query.GetTeamFeedQuery{
		TenantID: tenantID,
		Limit:    limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

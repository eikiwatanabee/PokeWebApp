package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type WeeklyEventHandler struct {
	getWeeklyEvent *query.GetWeeklyEventHandler
}

func NewWeeklyEventHandler(getWeeklyEvent *query.GetWeeklyEventHandler) *WeeklyEventHandler {
	return &WeeklyEventHandler{getWeeklyEvent: getWeeklyEvent}
}

func (h *WeeklyEventHandler) GetCurrentEvent(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	result, err := h.getWeeklyEvent.Handle(c.Request.Context(), &query.GetWeeklyEventQuery{
		TenantID: tenantID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

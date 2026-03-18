package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type DailyMissionHandler struct {
	getDailyMissions  *query.GetDailyMissionsHandler
	claimLoginBonus   *command.ClaimLoginBonusHandler
}

func NewDailyMissionHandler(
	gdm *query.GetDailyMissionsHandler,
	clb *command.ClaimLoginBonusHandler,
) *DailyMissionHandler {
	return &DailyMissionHandler{
		getDailyMissions: gdm,
		claimLoginBonus:  clb,
	}
}

func (h *DailyMissionHandler) GetMissions(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.getDailyMissions.Handle(c.Request.Context(), &query.GetDailyMissionsQuery{
		UserID: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *DailyMissionHandler) ClaimLoginBonus(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.claimLoginBonus.Handle(c.Request.Context(), &command.ClaimLoginBonusCommand{
		UserID: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

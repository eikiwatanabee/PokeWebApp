package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	createTeam     *command.CreateTeamHandler
	joinTeam       *command.JoinTeamHandler
	getTeamRanking *query.GetTeamRankingHandler
}

func NewTeamHandler(ct *command.CreateTeamHandler, jt *command.JoinTeamHandler, gtr *query.GetTeamRankingHandler) *TeamHandler {
	return &TeamHandler{createTeam: ct, joinTeam: jt, getTeamRanking: gtr}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.createTeam.Handle(c.Request.Context(), &command.CreateTeamCommand{
		TenantID: tenantID,
		Name:     req.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *TeamHandler) JoinTeam(c *gin.Context) {
	userID := c.GetString("user_id")
	teamID := c.Param("id")

	err := h.joinTeam.Handle(c.Request.Context(), &command.JoinTeamCommand{
		UserID: userID,
		TeamID: teamID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined team"})
}

func (h *TeamHandler) GetRanking(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	result, err := h.getTeamRanking.Handle(c.Request.Context(), &query.GetTeamRankingQuery{
		TenantID: tenantID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeployerHandler struct {
	deployerRepo repository.DeployerRepository
}

func NewDeployerHandler(deployerRepo repository.DeployerRepository) *DeployerHandler {
	return &DeployerHandler{deployerRepo: deployerRepo}
}

type addDeployerRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *DeployerHandler) AddDeployer(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	var req addDeployerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	if err := h.deployerRepo.AddDeployer(c.Request.Context(), tenantID.(uuid.UUID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deployer added"})
}

func (h *DeployerHandler) RemoveDeployer(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	if err := h.deployerRepo.RemoveDeployer(c.Request.Context(), tenantID.(uuid.UUID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deployer removed"})
}

func (h *DeployerHandler) GetDeployers(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	ids, err := h.deployerRepo.FindByTenantID(c.Request.Context(), tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deployer_ids": ids})
}

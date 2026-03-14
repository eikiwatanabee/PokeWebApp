package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	createTag *command.CreateTagHandler
	deleteTag *command.DeleteTagHandler
	getTags   *query.GetTagsHandler
}

func NewTagHandler(
	ct *command.CreateTagHandler,
	dt *command.DeleteTagHandler,
	gt *query.GetTagsHandler,
) *TagHandler {
	return &TagHandler{createTag: ct, deleteTag: dt, getTags: gt}
}

type createTagRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var req createTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := c.GetString("tenant_id")
	cmd := &command.CreateTagCommand{TenantID: tenantID, Name: req.Name}
	result, err := h.createTag.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *TagHandler) GetTags(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	q := &query.GetTagsQuery{TenantID: tenantID}
	result, err := h.getTags.Handle(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	tagID := c.Param("id")
	cmd := &command.DeleteTagCommand{TagID: tagID}
	if err := h.deleteTag.Handle(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag deleted"})
}

package handler

import (
	"net/http"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type MemoHandler struct {
	addMemo    *command.AddMemoHandler
	updateMemo *command.UpdateMemoHandler
	deleteMemo *command.DeleteMemoHandler
	getMemos   *query.GetMemosHandler
}

func NewMemoHandler(
	am *command.AddMemoHandler,
	um *command.UpdateMemoHandler,
	dm *command.DeleteMemoHandler,
	gm *query.GetMemosHandler,
) *MemoHandler {
	return &MemoHandler{addMemo: am, updateMemo: um, deleteMemo: dm, getMemos: gm}
}

type addMemoRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *MemoHandler) AddMemo(c *gin.Context) {
	var req addMemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bookID := c.Param("id")
	cmd := &command.AddMemoCommand{BookID: bookID, Content: req.Content}
	result, err := h.addMemo.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *MemoHandler) GetMemos(c *gin.Context) {
	bookID := c.Param("id")
	q := &query.GetMemosQuery{BookID: bookID}
	result, err := h.getMemos.Handle(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type updateMemoRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *MemoHandler) UpdateMemo(c *gin.Context) {
	var req updateMemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	memoID := c.Param("id")
	cmd := &command.UpdateMemoCommand{MemoID: memoID, Content: req.Content}
	if err := h.updateMemo.Handle(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "memo updated"})
}

func (h *MemoHandler) DeleteMemo(c *gin.Context) {
	memoID := c.Param("id")
	cmd := &command.DeleteMemoCommand{MemoID: memoID}
	if err := h.deleteMemo.Handle(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "memo deleted"})
}

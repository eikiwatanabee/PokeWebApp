package handler

import (
	"net/http"
	"strconv"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/command"
	"github.com/eikiwatanabee/PokeWebApp/backend/internal/application/query"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	registerBook  *command.RegisterBookHandler
	finishReading *command.FinishReadingHandler
	getBooks      *query.GetBooksHandler
}

func NewBookHandler(
	rb *command.RegisterBookHandler,
	fr *command.FinishReadingHandler,
	gb *query.GetBooksHandler,
) *BookHandler {
	return &BookHandler{registerBook: rb, finishReading: fr, getBooks: gb}
}

type registerBookRequest struct {
	Title  string   `json:"title" binding:"required"`
	Author string   `json:"author" binding:"required"`
	TagIDs []string `json:"tag_ids"`
}

func (h *BookHandler) Register(c *gin.Context) {
	var req registerBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	cmd := &command.RegisterBookCommand{
		UserID: userID,
		Title:  req.Title,
		Author: req.Author,
		TagIDs: req.TagIDs,
	}

	result, err := h.registerBook.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *BookHandler) FinishReading(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("id")

	cmd := &command.FinishReadingCommand{BookID: bookID, UserID: userID}
	result, err := h.finishReading.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book finished! You caught a Pokemon!",
		"pokemon": result,
	})
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	q := &query.GetBooksQuery{
		UserID: userID,
		Status: c.Query("status"),
		TagID:  c.Query("tag_id"),
		Page:   page,
		Limit:  limit,
	}

	result, err := h.getBooks.Handle(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

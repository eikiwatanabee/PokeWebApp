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
	startReading  *command.StartReadingHandler
	finishReading *command.FinishReadingHandler
	updateBook    *command.UpdateBookHandler
	deleteBook    *command.DeleteBookHandler
	getBooks      *query.GetBooksHandler
	getBookDetail *query.GetBookDetailHandler
}

func NewBookHandler(
	rb *command.RegisterBookHandler,
	sr *command.StartReadingHandler,
	fr *command.FinishReadingHandler,
	ub *command.UpdateBookHandler,
	db *command.DeleteBookHandler,
	gb *query.GetBooksHandler,
	gbd *query.GetBookDetailHandler,
) *BookHandler {
	return &BookHandler{
		registerBook: rb, startReading: sr, finishReading: fr,
		updateBook: ub, deleteBook: db, getBooks: gb, getBookDetail: gbd,
	}
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
		UserID: userID, Title: req.Title, Author: req.Author, TagIDs: req.TagIDs,
	}
	result, err := h.registerBook.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	q := &query.GetBooksQuery{
		UserID: userID, Status: c.Query("status"),
		TagID: c.Query("tag_id"), Page: page, Limit: limit,
	}
	result, err := h.getBooks.Handle(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *BookHandler) GetBookDetail(c *gin.Context) {
	bookID := c.Param("id")
	q := &query.GetBookDetailQuery{BookID: bookID}
	result, err := h.getBookDetail.Handle(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *BookHandler) StartReading(c *gin.Context) {
	bookID := c.Param("id")
	cmd := &command.StartReadingCommand{BookID: bookID}
	if err := h.startReading.Handle(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "started reading"})
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

type updateBookRequest struct {
	Title  string   `json:"title"`
	Author string   `json:"author"`
	TagIDs []string `json:"tag_ids"`
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	var req updateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bookID := c.Param("id")
	cmd := &command.UpdateBookCommand{
		BookID: bookID, Title: req.Title, Author: req.Author, TagIDs: req.TagIDs,
	}
	if err := h.updateBook.Handle(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book updated"})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	bookID := c.Param("id")
	cmd := &command.DeleteBookCommand{BookID: bookID}
	if err := h.deleteBook.Handle(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}

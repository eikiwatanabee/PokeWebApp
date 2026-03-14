---
name: add-handler
description: Generate a Gin HTTP handler following Clean Architecture presentation layer patterns. Use when adding new API endpoints.
argument-hint: [HandlerName]
---

# Gin Handler Generator

Generate a new HTTP handler for the PokeBookManager project.

## Handler Name: $ARGUMENTS

## Instructions

1. Read `docs/REQUIREMENTS.md` for API endpoint specs
2. Generate handler in `backend/internal/presentation/handler/`
3. Register routes in `backend/internal/presentation/router/router.go`

### Pattern

```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/your-org/pokebookmanager/internal/application/command"
    "github.com/your-org/pokebookmanager/internal/application/query"
)

type BookHandler struct {
    finishReadingHandler *command.FinishReadingHandler
    getBooksHandler      *query.GetBooksHandler
}

func NewBookHandler(
    frh *command.FinishReadingHandler,
    gbh *query.GetBooksHandler,
) *BookHandler {
    return &BookHandler{finishReadingHandler: frh, getBooksHandler: gbh}
}

// POST /api/books/:id/finish
func (h *BookHandler) FinishReading(c *gin.Context) {
    userID := c.GetString("user_id") // from auth middleware
    bookID := c.Param("id")

    cmd := &command.FinishReadingCommand{
        BookID: bookID,
        UserID: userID,
    }

    result, err := h.finishReadingHandler.Handle(c.Request.Context(), cmd)
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Book finished! You caught a Pokemon!",
        "pokemon": result,
    })
}

// GET /api/books
func (h *BookHandler) GetBooks(c *gin.Context) {
    userID := c.GetString("user_id")

    q := &query.GetBooksQuery{
        UserID: userID,
        Status: c.Query("status"),
        TagID:  c.Query("tag_id"),
        Page:   getPageParam(c),
        Limit:  getLimitParam(c),
    }

    result, err := h.getBooksHandler.Handle(c.Request.Context(), q)
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, result)
}
```

## Rules
- Handlers are thin: parse request → call command/query → return response
- NO business logic in handlers
- Use auth middleware context for user_id, tenant_id
- Consistent error handling with handleError helper
- JSON response format

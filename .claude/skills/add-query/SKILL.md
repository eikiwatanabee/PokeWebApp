---
name: add-query
description: Generate a CQRS Query with handler for read operations. Use when adding new read/fetch functionality.
argument-hint: [QueryName]
---

# CQRS Query Generator

Generate a new Query (read side of CQRS) for the PokeBookManager project.

## Query Name: $ARGUMENTS

## Instructions

1. Read `docs/REQUIREMENTS.md` for domain context
2. Generate the following files:

### Query DTO (`backend/internal/application/query/$ARGUMENTS.go`)
- Input struct (filter/pagination params)
- Response DTO struct
- Handler struct with read-only dependencies
- `Handle(ctx context.Context, q *{QueryName}Query) (*{Result}, error)` method

### Pattern

```go
package query

import (
    "context"
    "github.com/your-org/pokebookmanager/internal/domain/repository"
)

type GetBooksQuery struct {
    UserID string
    Status string // optional filter
    TagID  string // optional filter
    Page   int
    Limit  int
}

type BookDTO struct {
    ID       string   `json:"id"`
    Title    string   `json:"title"`
    Author   string   `json:"author"`
    Status   string   `json:"status"`
    Tags     []string `json:"tags"`
}

type GetBooksResult struct {
    Books      []BookDTO `json:"books"`
    TotalCount int       `json:"total_count"`
    Page       int       `json:"page"`
}

type GetBooksHandler struct {
    bookRepo repository.BookRepository
}

func (h *GetBooksHandler) Handle(ctx context.Context, q *GetBooksQuery) (*GetBooksResult, error) {
    // Read-only: no UoW, no locks
    books, total, err := h.bookRepo.FindByUserWithFilters(ctx, q.UserID, q.Status, q.TagID, q.Page, q.Limit)
    if err != nil {
        return nil, err
    }
    // Map domain entities to DTOs
    dtos := make([]BookDTO, len(books))
    for i, b := range books {
        dtos[i] = BookDTO{
            ID: b.ID.String(), Title: b.Title,
            Author: b.Author, Status: string(b.Status),
        }
    }
    return &GetBooksResult{Books: dtos, TotalCount: total, Page: q.Page}, nil
}
```

## Rules
- Queries are read-only, NEVER modify data
- No UoW needed (no transactions for reads)
- No SELECT ... FOR UPDATE (no locks)
- Return DTOs, not domain entities (prevent domain leakage)
- Support pagination where applicable

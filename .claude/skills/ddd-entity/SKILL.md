---
name: ddd-entity
description: Generate DDD Entity, Value Object, and Repository interface following Clean Architecture. Use when creating new domain models.
argument-hint: [EntityName]
---

# DDD Entity Generator

Generate a new DDD Entity for the PokeBookManager project.

## Entity Name: $ARGUMENTS

## Instructions

1. Read `docs/REQUIREMENTS.md` for domain context
2. Generate the following files:

### Entity (`backend/internal/domain/entity/$ARGUMENTS.go`)
- UUID as ID (use `github.com/google/uuid`)
- Proper field types and validation
- Constructor function `New{EntityName}(...) (*{EntityName}, error)`
- Domain methods with business logic
- No GORM tags (domain layer is infrastructure-agnostic)

### Repository Interface (`backend/internal/domain/repository/$ARGUMENTS_repository.go`)
- Define in domain layer (interface only)
- Methods: FindByID, FindAll, Save, Delete
- Use context.Context as first parameter
- Return domain entities, not GORM models

### Value Objects (if needed)
- Place in `backend/internal/domain/valueobject/`
- Immutable, equality by value
- Validation in constructor

## Architecture Rules
- Domain layer has ZERO external dependencies (no Gin, no GORM)
- Repository interfaces defined in domain, implemented in infrastructure
- Use domain errors, not infrastructure errors
- Follow lock order: Tenant(1) → User(2) → Book(3) → Memo(4) → Tag(5) → UserPokemon(6)

## Example Pattern

```go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type Book struct {
    ID         uuid.UUID
    UserID     uuid.UUID
    Title      string
    Author     string
    Status     valueobject.BookStatus
    FinishedAt *time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

func NewBook(userID uuid.UUID, title, author string) (*Book, error) {
    if title == "" {
        return nil, ErrEmptyTitle
    }
    return &Book{
        ID:        uuid.New(),
        UserID:    userID,
        Title:     title,
        Author:    author,
        Status:    valueobject.Unread,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}
```

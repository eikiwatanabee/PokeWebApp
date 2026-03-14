---
name: db-migrate
description: Generate GORM database migration code. Use when modifying database schema.
argument-hint: [migration-description]
---

# DB Migration Generator

Generate migration for: $ARGUMENTS

## Instructions

1. Review current schema in `docs/REQUIREMENTS.md` section 8
2. Generate GORM AutoMigrate compatible model changes
3. Handle both up and down migrations

## GORM Model Convention

```go
package persistence

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Infrastructure model (NOT domain entity)
type BookModel struct {
    ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID     uuid.UUID  `gorm:"type:uuid;not null;index"`
    Title      string     `gorm:"type:varchar(255);not null"`
    Author     string     `gorm:"type:varchar(255);not null"`
    Status     string     `gorm:"type:varchar(50);not null;default:'unread'"`
    FinishedAt *time.Time `gorm:"type:timestamp"`
    CreatedAt  time.Time  `gorm:"autoCreateTime"`
    UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
}

func (BookModel) TableName() string {
    return "books"
}
```

## Rules
- Infrastructure models are SEPARATE from domain entities
- Include proper indexes for query performance
- Use UUID for all primary keys
- Use JSONB for flexible fields (pokemon_types)
- Add foreign key constraints
- Consider lock order when adding new tables

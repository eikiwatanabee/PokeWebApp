package event

import (
	"time"

	"github.com/google/uuid"
)

// BookFinished is a domain event raised when a user finishes reading a book.
type BookFinished struct {
	BookID     uuid.UUID
	UserID     uuid.UUID
	OccurredAt time.Time
}

func NewBookFinished(bookID, userID uuid.UUID) BookFinished {
	return BookFinished{
		BookID:     bookID,
		UserID:     userID,
		OccurredAt: time.Now(),
	}
}

package entity

import (
	"time"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
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
	Tags       []Tag
	Memos      []Memo
}

func NewBook(userID uuid.UUID, title, author string) (*Book, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if author == "" {
		return nil, ErrEmptyAuthor
	}
	now := time.Now()
	return &Book{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Author:    author,
		Status:    valueobject.Unread,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (b *Book) StartReading() error {
	if !b.Status.CanTransitionTo(valueobject.Reading) {
		return ErrInvalidTransition
	}
	b.Status = valueobject.Reading
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Book) Finish() error {
	if !b.Status.CanTransitionTo(valueobject.Finished) {
		return ErrNotReading
	}
	now := time.Now()
	b.Status = valueobject.Finished
	b.FinishedAt = &now
	b.UpdatedAt = now
	return nil
}

func (b *Book) IsFinished() bool {
	return b.Status == valueobject.Finished
}

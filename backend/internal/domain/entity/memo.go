package entity

import (
	"time"

	"github.com/google/uuid"
)

type Memo struct {
	ID        uuid.UUID
	BookID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMemo(bookID uuid.UUID, content string) (*Memo, error) {
	if content == "" {
		return nil, ErrEmptyMemoContent
	}
	now := time.Now()
	return &Memo{
		ID:        uuid.New(),
		BookID:    bookID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (m *Memo) UpdateContent(content string) error {
	if content == "" {
		return ErrEmptyMemoContent
	}
	m.Content = content
	m.UpdatedAt = time.Now()
	return nil
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type CatalogBook struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Title     string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCatalogBook(tenantID uuid.UUID, title, author string) (*CatalogBook, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	if author == "" {
		return nil, ErrEmptyAuthor
	}
	now := time.Now()
	return &CatalogBook{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Title:     title,
		Author:    author,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (cb *CatalogBook) Update(title, author string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if author == "" {
		return ErrEmptyAuthor
	}
	cb.Title = title
	cb.Author = author
	cb.UpdatedAt = time.Now()
	return nil
}

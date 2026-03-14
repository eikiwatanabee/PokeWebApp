package entity

import (
	"testing"

	"github.com/eikiwatanabee/PokeWebApp/backend/internal/domain/valueobject"
	"github.com/google/uuid"
)

func TestNewBook(t *testing.T) {
	tests := []struct {
		name    string
		userID  uuid.UUID
		title   string
		author  string
		wantErr error
	}{
		{
			name:    "valid book",
			userID:  uuid.New(),
			title:   "Clean Architecture",
			author:  "Robert C. Martin",
			wantErr: nil,
		},
		{
			name:    "empty title",
			userID:  uuid.New(),
			title:   "",
			author:  "Robert C. Martin",
			wantErr: ErrEmptyTitle,
		},
		{
			name:    "empty author",
			userID:  uuid.New(),
			title:   "Clean Architecture",
			author:  "",
			wantErr: ErrEmptyAuthor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := NewBook(tt.userID, tt.title, tt.author)
			if err != tt.wantErr {
				t.Errorf("NewBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if book.ID == uuid.Nil {
					t.Error("NewBook() ID should not be nil")
				}
				if book.UserID != tt.userID {
					t.Errorf("NewBook() UserID = %v, want %v", book.UserID, tt.userID)
				}
				if book.Title != tt.title {
					t.Errorf("NewBook() Title = %v, want %v", book.Title, tt.title)
				}
				if book.Author != tt.author {
					t.Errorf("NewBook() Author = %v, want %v", book.Author, tt.author)
				}
				if book.Status != valueobject.Unread {
					t.Errorf("NewBook() Status = %v, want %v", book.Status, valueobject.Unread)
				}
				if book.FinishedAt != nil {
					t.Error("NewBook() FinishedAt should be nil")
				}
			}
		})
	}
}

func TestBook_StartReading(t *testing.T) {
	tests := []struct {
		name    string
		status  valueobject.BookStatus
		wantErr bool
	}{
		{"from unread", valueobject.Unread, false},
		{"from reading", valueobject.Reading, true},
		{"from finished", valueobject.Finished, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := &Book{Status: tt.status}
			err := book.StartReading()
			if (err != nil) != tt.wantErr {
				t.Errorf("StartReading() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && book.Status != valueobject.Reading {
				t.Errorf("StartReading() Status = %v, want %v", book.Status, valueobject.Reading)
			}
		})
	}
}

func TestBook_Finish(t *testing.T) {
	tests := []struct {
		name    string
		status  valueobject.BookStatus
		wantErr bool
	}{
		{"from reading", valueobject.Reading, false},
		{"from unread", valueobject.Unread, true},
		{"from finished", valueobject.Finished, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := &Book{Status: tt.status}
			err := book.Finish()
			if (err != nil) != tt.wantErr {
				t.Errorf("Finish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if book.Status != valueobject.Finished {
					t.Errorf("Finish() Status = %v, want %v", book.Status, valueobject.Finished)
				}
				if book.FinishedAt == nil {
					t.Error("Finish() FinishedAt should not be nil")
				}
			}
		})
	}
}

func TestBook_IsFinished(t *testing.T) {
	tests := []struct {
		name   string
		status valueobject.BookStatus
		want   bool
	}{
		{"unread", valueobject.Unread, false},
		{"reading", valueobject.Reading, false},
		{"finished", valueobject.Finished, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := &Book{Status: tt.status}
			if got := book.IsFinished(); got != tt.want {
				t.Errorf("IsFinished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBook_StatusTransitionFlow(t *testing.T) {
	book, err := NewBook(uuid.New(), "DDD", "Eric Evans")
	if err != nil {
		t.Fatalf("NewBook() error = %v", err)
	}

	if err := book.StartReading(); err != nil {
		t.Fatalf("StartReading() error = %v", err)
	}
	if book.Status != valueobject.Reading {
		t.Fatalf("expected Reading, got %v", book.Status)
	}

	if err := book.Finish(); err != nil {
		t.Fatalf("Finish() error = %v", err)
	}
	if book.Status != valueobject.Finished {
		t.Fatalf("expected Finished, got %v", book.Status)
	}

	if err := book.StartReading(); err == nil {
		t.Error("should not be able to start reading a finished book")
	}
}

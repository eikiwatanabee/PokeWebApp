package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewMemo(t *testing.T) {
	tests := []struct {
		name    string
		bookID  uuid.UUID
		content string
		wantErr error
	}{
		{
			name:    "valid memo",
			bookID:  uuid.New(),
			content: "Great chapter about DDD aggregates",
			wantErr: nil,
		},
		{
			name:    "empty content",
			bookID:  uuid.New(),
			content: "",
			wantErr: ErrEmptyMemoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memo, err := NewMemo(tt.bookID, tt.content)
			if err != tt.wantErr {
				t.Errorf("NewMemo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if memo.ID == uuid.Nil {
					t.Error("NewMemo() ID should not be nil")
				}
				if memo.BookID != tt.bookID {
					t.Errorf("NewMemo() BookID = %v, want %v", memo.BookID, tt.bookID)
				}
				if memo.Content != tt.content {
					t.Errorf("NewMemo() Content = %v, want %v", memo.Content, tt.content)
				}
			}
		})
	}
}

func TestMemo_UpdateContent(t *testing.T) {
	memo := &Memo{Content: "original"}

	if err := memo.UpdateContent("updated"); err != nil {
		t.Errorf("UpdateContent() error = %v", err)
	}
	if memo.Content != "updated" {
		t.Errorf("UpdateContent() Content = %v, want %v", memo.Content, "updated")
	}

	if err := memo.UpdateContent(""); err != ErrEmptyMemoContent {
		t.Errorf("UpdateContent() error = %v, wantErr %v", err, ErrEmptyMemoContent)
	}
}

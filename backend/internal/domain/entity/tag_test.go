package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTag(t *testing.T) {
	tests := []struct {
		name     string
		tenantID uuid.UUID
		tagName  string
		wantErr  error
	}{
		{
			name:     "valid tag",
			tenantID: uuid.New(),
			tagName:  "技術書",
			wantErr:  nil,
		},
		{
			name:     "empty name",
			tenantID: uuid.New(),
			tagName:  "",
			wantErr:  ErrEmptyTagName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, err := NewTag(tt.tenantID, tt.tagName)
			if err != tt.wantErr {
				t.Errorf("NewTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if tag.ID == uuid.Nil {
					t.Error("NewTag() ID should not be nil")
				}
				if tag.TenantID != tt.tenantID {
					t.Errorf("NewTag() TenantID = %v, want %v", tag.TenantID, tt.tenantID)
				}
				if tag.Name != tt.tagName {
					t.Errorf("NewTag() Name = %v, want %v", tag.Name, tt.tagName)
				}
			}
		})
	}
}

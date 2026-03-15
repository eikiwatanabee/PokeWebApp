package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		tenantID uuid.UUID
		googleID string
		email    string
		userName string
		wantErr  error
	}{
		{
			name:     "valid user",
			tenantID: uuid.New(),
			googleID: "google-123",
			email:    "test@example.com",
			userName: "Test User",
			wantErr:  nil,
		},
		{
			name:     "empty google ID",
			tenantID: uuid.New(),
			googleID: "",
			email:    "test@example.com",
			userName: "Test User",
			wantErr:  ErrEmptyGoogleID,
		},
		{
			name:     "empty email",
			tenantID: uuid.New(),
			googleID: "google-123",
			email:    "",
			userName: "Test User",
			wantErr:  ErrEmptyEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.tenantID, tt.googleID, tt.email, tt.userName)
			if err != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if user.ID == uuid.Nil {
					t.Error("NewUser() ID should not be nil")
				}
				if user.Role != RoleMember {
					t.Errorf("NewUser() Role = %v, want %v", user.Role, RoleMember)
				}
				if user.TenantID != tt.tenantID {
					t.Errorf("NewUser() TenantID = %v, want %v", user.TenantID, tt.tenantID)
				}
			}
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	user := &User{Role: RoleMember}
	if user.IsAdmin() {
		t.Error("member should not be admin")
	}

	user.Role = RoleAdmin
	if !user.IsAdmin() {
		t.Error("admin should be admin")
	}
}

func TestUser_PromoteToAdmin(t *testing.T) {
	user := &User{Role: RoleMember}
	user.PromoteToAdmin()

	if user.Role != RoleAdmin {
		t.Errorf("PromoteToAdmin() Role = %v, want %v", user.Role, RoleAdmin)
	}
}

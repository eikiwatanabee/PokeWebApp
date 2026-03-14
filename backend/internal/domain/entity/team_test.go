package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTeam(t *testing.T) {
	tests := []struct {
		name     string
		tenantID uuid.UUID
		teamName string
		wantErr  error
	}{
		{
			name:     "valid team",
			tenantID: uuid.New(),
			teamName: "Team Rocket",
			wantErr:  nil,
		},
		{
			name:     "empty name",
			tenantID: uuid.New(),
			teamName: "",
			wantErr:  ErrEmptyTeamName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, err := NewTeam(tt.tenantID, tt.teamName)
			if err != tt.wantErr {
				t.Errorf("NewTeam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if team.ID == uuid.Nil {
					t.Error("NewTeam() ID should not be nil")
				}
				if team.TenantID != tt.tenantID {
					t.Errorf("NewTeam() TenantID = %v, want %v", team.TenantID, tt.tenantID)
				}
				if team.Name != tt.teamName {
					t.Errorf("NewTeam() Name = %v, want %v", team.Name, tt.teamName)
				}
			}
		})
	}
}

func TestUser_JoinTeam(t *testing.T) {
	user := &User{Role: RoleMember}
	teamID := uuid.New()
	user.JoinTeam(teamID)

	if user.TeamID == nil {
		t.Error("JoinTeam() TeamID should not be nil")
	}
	if *user.TeamID != teamID {
		t.Errorf("JoinTeam() TeamID = %v, want %v", *user.TeamID, teamID)
	}
}

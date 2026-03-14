package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTenant(t *testing.T) {
	tests := []struct {
		name       string
		tenantName string
		wantErr    error
	}{
		{
			name:       "valid tenant",
			tenantName: "Acme Corp",
			wantErr:    nil,
		},
		{
			name:       "empty name",
			tenantName: "",
			wantErr:    ErrEmptyTenantName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenant, err := NewTenant(tt.tenantName)
			if err != tt.wantErr {
				t.Errorf("NewTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if tenant.ID == uuid.Nil {
					t.Error("NewTenant() ID should not be nil")
				}
				if tenant.Name != tt.tenantName {
					t.Errorf("NewTenant() Name = %v, want %v", tenant.Name, tt.tenantName)
				}
			}
		})
	}
}

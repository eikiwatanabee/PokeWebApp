package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
)

type User struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	GoogleID  string
	Email     string
	Name      string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(tenantID uuid.UUID, googleID, email, name string) (*User, error) {
	if googleID == "" {
		return nil, ErrEmptyGoogleID
	}
	if email == "" {
		return nil, ErrEmptyEmail
	}
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		GoogleID:  googleID,
		Email:     email,
		Name:      name,
		Role:      RoleMember,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) PromoteToAdmin() {
	u.Role = RoleAdmin
	u.UpdatedAt = time.Now()
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type AdminRole string
type AdminStatus string

const (
	// Admin roles
	RoleSuperAdmin   AdminRole = "super_admin"
	RoleOpsAdmin     AdminRole = "ops_admin"
	RoleFinanceAdmin AdminRole = "finance_admin"

	// Admin status
	AdminStatusActive    AdminStatus = "active"
	AdminStatusSuspended AdminStatus = "suspended"
	AdminStatusInactive  AdminStatus = "inactive"
)

type Admin struct {
	ID           uuid.UUID   `db:"id" json:"id"`
	Username     string      `db:"username" json:"username"`
	Email        string      `db:"email" json:"email"`
	PasswordHash string      `db:"password_hash" json:"-"`
	FullName     string      `db:"full_name" json:"full_name"`
	Role         AdminRole   `db:"role" json:"role"`
	Status       AdminStatus `db:"status" json:"status"`
	LastLoginAt  *time.Time  `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedBy    *uuid.UUID  `db:"created_by" json:"created_by,omitempty"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
}

// IsActive checks if admin is active
func (a *Admin) IsActive() bool {
	return a.Status == AdminStatusActive
}

// HasPermission checks if admin has specific permission
func (a *Admin) HasPermission(requiredRole AdminRole) bool {
	// Super admin has all permissions
	if a.Role == RoleSuperAdmin {
		return true
	}
	// Exact role match
	return a.Role == requiredRole
}

// CanManageAdmins checks if admin can manage other admins
func (a *Admin) CanManageAdmins() bool {
	return a.Role == RoleSuperAdmin
}

// CanRefund checks if admin can perform refunds
func (a *Admin) CanRefund() bool {
	return a.Role == RoleSuperAdmin || a.Role == RoleFinanceAdmin
}

// CanFreezeWallet checks if admin can freeze wallets
func (a *Admin) CanFreezeWallet() bool {
	return a.Role == RoleSuperAdmin || a.Role == RoleOpsAdmin
}

// Validate validates admin data
func (a *Admin) Validate() error {
	if a.Username == "" {
		return ErrInvalidUsername
	}
	if a.Email == "" {
		return ErrInvalidEmail
	}
	if a.FullName == "" {
		return ErrInvalidFullName
	}
	return nil
}

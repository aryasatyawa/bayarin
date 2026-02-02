package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/aryasatyawa/bayarin/internal/config"
	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/pkg/crypto"
	"github.com/aryasatyawa/bayarin/internal/pkg/jwt"
	"github.com/aryasatyawa/bayarin/internal/pkg/validator"
	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AdminUsecase interface {
	Login(ctx context.Context, req AdminLoginRequest) (*AdminLoginResponse, error)
	CreateAdmin(ctx context.Context, creatorID uuid.UUID, req CreateAdminRequest) (*AdminResponse, error)
	GetAdminByID(ctx context.Context, id uuid.UUID) (*AdminResponse, error)
	ListAdmins(ctx context.Context, limit, offset int) ([]*AdminResponse, error)
	UpdateAdmin(ctx context.Context, id uuid.UUID, req UpdateAdminRequest) error
	UpdateAdminStatus(ctx context.Context, actorID, targetID uuid.UUID, status domain.AdminStatus) error
	GetAuditLogs(ctx context.Context, adminID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error)
}

type adminUsecase struct {
	db           *sqlx.DB
	adminRepo    repository.AdminRepository
	auditLogRepo repository.AuditLogRepository
	tokenManager *jwt.TokenManager
	cfg          *config.Config
}

func NewAdminUsecase(
	db *sqlx.DB,
	adminRepo repository.AdminRepository,
	auditLogRepo repository.AuditLogRepository,
	tokenManager *jwt.TokenManager,
	cfg *config.Config,
) AdminUsecase {
	return &adminUsecase{
		db:           db,
		adminRepo:    adminRepo,
		auditLogRepo: auditLogRepo,
		tokenManager: tokenManager,
		cfg:          cfg,
	}
}

// DTOs
type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	AdminID  uuid.UUID        `json:"admin_id"`
	Username string           `json:"username"`
	FullName string           `json:"full_name"`
	Role     domain.AdminRole `json:"role"`
	Token    string           `json:"token"`
}

type CreateAdminRequest struct {
	Username string           `json:"username" validate:"required,min=3,max=100"`
	Email    string           `json:"email" validate:"required,email"`
	Password string           `json:"password" validate:"required,min=8"`
	FullName string           `json:"full_name" validate:"required,min=3,max=255"`
	Role     domain.AdminRole `json:"role" validate:"required"`
}

type UpdateAdminRequest struct {
	Username string           `json:"username" validate:"omitempty,min=3,max=100"`
	Email    string           `json:"email" validate:"omitempty,email"`
	FullName string           `json:"full_name" validate:"omitempty,min=3,max=255"`
	Role     domain.AdminRole `json:"role" validate:"omitempty"`
}

type AdminResponse struct {
	ID          uuid.UUID          `json:"id"`
	Username    string             `json:"username"`
	Email       string             `json:"email"`
	FullName    string             `json:"full_name"`
	Role        domain.AdminRole   `json:"role"`
	Status      domain.AdminStatus `json:"status"`
	LastLoginAt *time.Time         `json:"last_login_at,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
}

// Login authenticates admin
func (uc *adminUsecase) Login(ctx context.Context, req AdminLoginRequest) (*AdminLoginResponse, error) {
	// Validate input
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get admin by username
	admin, err := uc.adminRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		// Create audit log for failed login
		_ = uc.createAuditLog(ctx, uuid.Nil, domain.AuditActionLogin, "", nil,
			fmt.Sprintf("Failed login attempt for username: %s", req.Username), "", "")
		return nil, domain.ErrInvalidCredentials
	}

	// Verify password
	if !crypto.VerifyPassword(req.Password, admin.PasswordHash) {
		_ = uc.createAuditLog(ctx, admin.ID, domain.AuditActionLogin, "", nil,
			"Failed login: invalid password", "", "")
		return nil, domain.ErrInvalidCredentials
	}

	// Check if admin is active
	if !admin.IsActive() {
		_ = uc.createAuditLog(ctx, admin.ID, domain.AuditActionLogin, "", nil,
			"Failed login: admin not active", "", "")
		return nil, domain.ErrAdminNotActive
	}

	// Update last login
	if err := uc.adminRepo.UpdateLastLogin(ctx, admin.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("failed to update last login: %v\n", err)
	}

	// Generate JWT token
	token, err := uc.tokenManager.GenerateAdminToken(admin.ID, admin.Username, string(admin.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create audit log for successful login
	_ = uc.createAuditLog(ctx, admin.ID, domain.AuditActionLogin, "", nil,
		"Successful login", "", "")

	return &AdminLoginResponse{
		AdminID:  admin.ID,
		Username: admin.Username,
		FullName: admin.FullName,
		Role:     admin.Role,
		Token:    token,
	}, nil
}

// CreateAdmin creates new admin (only by super admin)
func (uc *adminUsecase) CreateAdmin(ctx context.Context, creatorID uuid.UUID, req CreateAdminRequest) (*AdminResponse, error) {
	// Validate input
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if username already exists
	existingAdmin, _ := uc.adminRepo.GetByUsername(ctx, req.Username)
	if existingAdmin != nil {
		return nil, domain.ErrAdminAlreadyExist
	}

	// Check if email already exists
	existingAdmin, _ = uc.adminRepo.GetByEmail(ctx, req.Email)
	if existingAdmin != nil {
		return nil, domain.ErrAdminAlreadyExist
	}

	// Hash password
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin
	now := time.Now()
	admin := &domain.Admin{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		Role:         req.Role,
		Status:       domain.AdminStatusActive,
		CreatedBy:    &creatorID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.adminRepo.Create(ctx, admin); err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	// Create audit log
	_ = uc.createAuditLog(ctx, creatorID, domain.AuditActionViewUser, "admin", &admin.ID,
		fmt.Sprintf("Created new admin: %s (role: %s)", admin.Username, admin.Role), "", "")

	return &AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		FullName:  admin.FullName,
		Role:      admin.Role,
		Status:    admin.Status,
		CreatedAt: admin.CreatedAt,
	}, nil
}

// GetAdminByID gets admin by ID
func (uc *adminUsecase) GetAdminByID(ctx context.Context, id uuid.UUID) (*AdminResponse, error) {
	admin, err := uc.adminRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &AdminResponse{
		ID:          admin.ID,
		Username:    admin.Username,
		Email:       admin.Email,
		FullName:    admin.FullName,
		Role:        admin.Role,
		Status:      admin.Status,
		LastLoginAt: admin.LastLoginAt,
		CreatedAt:   admin.CreatedAt,
	}, nil
}

// ListAdmins lists all admins
func (uc *adminUsecase) ListAdmins(ctx context.Context, limit, offset int) ([]*AdminResponse, error) {
	admins, err := uc.adminRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*AdminResponse, 0, len(admins))
	for _, admin := range admins {
		responses = append(responses, &AdminResponse{
			ID:          admin.ID,
			Username:    admin.Username,
			Email:       admin.Email,
			FullName:    admin.FullName,
			Role:        admin.Role,
			Status:      admin.Status,
			LastLoginAt: admin.LastLoginAt,
			CreatedAt:   admin.CreatedAt,
		})
	}

	return responses, nil
}

// UpdateAdmin updates admin information
func (uc *adminUsecase) UpdateAdmin(ctx context.Context, id uuid.UUID, req UpdateAdminRequest) error {
	// Validate input
	if err := validator.ValidateStruct(req); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Get existing admin
	admin, err := uc.adminRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields if provided
	if req.Username != "" {
		admin.Username = req.Username
	}
	if req.Email != "" {
		admin.Email = req.Email
	}
	if req.FullName != "" {
		admin.FullName = req.FullName
	}
	if req.Role != "" {
		admin.Role = req.Role
	}

	admin.UpdatedAt = time.Now()

	// Update admin
	if err := uc.adminRepo.Update(ctx, admin); err != nil {
		return fmt.Errorf("failed to update admin: %w", err)
	}

	return nil
}

// UpdateAdminStatus updates admin status (suspend/activate)
func (uc *adminUsecase) UpdateAdminStatus(ctx context.Context, actorID, targetID uuid.UUID, status domain.AdminStatus) error {
	// Prevent self-action
	if actorID == targetID {
		return domain.ErrSelfAction
	}

	// Update status
	if err := uc.adminRepo.UpdateStatus(ctx, targetID, status); err != nil {
		return err
	}

	// Create audit log
	_ = uc.createAuditLog(ctx, actorID, domain.AuditActionViewUser, "admin", &targetID,
		fmt.Sprintf("Updated admin status to: %s", status), "", "")

	return nil
}

// GetAuditLogs gets audit logs for admin
func (uc *adminUsecase) GetAuditLogs(ctx context.Context, adminID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error) {
	return uc.auditLogRepo.GetByAdminID(ctx, adminID, limit, offset)
}

// Helper: Create audit log
func (uc *adminUsecase) createAuditLog(
	ctx context.Context,
	adminID uuid.UUID,
	action domain.AuditAction,
	resourceType string,
	resourceID *uuid.UUID,
	description string,
	ipAddress string,
	userAgent string,
) error {
	log := &domain.AuditLog{
		ID:           uuid.New(),
		AdminID:      adminID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Description:  description,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}

	return uc.auditLogRepo.Create(ctx, log)
}

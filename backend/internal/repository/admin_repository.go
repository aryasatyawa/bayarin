package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *domain.Admin) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Admin, error)
	GetByUsername(ctx context.Context, username string) (*domain.Admin, error)
	GetByEmail(ctx context.Context, email string) (*domain.Admin, error)
	Update(ctx context.Context, admin *domain.Admin) error
	UpdateLastLogin(ctx context.Context, adminID uuid.UUID) error
	UpdateStatus(ctx context.Context, adminID uuid.UUID, status domain.AdminStatus) error
	List(ctx context.Context, limit, offset int) ([]*domain.Admin, error)
}

type adminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) Create(ctx context.Context, admin *domain.Admin) error {
	query := `
		INSERT INTO admins (id, username, email, password_hash, full_name, role, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(
		ctx, query,
		admin.ID, admin.Username, admin.Email, admin.PasswordHash,
		admin.FullName, admin.Role, admin.Status, admin.CreatedBy,
		admin.CreatedAt, admin.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	return nil
}

func (r *adminRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Admin, error) {
	var admin domain.Admin
	query := `
		SELECT id, username, email, password_hash, full_name, role, status,
		       last_login_at, created_by, created_at, updated_at
		FROM admins
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &admin, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAdminNotFound
		}
		return nil, fmt.Errorf("failed to get admin by id: %w", err)
	}

	return &admin, nil
}

func (r *adminRepository) GetByUsername(ctx context.Context, username string) (*domain.Admin, error) {
	var admin domain.Admin
	query := `
		SELECT id, username, email, password_hash, full_name, role, status,
		       last_login_at, created_by, created_at, updated_at
		FROM admins
		WHERE username = $1
	`

	err := r.db.GetContext(ctx, &admin, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAdminNotFound
		}
		return nil, fmt.Errorf("failed to get admin by username: %w", err)
	}

	return &admin, nil
}

func (r *adminRepository) GetByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	var admin domain.Admin
	query := `
		SELECT id, username, email, password_hash, full_name, role, status,
		       last_login_at, created_by, created_at, updated_at
		FROM admins
		WHERE email = $1
	`

	err := r.db.GetContext(ctx, &admin, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAdminNotFound
		}
		return nil, fmt.Errorf("failed to get admin by email: %w", err)
	}

	return &admin, nil
}

func (r *adminRepository) Update(ctx context.Context, admin *domain.Admin) error {
	query := `
		UPDATE admins
		SET username = $1, email = $2, full_name = $3, role = $4, status = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(
		ctx, query,
		admin.Username, admin.Email, admin.FullName, admin.Role, admin.Status,
		time.Now(), admin.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update admin: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrAdminNotFound
	}

	return nil
}

func (r *adminRepository) UpdateLastLogin(ctx context.Context, adminID uuid.UUID) error {
	query := `
		UPDATE admins
		SET last_login_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, adminID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

func (r *adminRepository) UpdateStatus(ctx context.Context, adminID uuid.UUID, status domain.AdminStatus) error {
	query := `
		UPDATE admins
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, adminID)
	if err != nil {
		return fmt.Errorf("failed to update admin status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrAdminNotFound
	}

	return nil
}

func (r *adminRepository) List(ctx context.Context, limit, offset int) ([]*domain.Admin, error) {
	var admins []*domain.Admin
	query := `
		SELECT id, username, email, password_hash, full_name, role, status,
		       last_login_at, created_by, created_at, updated_at
		FROM admins
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	err := r.db.SelectContext(ctx, &admins, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list admins: %w", err)
	}

	return admins, nil
}

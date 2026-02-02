package repository

import (
	"context"
	"fmt"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	GetByAdminID(ctx context.Context, adminID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error)
	GetByAction(ctx context.Context, action domain.AuditAction, limit, offset int) ([]*domain.AuditLog, error)
	GetByResourceID(ctx context.Context, resourceType string, resourceID uuid.UUID) ([]*domain.AuditLog, error)
}

type auditLogRepository struct {
	db *sqlx.DB
}

func NewAuditLogRepository(db *sqlx.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			id, admin_id, action, resource_type, resource_id,
			description, ip_address, user_agent, before_value, after_value, metadata, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(
		ctx, query,
		log.ID, log.AdminID, log.Action, log.ResourceType, log.ResourceID,
		log.Description, log.IPAddress, log.UserAgent, log.BeforeValue, log.AfterValue,
		log.Metadata, log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

func (r *auditLogRepository) GetByAdminID(ctx context.Context, adminID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	query := `
		SELECT id, admin_id, action, resource_type, resource_id,
		       description, ip_address, user_agent, before_value, after_value, metadata, created_at
		FROM audit_logs
		WHERE admin_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &logs, query, adminID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, nil
}

func (r *auditLogRepository) GetByAction(ctx context.Context, action domain.AuditAction, limit, offset int) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	query := `
		SELECT id, admin_id, action, resource_type, resource_id,
		       description, ip_address, user_agent, before_value, after_value, metadata, created_at
		FROM audit_logs
		WHERE action = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &logs, query, action, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by action: %w", err)
	}

	return logs, nil
}

func (r *auditLogRepository) GetByResourceID(ctx context.Context, resourceType string, resourceID uuid.UUID) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	query := `
		SELECT id, admin_id, action, resource_type, resource_id,
		       description, ip_address, user_agent, before_value, after_value, metadata, created_at
		FROM audit_logs
		WHERE resource_type = $1 AND resource_id = $2
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &logs, query, resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by resource: %w", err)
	}

	return logs, nil
}

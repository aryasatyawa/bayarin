package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditAction string

const (
	AuditActionLogin              AuditAction = "login"
	AuditActionLogout             AuditAction = "logout"
	AuditActionViewUser           AuditAction = "view_user"
	AuditActionFreezeWallet       AuditAction = "freeze_wallet"
	AuditActionUnfreezeWallet     AuditAction = "unfreeze_wallet"
	AuditActionViewTransaction    AuditAction = "view_transaction"
	AuditActionRefundTransaction  AuditAction = "refund_transaction"
	AuditActionReverseTransaction AuditAction = "reverse_transaction"
	AuditActionViewLedger         AuditAction = "view_ledger"
	AuditActionDailySettlement    AuditAction = "daily_settlement"
	AuditActionCreateQR           AuditAction = "create_qr"
	AuditActionUpdateQR           AuditAction = "update_qr"
	AuditActionDeleteQR           AuditAction = "delete_qr"
)

type AuditLog struct {
	ID           uuid.UUID   `db:"id" json:"id"`
	AdminID      uuid.UUID   `db:"admin_id" json:"admin_id"`
	Action       AuditAction `db:"action" json:"action"`
	ResourceType string      `db:"resource_type" json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID  `db:"resource_id" json:"resource_id,omitempty"`
	Description  string      `db:"description" json:"description"`
	IPAddress    string      `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent    string      `db:"user_agent" json:"user_agent,omitempty"`
	BeforeValue  []byte      `db:"before_value" json:"before_value,omitempty"`
	AfterValue   []byte      `db:"after_value" json:"after_value,omitempty"`
	Metadata     []byte      `db:"metadata" json:"metadata,omitempty"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
}

// NewAuditLog creates new audit log entry
func NewAuditLog(adminID uuid.UUID, action AuditAction, description string) *AuditLog {
	return &AuditLog{
		ID:          uuid.New(),
		AdminID:     adminID,
		Action:      action,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

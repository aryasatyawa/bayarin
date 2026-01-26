package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID         `db:"id" json:"id"`
	IdempotencyKey  string            `db:"idempotency_key" json:"idempotency_key"`
	UserID          uuid.UUID         `db:"user_id" json:"user_id"`
	TransactionType TransactionType   `db:"transaction_type" json:"transaction_type"`
	Amount          int64             `db:"amount" json:"amount"` // WAJIB INTEGER
	Currency        string            `db:"currency" json:"currency"`
	Status          TransactionStatus `db:"status" json:"status"`
	FromWalletID    *uuid.UUID        `db:"from_wallet_id" json:"from_wallet_id,omitempty"`
	ToWalletID      *uuid.UUID        `db:"to_wallet_id" json:"to_wallet_id,omitempty"`
	ReferenceID     *string           `db:"reference_id" json:"reference_id,omitempty"`
	Description     string            `db:"description" json:"description"`
	Metadata        []byte            `db:"metadata" json:"metadata,omitempty"` // JSONB
	CreatedAt       time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at"`
	CompletedAt     *time.Time        `db:"completed_at" json:"completed_at,omitempty"`
}

type TransactionType string

const (
	TransactionTypeTopup      TransactionType = "topup"
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

type TransactionStatus string

const (
	TransactionStatusPending  TransactionStatus = "pending"
	TransactionStatusSuccess  TransactionStatus = "success"
	TransactionStatusFailed   TransactionStatus = "failed"
	TransactionStatusReversed TransactionStatus = "reversed"
)

// IsCompleted checks if transaction is in final state
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusSuccess ||
		t.Status == TransactionStatusFailed ||
		t.Status == TransactionStatusReversed
}

// IsPending checks if transaction is pending
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// MarkSuccess marks transaction as success
func (t *Transaction) MarkSuccess() {
	now := time.Now()
	t.Status = TransactionStatusSuccess
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// MarkFailed marks transaction as failed
func (t *Transaction) MarkFailed() {
	now := time.Now()
	t.Status = TransactionStatusFailed
	t.CompletedAt = &now
	t.UpdatedAt = now
}

package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, transaction *domain.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Transaction, error)
	UpdateStatus(ctx context.Context, tx *sqlx.Tx, transactionID uuid.UUID, status domain.TransactionStatus) error
}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *sqlx.Tx, transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, idempotency_key, user_id, transaction_type, amount, currency,
			status, from_wallet_id, to_wallet_id, reference_id, description,
			metadata, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		transaction.ID,
		transaction.IdempotencyKey,
		transaction.UserID,
		transaction.TransactionType,
		transaction.Amount,
		transaction.Currency,
		transaction.Status,
		transaction.FromWalletID,
		transaction.ToWalletID,
		transaction.ReferenceID,
		transaction.Description,
		transaction.Metadata,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var transaction domain.Transaction
	query := `
		SELECT id, idempotency_key, user_id, transaction_type, amount, currency,
			   status, from_wallet_id, to_wallet_id, reference_id, description,
			   metadata, created_at, updated_at, completed_at
		FROM transactions
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &transaction, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}

	return &transaction, nil
}

func (r *transactionRepository) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	query := `
		SELECT id, idempotency_key, user_id, transaction_type, amount, currency,
			   status, from_wallet_id, to_wallet_id, reference_id, description,
			   metadata, created_at, updated_at, completed_at
		FROM transactions
		WHERE idempotency_key = $1
	`

	err := r.db.GetContext(ctx, &transaction, query, idempotencyKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction by idempotency key: %w", err)
	}

	return &transaction, nil
}

func (r *transactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	query := `
		SELECT id, idempotency_key, user_id, transaction_type, amount, currency,
			   status, from_wallet_id, to_wallet_id, reference_id, description,
			   metadata, created_at, updated_at, completed_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &transactions, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by user id: %w", err)
	}

	return transactions, nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, tx *sqlx.Tx, transactionID uuid.UUID, status domain.TransactionStatus) error {
	query := `
		UPDATE transactions
		SET status = $1, 
		    updated_at = NOW(),
		    completed_at = CASE WHEN $1 IN ('success', 'failed', 'reversed') THEN NOW() ELSE completed_at END
		WHERE id = $2
	`

	result, err := tx.ExecContext(ctx, query, status, transactionID)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrTransactionNotFound
	}

	return nil
}

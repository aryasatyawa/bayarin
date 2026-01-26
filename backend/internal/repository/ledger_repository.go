package repository

import (
	"context"
	"fmt"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LedgerRepository interface {
	CreateEntry(ctx context.Context, tx *sqlx.Tx, entry *domain.LedgerEntry) error
	CreateEntries(ctx context.Context, tx *sqlx.Tx, entries []*domain.LedgerEntry) error
	GetByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*domain.LedgerEntry, error)
	GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*domain.LedgerEntry, error)
}

type ledgerRepository struct {
	db *sqlx.DB
}

func NewLedgerRepository(db *sqlx.DB) LedgerRepository {
	return &ledgerRepository{db: db}
}

func (r *ledgerRepository) CreateEntry(ctx context.Context, tx *sqlx.Tx, entry *domain.LedgerEntry) error {
	query := `
		INSERT INTO ledger_entries (
			id, transaction_id, wallet_id, entry_type, amount,
			balance_before, balance_after, description, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		entry.ID,
		entry.TransactionID,
		entry.WalletID,
		entry.EntryType,
		entry.Amount,
		entry.BalanceBefore,
		entry.BalanceAfter,
		entry.Description,
		entry.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create ledger entry: %w", err)
	}

	return nil
}

func (r *ledgerRepository) CreateEntries(ctx context.Context, tx *sqlx.Tx, entries []*domain.LedgerEntry) error {
	query := `
		INSERT INTO ledger_entries (
			id, transaction_id, wallet_id, entry_type, amount,
			balance_before, balance_after, description, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	for _, entry := range entries {
		_, err := tx.ExecContext(
			ctx,
			query,
			entry.ID,
			entry.TransactionID,
			entry.WalletID,
			entry.EntryType,
			entry.Amount,
			entry.BalanceBefore,
			entry.BalanceAfter,
			entry.Description,
			entry.CreatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to create ledger entries: %w", err)
		}
	}

	return nil
}

func (r *ledgerRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*domain.LedgerEntry, error) {
	var entries []*domain.LedgerEntry
	query := `
		SELECT id, transaction_id, wallet_id, entry_type, amount,
			   balance_before, balance_after, description, created_at
		FROM ledger_entries
		WHERE transaction_id = $1
		ORDER BY created_at ASC
	`

	err := r.db.SelectContext(ctx, &entries, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger entries by transaction id: %w", err)
	}

	return entries, nil
}

func (r *ledgerRepository) GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*domain.LedgerEntry, error) {
	var entries []*domain.LedgerEntry
	query := `
		SELECT id, transaction_id, wallet_id, entry_type, amount,
			   balance_before, balance_after, description, created_at
		FROM ledger_entries
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &entries, query, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger entries by wallet id: %w", err)
	}

	return entries, nil
}

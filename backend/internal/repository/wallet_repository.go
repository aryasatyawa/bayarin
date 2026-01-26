package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, wallet *domain.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	GetByIDWithTx(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*domain.Wallet, error)
	GetByUserIDAndType(ctx context.Context, userID uuid.UUID, walletType domain.WalletType) (*domain.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Wallet, error)
	UpdateBalance(ctx context.Context, tx *sqlx.Tx, walletID uuid.UUID, newBalance int64) error
	LockForUpdate(ctx context.Context, tx *sqlx.Tx, walletID uuid.UUID) (*domain.Wallet, error)
}

type walletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, wallet_type, balance, currency, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		wallet.ID,
		wallet.UserID,
		wallet.WalletType,
		wallet.Balance,
		wallet.Currency,
		wallet.Status,
		wallet.CreatedAt,
		wallet.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (r *walletRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, wallet *domain.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, wallet_type, balance, currency, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		wallet.ID,
		wallet.UserID,
		wallet.WalletType,
		wallet.Balance,
		wallet.Currency,
		wallet.Status,
		wallet.CreatedAt,
		wallet.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create wallet with tx: %w", err)
	}

	return nil
}

func (r *walletRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := `
		SELECT id, user_id, wallet_type, balance, currency, status, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &wallet, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet by id: %w", err)
	}

	return &wallet, nil
}

func (r *walletRepository) GetByIDWithTx(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := `
		SELECT id, user_id, wallet_type, balance, currency, status, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`

	err := tx.GetContext(ctx, &wallet, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet by id with tx: %w", err)
	}

	return &wallet, nil
}

func (r *walletRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, walletType domain.WalletType) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := `
		SELECT id, user_id, wallet_type, balance, currency, status, created_at, updated_at
		FROM wallets
		WHERE user_id = $1 AND wallet_type = $2
	`

	err := r.db.GetContext(ctx, &wallet, query, userID, walletType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet by user and type: %w", err)
	}

	return &wallet, nil
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Wallet, error) {
	var wallets []*domain.Wallet
	query := `
		SELECT id, user_id, wallet_type, balance, currency, status, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	err := r.db.SelectContext(ctx, &wallets, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets by user id: %w", err)
	}

	return wallets, nil
}

// UpdateBalance updates wallet balance - MUST be called within transaction
func (r *walletRepository) UpdateBalance(ctx context.Context, tx *sqlx.Tx, walletID uuid.UUID, newBalance int64) error {
	query := `
		UPDATE wallets
		SET balance = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := tx.ExecContext(ctx, query, newBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrWalletNotFound
	}

	return nil
}

// LockForUpdate locks wallet row for update (SELECT ... FOR UPDATE)
// CRITICAL: Prevents race condition dalam concurrent transactions
func (r *walletRepository) LockForUpdate(ctx context.Context, tx *sqlx.Tx, walletID uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := `
		SELECT id, user_id, wallet_type, balance, currency, status, created_at, updated_at
		FROM wallets
		WHERE id = $1
		FOR UPDATE
	`

	err := tx.GetContext(ctx, &wallet, query, walletID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to lock wallet for update: %w", err)
	}

	return &wallet, nil
}

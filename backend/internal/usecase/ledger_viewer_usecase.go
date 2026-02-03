package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LedgerViewerUsecase interface {
	GetLedgerEntries(ctx context.Context, filter LedgerFilter) (*LedgerEntriesResponse, error)
	GetLedgerByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*LedgerEntryDetail, error)
	GetLedgerByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) (*WalletLedgerResponse, error)
	ValidateBalance(ctx context.Context, walletID uuid.UUID) (*BalanceValidation, error)
}

type ledgerViewerUsecase struct {
	db         *sqlx.DB
	ledgerRepo repository.LedgerRepository
	walletRepo repository.WalletRepository
	txRepo     repository.TransactionRepository
}

func NewLedgerViewerUsecase(
	db *sqlx.DB,
	ledgerRepo repository.LedgerRepository,
	walletRepo repository.WalletRepository,
	txRepo repository.TransactionRepository,
) LedgerViewerUsecase {
	return &ledgerViewerUsecase{
		db:         db,
		ledgerRepo: ledgerRepo,
		walletRepo: walletRepo,
		txRepo:     txRepo,
	}
}

// DTOs
type LedgerFilter struct {
	UserID        *uuid.UUID        `json:"user_id,omitempty"`
	WalletID      *uuid.UUID        `json:"wallet_id,omitempty"`
	TransactionID *uuid.UUID        `json:"transaction_id,omitempty"`
	EntryType     *domain.EntryType `json:"entry_type,omitempty"`
	StartDate     *time.Time        `json:"start_date,omitempty"`
	EndDate       *time.Time        `json:"end_date,omitempty"`
	Limit         int               `json:"limit"`
	Offset        int               `json:"offset"`
}

type LedgerEntriesResponse struct {
	Entries []*LedgerEntryDetail `json:"entries"`
	Total   int                  `json:"total"`
	Filter  LedgerFilter         `json:"filter"`
}

type LedgerEntryDetail struct {
	ID              uuid.UUID        `json:"id"`
	TransactionID   uuid.UUID        `json:"transaction_id"`
	WalletID        uuid.UUID        `json:"wallet_id"`
	UserID          uuid.UUID        `json:"user_id"` // From wallet
	EntryType       domain.EntryType `json:"entry_type"`
	Amount          int64            `json:"amount"`
	BalanceBefore   int64            `json:"balance_before"`
	BalanceAfter    int64            `json:"balance_after"`
	Description     string           `json:"description"`
	TransactionType string           `json:"transaction_type"` // From transaction
	CreatedAt       time.Time        `json:"created_at"`
}

type WalletLedgerResponse struct {
	WalletID       uuid.UUID            `json:"wallet_id"`
	CurrentBalance int64                `json:"current_balance"`
	Entries        []*LedgerEntryDetail `json:"entries"`
	Total          int                  `json:"total"`
}

type BalanceValidation struct {
	WalletID          uuid.UUID `json:"wallet_id"`
	CurrentBalance    int64     `json:"current_balance"`
	CalculatedBalance int64     `json:"calculated_balance"`
	IsValid           bool      `json:"is_valid"`
	Difference        int64     `json:"difference"`
	Message           string    `json:"message"`
}

// GetLedgerEntries returns filtered ledger entries
func (uc *ledgerViewerUsecase) GetLedgerEntries(ctx context.Context, filter LedgerFilter) (*LedgerEntriesResponse, error) {
	// Build query
	query := `
		SELECT 
			le.id, le.transaction_id, le.wallet_id, le.entry_type,
			le.amount, le.balance_before, le.balance_after, le.description,
			le.created_at, w.user_id, t.transaction_type
		FROM ledger_entries le
		INNER JOIN wallets w ON le.wallet_id = w.id
		INNER JOIN transactions t ON le.transaction_id = t.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// Apply filters
	if filter.UserID != nil {
		query += fmt.Sprintf(" AND w.user_id = $%d", argCount)
		args = append(args, *filter.UserID)
		argCount++
	}

	if filter.WalletID != nil {
		query += fmt.Sprintf(" AND le.wallet_id = $%d", argCount)
		args = append(args, *filter.WalletID)
		argCount++
	}

	if filter.TransactionID != nil {
		query += fmt.Sprintf(" AND le.transaction_id = $%d", argCount)
		args = append(args, *filter.TransactionID)
		argCount++
	}

	if filter.EntryType != nil {
		query += fmt.Sprintf(" AND le.entry_type = $%d", argCount)
		args = append(args, *filter.EntryType)
		argCount++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND le.created_at >= $%d", argCount)
		args = append(args, *filter.StartDate)
		argCount++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND le.created_at <= $%d", argCount)
		args = append(args, *filter.EndDate)
		argCount++
	}

	// Order and pagination
	query += " ORDER BY le.created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}

	// Execute query
	var entries []*LedgerEntryDetail
	if err := uc.db.SelectContext(ctx, &entries, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get ledger entries: %w", err)
	}

	// Get total count (without pagination)
	countQuery := `
		SELECT COUNT(*)
		FROM ledger_entries le
		INNER JOIN wallets w ON le.wallet_id = w.id
		INNER JOIN transactions t ON le.transaction_id = t.id
		WHERE 1=1
	`

	countArgs := []interface{}{}
	countArgCount := 1

	if filter.UserID != nil {
		countQuery += fmt.Sprintf(" AND w.user_id = $%d", countArgCount)
		countArgs = append(countArgs, *filter.UserID)
		countArgCount++
	}

	if filter.WalletID != nil {
		countQuery += fmt.Sprintf(" AND le.wallet_id = $%d", countArgCount)
		countArgs = append(countArgs, *filter.WalletID)
		countArgCount++
	}

	if filter.TransactionID != nil {
		countQuery += fmt.Sprintf(" AND le.transaction_id = $%d", countArgCount)
		countArgs = append(countArgs, *filter.TransactionID)
		countArgCount++
	}

	var total int
	if err := uc.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &LedgerEntriesResponse{
		Entries: entries,
		Total:   total,
		Filter:  filter,
	}, nil
}

// GetLedgerByTransactionID returns all ledger entries for a transaction
func (uc *ledgerViewerUsecase) GetLedgerByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*LedgerEntryDetail, error) {
	query := `
		SELECT 
			le.id, le.transaction_id, le.wallet_id, le.entry_type,
			le.amount, le.balance_before, le.balance_after, le.description,
			le.created_at, w.user_id, t.transaction_type
		FROM ledger_entries le
		INNER JOIN wallets w ON le.wallet_id = w.id
		INNER JOIN transactions t ON le.transaction_id = t.id
		WHERE le.transaction_id = $1
		ORDER BY le.created_at ASC
	`

	var entries []*LedgerEntryDetail
	if err := uc.db.SelectContext(ctx, &entries, query, transactionID); err != nil {
		return nil, fmt.Errorf("failed to get ledger by transaction: %w", err)
	}

	return entries, nil
}

// GetLedgerByWalletID returns ledger entries for a wallet
func (uc *ledgerViewerUsecase) GetLedgerByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) (*WalletLedgerResponse, error) {
	// Get wallet current balance
	wallet, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	// Get ledger entries
	query := `
		SELECT 
			le.id, le.transaction_id, le.wallet_id, le.entry_type,
			le.amount, le.balance_before, le.balance_after, le.description,
			le.created_at, w.user_id, t.transaction_type
		FROM ledger_entries le
		INNER JOIN wallets w ON le.wallet_id = w.id
		INNER JOIN transactions t ON le.transaction_id = t.id
		WHERE le.wallet_id = $1
		ORDER BY le.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var entries []*LedgerEntryDetail
	if err := uc.db.SelectContext(ctx, &entries, query, walletID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get wallet ledger: %w", err)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM ledger_entries WHERE wallet_id = $1`
	if err := uc.db.GetContext(ctx, &total, countQuery, walletID); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &WalletLedgerResponse{
		WalletID:       walletID,
		CurrentBalance: wallet.Balance,
		Entries:        entries,
		Total:          total,
	}, nil
}

// ValidateBalance validates wallet balance against ledger
func (uc *ledgerViewerUsecase) ValidateBalance(ctx context.Context, walletID uuid.UUID) (*BalanceValidation, error) {
	// Get current wallet balance
	wallet, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	// Calculate balance from ledger
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN entry_type = 'credit' THEN amount ELSE 0 END), 0) as total_credit,
			COALESCE(SUM(CASE WHEN entry_type = 'debit' THEN amount ELSE 0 END), 0) as total_debit
		FROM ledger_entries
		WHERE wallet_id = $1
	`

	var result struct {
		TotalCredit int64 `db:"total_credit"`
		TotalDebit  int64 `db:"total_debit"`
	}

	if err := uc.db.GetContext(ctx, &result, query, walletID); err != nil {
		return nil, fmt.Errorf("failed to calculate balance: %w", err)
	}

	calculatedBalance := result.TotalCredit - result.TotalDebit
	isValid := wallet.Balance == calculatedBalance
	difference := wallet.Balance - calculatedBalance

	message := "Balance is valid"
	if !isValid {
		message = fmt.Sprintf("Balance mismatch! Difference: %d", difference)
	}

	return &BalanceValidation{
		WalletID:          walletID,
		CurrentBalance:    wallet.Balance,
		CalculatedBalance: calculatedBalance,
		IsValid:           isValid,
		Difference:        difference,
		Message:           message,
	}, nil
}

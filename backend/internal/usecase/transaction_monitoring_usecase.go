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

type TransactionMonitoringUsecase interface {
	GetAllTransactions(ctx context.Context, filter TransactionFilter) (*TransactionListResponse, error)
	GetTransactionDetail(ctx context.Context, transactionID uuid.UUID) (*TransactionDetailResponse, error)
	GetPendingTransactions(ctx context.Context, limit, offset int) ([]*TransactionDetailResponse, error)
	GetFailedTransactions(ctx context.Context, days int, limit, offset int) ([]*TransactionDetailResponse, error)
	GetTransactionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*TransactionDetailResponse, error)
}

type transactionMonitoringUsecase struct {
	db         *sqlx.DB
	txRepo     repository.TransactionRepository
	ledgerRepo repository.LedgerRepository
	userRepo   repository.UserRepository
}

func NewTransactionMonitoringUsecase(
	db *sqlx.DB,
	txRepo repository.TransactionRepository,
	ledgerRepo repository.LedgerRepository,
	userRepo repository.UserRepository,
) TransactionMonitoringUsecase {
	return &transactionMonitoringUsecase{
		db:         db,
		txRepo:     txRepo,
		ledgerRepo: ledgerRepo,
		userRepo:   userRepo,
	}
}

// DTOs
type TransactionFilter struct {
	UserID          *uuid.UUID                `json:"user_id,omitempty"`
	TransactionType *domain.TransactionType   `json:"transaction_type,omitempty"`
	Status          *domain.TransactionStatus `json:"status,omitempty"`
	StartDate       *time.Time                `json:"start_date,omitempty"`
	EndDate         *time.Time                `json:"end_date,omitempty"`
	MinAmount       *int64                    `json:"min_amount,omitempty"`
	MaxAmount       *int64                    `json:"max_amount,omitempty"`
	Limit           int                       `json:"limit"`
	Offset          int                       `json:"offset"`
}

type TransactionListResponse struct {
	Transactions []*TransactionDetailResponse `json:"transactions"`
	Total        int                          `json:"total"`
	Filter       TransactionFilter            `json:"filter"`
}

type TransactionDetailResponse struct {
	ID              uuid.UUID                `json:"id"`
	IdempotencyKey  string                   `json:"idempotency_key"`
	UserID          uuid.UUID                `json:"user_id"`
	UserEmail       string                   `json:"user_email"`
	TransactionType domain.TransactionType   `json:"transaction_type"`
	Amount          int64                    `json:"amount"`
	Currency        string                   `json:"currency"`
	Status          domain.TransactionStatus `json:"status"`
	FromWalletID    *uuid.UUID               `json:"from_wallet_id,omitempty"`
	ToWalletID      *uuid.UUID               `json:"to_wallet_id,omitempty"`
	ReferenceID     *string                  `json:"reference_id,omitempty"`
	Description     string                   `json:"description"`
	Metadata        []byte                   `json:"metadata,omitempty"`
	LedgerEntries   []*LedgerEntryDetail     `json:"ledger_entries,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	CompletedAt     *time.Time               `json:"completed_at,omitempty"`
}

// GetAllTransactions returns filtered transactions
func (uc *transactionMonitoringUsecase) GetAllTransactions(ctx context.Context, filter TransactionFilter) (*TransactionListResponse, error) {
	// Build query
	query := `
		SELECT 
			t.id, t.idempotency_key, t.user_id, t.transaction_type,
			t.amount, t.currency, t.status, t.from_wallet_id, t.to_wallet_id,
			t.reference_id, t.description, t.metadata,
			t.created_at, t.updated_at, t.completed_at,
			u.email as user_email
		FROM transactions t
		INNER JOIN users u ON t.user_id = u.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	// Apply filters
	if filter.UserID != nil {
		query += fmt.Sprintf(" AND t.user_id = $%d", argCount)
		args = append(args, *filter.UserID)
		argCount++
	}

	if filter.TransactionType != nil {
		query += fmt.Sprintf(" AND t.transaction_type = $%d", argCount)
		args = append(args, *filter.TransactionType)
		argCount++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND t.status = $%d", argCount)
		args = append(args, *filter.Status)
		argCount++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND t.created_at >= $%d", argCount)
		args = append(args, *filter.StartDate)
		argCount++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND t.created_at <= $%d", argCount)
		args = append(args, *filter.EndDate)
		argCount++
	}

	if filter.MinAmount != nil {
		query += fmt.Sprintf(" AND t.amount >= $%d", argCount)
		args = append(args, *filter.MinAmount)
		argCount++
	}

	if filter.MaxAmount != nil {
		query += fmt.Sprintf(" AND t.amount <= $%d", argCount)
		args = append(args, *filter.MaxAmount)
		argCount++
	}

	// Order and pagination
	query += " ORDER BY t.created_at DESC"

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
	var transactions []*TransactionDetailResponse
	if err := uc.db.SelectContext(ctx, &transactions, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Get total count (simplified count query without joins for performance)
	countQuery := `SELECT COUNT(*) FROM transactions WHERE 1=1`
	countArgs := []interface{}{}
	countArgCount := 1

	if filter.UserID != nil {
		countQuery += fmt.Sprintf(" AND user_id = $%d", countArgCount)
		countArgs = append(countArgs, *filter.UserID)
		countArgCount++
	}

	if filter.TransactionType != nil {
		countQuery += fmt.Sprintf(" AND transaction_type = $%d", countArgCount)
		countArgs = append(countArgs, *filter.TransactionType)
		countArgCount++
	}

	if filter.Status != nil {
		countQuery += fmt.Sprintf(" AND status = $%d", countArgCount)
		countArgs = append(countArgs, *filter.Status)
	}

	var total int
	if err := uc.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &TransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Filter:       filter,
	}, nil
}

// GetTransactionDetail returns detailed transaction with ledger entries
func (uc *transactionMonitoringUsecase) GetTransactionDetail(ctx context.Context, transactionID uuid.UUID) (*TransactionDetailResponse, error) {
	// Get transaction
	query := `
		SELECT 
			t.id, t.idempotency_key, t.user_id, t.transaction_type,
			t.amount, t.currency, t.status, t.from_wallet_id, t.to_wallet_id,
			t.reference_id, t.description, t.metadata,
			t.created_at, t.updated_at, t.completed_at,
			u.email as user_email
		FROM transactions t
		INNER JOIN users u ON t.user_id = u.id
		WHERE t.id = $1
	`

	var tx TransactionDetailResponse
	if err := uc.db.GetContext(ctx, &tx, query, transactionID); err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Get ledger entries
	ledgerEntries, err := uc.ledgerRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger entries: %w", err)
	}

	// Convert to detail format
	tx.LedgerEntries = make([]*LedgerEntryDetail, 0, len(ledgerEntries))
	for _, entry := range ledgerEntries {
		tx.LedgerEntries = append(tx.LedgerEntries, &LedgerEntryDetail{
			ID:            entry.ID,
			TransactionID: entry.TransactionID,
			WalletID:      entry.WalletID,
			EntryType:     entry.EntryType,
			Amount:        entry.Amount,
			BalanceBefore: entry.BalanceBefore,
			BalanceAfter:  entry.BalanceAfter,
			Description:   entry.Description,
			CreatedAt:     entry.CreatedAt,
		})
	}

	return &tx, nil
}

// GetPendingTransactions returns all pending transactions
func (uc *transactionMonitoringUsecase) GetPendingTransactions(ctx context.Context, limit, offset int) ([]*TransactionDetailResponse, error) {
	query := `
		SELECT 
			t.id, t.idempotency_key, t.user_id, t.transaction_type,
			t.amount, t.currency, t.status, t.from_wallet_id, t.to_wallet_id,
			t.reference_id, t.description, t.metadata,
			t.created_at, t.updated_at, t.completed_at,
			u.email as user_email
		FROM transactions t
		INNER JOIN users u ON t.user_id = u.id
		WHERE t.status = 'pending'
		ORDER BY t.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var transactions []*TransactionDetailResponse
	if err := uc.db.SelectContext(ctx, &transactions, query, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get pending transactions: %w", err)
	}

	return transactions, nil
}

// GetFailedTransactions returns failed transactions within specified days
func (uc *transactionMonitoringUsecase) GetFailedTransactions(ctx context.Context, days int, limit, offset int) ([]*TransactionDetailResponse, error) {
	query := `
		SELECT 
			t.id, t.idempotency_key, t.user_id, t.transaction_type,
			t.amount, t.currency, t.status, t.from_wallet_id, t.to_wallet_id,
			t.reference_id, t.description, t.metadata,
			t.created_at, t.updated_at, t.completed_at,
			u.email as user_email
		FROM transactions t
		INNER JOIN users u ON t.user_id = u.id
		WHERE t.status = 'failed'
		AND t.created_at >= NOW() - INTERVAL '%d days'
		ORDER BY t.created_at DESC
		LIMIT $1 OFFSET $2
	`

	formattedQuery := fmt.Sprintf(query, days)

	var transactions []*TransactionDetailResponse
	if err := uc.db.SelectContext(ctx, &transactions, formattedQuery, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get failed transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionsByUser returns transactions for specific user
func (uc *transactionMonitoringUsecase) GetTransactionsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*TransactionDetailResponse, error) {
	query := `
		SELECT 
			t.id, t.idempotency_key, t.user_id, t.transaction_type,
			t.amount, t.currency, t.status, t.from_wallet_id, t.to_wallet_id,
			t.reference_id, t.description, t.metadata,
			t.created_at, t.updated_at, t.completed_at,
			u.email as user_email
		FROM transactions t
		INNER JOIN users u ON t.user_id = u.id
		WHERE t.user_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var transactions []*TransactionDetailResponse
	if err := uc.db.SelectContext(ctx, &transactions, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}

	return transactions, nil
}

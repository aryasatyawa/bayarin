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

type UserInspectorUsecase interface {
	GetUserDetails(ctx context.Context, userID uuid.UUID) (*UserInspectorDetail, error)
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*WalletInspectorDetail, error)
	FreezeWallet(ctx context.Context, adminID, walletID uuid.UUID, reason string) error
	UnfreezeWallet(ctx context.Context, adminID, walletID uuid.UUID, reason string) error
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]*UserSearchResult, error)
}

type userInspectorUsecase struct {
	db           *sqlx.DB
	userRepo     repository.UserRepository
	walletRepo   repository.WalletRepository
	txRepo       repository.TransactionRepository
	auditLogRepo repository.AuditLogRepository
}

func NewUserInspectorUsecase(
	db *sqlx.DB,
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
	txRepo repository.TransactionRepository,
	auditLogRepo repository.AuditLogRepository,
) UserInspectorUsecase {
	return &userInspectorUsecase{
		db:           db,
		userRepo:     userRepo,
		walletRepo:   walletRepo,
		txRepo:       txRepo,
		auditLogRepo: auditLogRepo,
	}
}

// DTOs
type UserInspectorDetail struct {
	User                *domain.User             `json:"user"`
	Wallets             []*WalletInspectorDetail `json:"wallets"`
	TotalBalance        int64                    `json:"total_balance"`
	TotalTransactions   int64                    `json:"total_transactions"`
	SuccessTransactions int64                    `json:"success_transactions"`
	FailedTransactions  int64                    `json:"failed_transactions"`
	LastTransactionAt   *time.Time               `json:"last_transaction_at,omitempty"`
}

type WalletInspectorDetail struct {
	ID               uuid.UUID           `json:"id"`
	WalletType       domain.WalletType   `json:"wallet_type"`
	Balance          int64               `json:"balance"`
	Status           domain.WalletStatus `json:"status"`
	TransactionCount int64               `json:"transaction_count"`
	LastActivityAt   *time.Time          `json:"last_activity_at,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
}

type UserSearchResult struct {
	ID       uuid.UUID         `json:"id"`
	Email    string            `json:"email"`
	Phone    string            `json:"phone"`
	FullName string            `json:"full_name"`
	Status   domain.UserStatus `json:"status"`
}

// GetUserDetails returns comprehensive user details
func (uc *userInspectorUsecase) GetUserDetails(ctx context.Context, userID uuid.UUID) (*UserInspectorDetail, error) {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get wallets
	wallets, err := uc.GetUserWallets(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate total balance
	var totalBalance int64
	for _, wallet := range wallets {
		totalBalance += wallet.Balance
	}

	// Get transaction statistics
	var txStats struct {
		TotalCount   int64      `db:"total_count"`
		SuccessCount int64      `db:"success_count"`
		FailedCount  int64      `db:"failed_count"`
		LastTxAt     *time.Time `db:"last_tx_at"`
	}

	statsQuery := `
		SELECT 
			COUNT(*) as total_count,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as success_count,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count,
			MAX(created_at) as last_tx_at
		FROM transactions
		WHERE user_id = $1
	`

	if err := uc.db.GetContext(ctx, &txStats, statsQuery, userID); err != nil {
		return nil, fmt.Errorf("failed to get transaction stats: %w", err)
	}

	return &UserInspectorDetail{
		User:                user,
		Wallets:             wallets,
		TotalBalance:        totalBalance,
		TotalTransactions:   txStats.TotalCount,
		SuccessTransactions: txStats.SuccessCount,
		FailedTransactions:  txStats.FailedCount,
		LastTransactionAt:   txStats.LastTxAt,
	}, nil
}

// GetUserWallets returns user wallets with statistics
func (uc *userInspectorUsecase) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*WalletInspectorDetail, error) {
	query := `
		SELECT 
			w.id,
			w.wallet_type,
			w.balance,
			w.status,
			w.created_at,
			COUNT(DISTINCT le.id) as transaction_count,
			MAX(le.created_at) as last_activity_at
		FROM wallets w
		LEFT JOIN ledger_entries le ON w.id = le.wallet_id
		WHERE w.user_id = $1
		GROUP BY w.id, w.wallet_type, w.balance, w.status, w.created_at
		ORDER BY w.created_at ASC
	`

	var wallets []*WalletInspectorDetail
	if err := uc.db.SelectContext(ctx, &wallets, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}

	return wallets, nil
}

// FreezeWallet freezes a wallet (prevents transactions)
func (uc *userInspectorUsecase) FreezeWallet(ctx context.Context, adminID, walletID uuid.UUID, reason string) error {
	// Get wallet
	wallet, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	// Check if already frozen
	if wallet.Status == domain.WalletStatusFrozen {
		return fmt.Errorf("wallet is already frozen")
	}

	// Begin transaction
	tx, err := uc.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update wallet status
	query := `
		UPDATE wallets
		SET status = 'frozen', updated_at = NOW()
		WHERE id = $1
	`

	if _, err := tx.ExecContext(ctx, query, walletID); err != nil {
		return fmt.Errorf("failed to freeze wallet: %w", err)
	}

	// Create audit log
	auditLog := &domain.AuditLog{
		ID:           uuid.New(),
		AdminID:      adminID,
		Action:       domain.AuditActionFreezeWallet,
		ResourceType: "wallet",
		ResourceID:   &walletID,
		Description:  fmt.Sprintf("Froze wallet %s. Reason: %s", walletID.String()[:8], reason),
		CreatedAt:    time.Now(),
	}

	if err := uc.auditLogRepo.Create(ctx, auditLog); err != nil {
		// Log error but don't fail freeze
		fmt.Printf("failed to create audit log: %v\n", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UnfreezeWallet unfreezes a wallet
func (uc *userInspectorUsecase) UnfreezeWallet(ctx context.Context, adminID, walletID uuid.UUID, reason string) error {
	// Get wallet
	wallet, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	// Check if already active
	if wallet.Status == domain.WalletStatusActive {
		return fmt.Errorf("wallet is already active")
	}

	// Begin transaction
	tx, err := uc.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update wallet status
	query := `
		UPDATE wallets
		SET status = 'active', updated_at = NOW()
		WHERE id = $1
	`

	if _, err := tx.ExecContext(ctx, query, walletID); err != nil {
		return fmt.Errorf("failed to unfreeze wallet: %w", err)
	}

	// Create audit log
	auditLog := &domain.AuditLog{
		ID:           uuid.New(),
		AdminID:      adminID,
		Action:       domain.AuditActionUnfreezeWallet,
		ResourceType: "wallet",
		ResourceID:   &walletID,
		Description:  fmt.Sprintf("Unfroze wallet %s. Reason: %s", walletID.String()[:8], reason),
		CreatedAt:    time.Now(),
	}

	if err := uc.auditLogRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("failed to create audit log: %v\n", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SearchUsers searches users by email, phone, or name
func (uc *userInspectorUsecase) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*UserSearchResult, error) {
	searchQuery := `
		SELECT id, email, phone, full_name, status
		FROM users
		WHERE 
			email ILIKE $1 OR
			phone ILIKE $1 OR
			full_name ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	searchPattern := "%" + query + "%"

	var users []*UserSearchResult
	if err := uc.db.SelectContext(ctx, &users, searchQuery, searchPattern, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

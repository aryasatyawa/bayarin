package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefundUsecase interface {
	RefundTransaction(ctx context.Context, adminID uuid.UUID, req RefundRequest) (*RefundResponse, error)
	ReverseTransaction(ctx context.Context, adminID uuid.UUID, req ReverseRequest) (*RefundResponse, error)
	GetRefundHistory(ctx context.Context, originalTxID uuid.UUID) ([]*RefundHistoryItem, error)
}

type refundUsecase struct {
	db           *sqlx.DB
	txRepo       repository.TransactionRepository
	walletRepo   repository.WalletRepository
	ledgerRepo   repository.LedgerRepository
	auditLogRepo repository.AuditLogRepository
}

func NewRefundUsecase(
	db *sqlx.DB,
	txRepo repository.TransactionRepository,
	walletRepo repository.WalletRepository,
	ledgerRepo repository.LedgerRepository,
	auditLogRepo repository.AuditLogRepository,
) RefundUsecase {
	return &refundUsecase{
		db:           db,
		txRepo:       txRepo,
		walletRepo:   walletRepo,
		ledgerRepo:   ledgerRepo,
		auditLogRepo: auditLogRepo,
	}
}

// DTOs
type RefundRequest struct {
	OriginalTransactionID uuid.UUID `json:"original_transaction_id" validate:"required"`
	Reason                string    `json:"reason" validate:"required,min=10"`
	Amount                *int64    `json:"amount,omitempty"` // Partial refund jika ada
	IdempotencyKey        string    `json:"idempotency_key" validate:"required"`
}

type ReverseRequest struct {
	OriginalTransactionID uuid.UUID `json:"original_transaction_id" validate:"required"`
	Reason                string    `json:"reason" validate:"required,min=10"`
	IdempotencyKey        string    `json:"idempotency_key" validate:"required"`
}

type RefundResponse struct {
	RefundTransactionID   uuid.UUID                `json:"refund_transaction_id"`
	OriginalTransactionID uuid.UUID                `json:"original_transaction_id"`
	Amount                int64                    `json:"amount"`
	Status                domain.TransactionStatus `json:"status"`
	Reason                string                   `json:"reason"`
	CreatedAt             time.Time                `json:"created_at"`
}

type RefundHistoryItem struct {
	RefundTransactionID uuid.UUID                `json:"refund_transaction_id"`
	Amount              int64                    `json:"amount"`
	Reason              string                   `json:"reason"`
	Status              domain.TransactionStatus `json:"status"`
	CreatedAt           time.Time                `json:"created_at"`
}

// RefundTransaction creates refund transaction (full or partial)
// CRITICAL: Membuat transaksi BARU dengan ledger entries BARU
func (uc *refundUsecase) RefundTransaction(ctx context.Context, adminID uuid.UUID, req RefundRequest) (*RefundResponse, error) {
	// Check idempotency
	existingTx, _ := uc.txRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if existingTx != nil {
		return &RefundResponse{
			RefundTransactionID:   existingTx.ID,
			OriginalTransactionID: *existingTx.FromWalletID, // Temporary, should store in metadata
			Amount:                existingTx.Amount,
			Status:                existingTx.Status,
			CreatedAt:             existingTx.CreatedAt,
		}, nil
	}

	// Get original transaction
	originalTx, err := uc.txRepo.GetByID(ctx, req.OriginalTransactionID)
	if err != nil {
		return nil, fmt.Errorf("original transaction not found: %w", err)
	}

	// Validate original transaction status
	if originalTx.Status != domain.TransactionStatusSuccess {
		return nil, fmt.Errorf("can only refund successful transactions")
	}

	// Determine refund amount
	refundAmount := originalTx.Amount
	if req.Amount != nil && *req.Amount > 0 {
		// Partial refund
		if *req.Amount > originalTx.Amount {
			return nil, fmt.Errorf("refund amount cannot exceed original amount")
		}
		refundAmount = *req.Amount
	}

	// Begin transaction
	tx, err := uc.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare refund metadata
	metadata, _ := json.Marshal(map[string]interface{}{
		"original_transaction_id": originalTx.ID.String(),
		"refund_type":             "admin_refund",
		"reason":                  req.Reason,
		"admin_id":                adminID.String(),
		"is_partial":              req.Amount != nil,
	})

	now := time.Now()
	refundTx := &domain.Transaction{
		ID:              uuid.New(),
		IdempotencyKey:  req.IdempotencyKey,
		UserID:          originalTx.UserID,
		TransactionType: domain.TransactionTypeTopup, // Refund treated as topup
		Amount:          refundAmount,
		Currency:        originalTx.Currency,
		Status:          domain.TransactionStatusPending,
		ReferenceID:     stringPtr(fmt.Sprintf("REFUND-%s", originalTx.ID.String()[:8])),
		Description:     fmt.Sprintf("Refund for transaction %s: %s", originalTx.ID.String()[:8], req.Reason),
		Metadata:        metadata,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Determine which wallet to refund to
	var targetWalletID uuid.UUID
	if originalTx.TransactionType == domain.TransactionTypeTransfer && originalTx.FromWalletID != nil {
		// Refund to sender wallet
		targetWalletID = *originalTx.FromWalletID
	} else if originalTx.ToWalletID != nil {
		// Refund to receiver wallet (for topup reversals)
		targetWalletID = *originalTx.ToWalletID
	} else {
		return nil, fmt.Errorf("cannot determine target wallet for refund")
	}

	refundTx.ToWalletID = &targetWalletID

	// Lock wallet for update
	wallet, err := uc.walletRepo.LockForUpdate(ctx, tx, targetWalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to lock wallet: %w", err)
	}

	// Check wallet status
	if !wallet.IsActive() {
		return nil, domain.ErrWalletNotActive
	}

	// Create refund transaction
	if err := uc.txRepo.Create(ctx, tx, refundTx); err != nil {
		return nil, fmt.Errorf("failed to create refund transaction: %w", err)
	}

	// Create ledger entry (CREDIT - menambah saldo)
	ledgerEntry := domain.NewCreditEntry(
		refundTx.ID,
		targetWalletID,
		refundAmount,
		wallet.Balance,
		refundTx.Description,
	)

	if err := uc.ledgerRepo.CreateEntry(ctx, tx, ledgerEntry); err != nil {
		return nil, fmt.Errorf("failed to create ledger entry: %w", err)
	}

	// Update wallet balance
	newBalance := wallet.Balance + refundAmount
	if err := uc.walletRepo.UpdateBalance(ctx, tx, targetWalletID, newBalance); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Mark refund transaction as success
	refundTx.MarkSuccess()
	if err := uc.txRepo.UpdateStatus(ctx, tx, refundTx.ID, domain.TransactionStatusSuccess); err != nil {
		return nil, fmt.Errorf("failed to update refund status: %w", err)
	}

	// Create audit log
	auditLog := &domain.AuditLog{
		ID:           uuid.New(),
		AdminID:      adminID,
		Action:       domain.AuditActionRefundTransaction,
		ResourceType: "transaction",
		ResourceID:   &originalTx.ID,
		Description:  fmt.Sprintf("Refund transaction %s for amount %d. Reason: %s", originalTx.ID.String()[:8], refundAmount, req.Reason),
		CreatedAt:    now,
	}

	if err := uc.auditLogRepo.Create(ctx, auditLog); err != nil {
		// Log error but don't fail refund
		fmt.Printf("failed to create audit log: %v\n", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &RefundResponse{
		RefundTransactionID:   refundTx.ID,
		OriginalTransactionID: originalTx.ID,
		Amount:                refundAmount,
		Status:                refundTx.Status,
		Reason:                req.Reason,
		CreatedAt:             refundTx.CreatedAt,
	}, nil
}

// ReverseTransaction reverses a transaction (full reversal only)
func (uc *refundUsecase) ReverseTransaction(ctx context.Context, adminID uuid.UUID, req ReverseRequest) (*RefundResponse, error) {
	// Reverse is same as full refund
	refundReq := RefundRequest{
		OriginalTransactionID: req.OriginalTransactionID,
		Reason:                req.Reason,
		Amount:                nil, // Full amount
		IdempotencyKey:        req.IdempotencyKey,
	}

	response, err := uc.RefundTransaction(ctx, adminID, refundReq)
	if err != nil {
		return nil, err
	}

	// Create audit log for reversal
	auditLog := &domain.AuditLog{
		ID:           uuid.New(),
		AdminID:      adminID,
		Action:       domain.AuditActionReverseTransaction,
		ResourceType: "transaction",
		ResourceID:   &req.OriginalTransactionID,
		Description:  fmt.Sprintf("Reversed transaction %s. Reason: %s", req.OriginalTransactionID.String()[:8], req.Reason),
		CreatedAt:    time.Now(),
	}

	_ = uc.auditLogRepo.Create(ctx, auditLog)

	return response, nil
}

// GetRefundHistory returns refund history for original transaction
func (uc *refundUsecase) GetRefundHistory(ctx context.Context, originalTxID uuid.UUID) ([]*RefundHistoryItem, error) {
	query := `
		SELECT 
			id as refund_transaction_id,
			amount,
			description as reason,
			status,
			created_at
		FROM transactions
		WHERE metadata @> $1
		ORDER BY created_at DESC
	`

	metadataFilter := fmt.Sprintf(`{"original_transaction_id": "%s"}`, originalTxID.String())

	var history []*RefundHistoryItem
	if err := uc.db.SelectContext(ctx, &history, query, metadataFilter); err != nil {
		return nil, fmt.Errorf("failed to get refund history: %w", err)
	}

	return history, nil
}

// Helper
func stringPtr(s string) *string {
	return &s
}

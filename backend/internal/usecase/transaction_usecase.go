package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aryasatyawa/bayarin/internal/config"
	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/pkg/crypto"
	"github.com/aryasatyawa/bayarin/internal/pkg/validator"
	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionUsecase interface {
	Topup(ctx context.Context, req TopupRequest) (*TransactionResponse, error)
	Transfer(ctx context.Context, req TransferRequest) (*TransactionResponse, error)
	GetTransaction(ctx context.Context, transactionID uuid.UUID) (*TransactionDetail, error)
	GetUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*TransactionDetail, error)
}

type transactionUsecase struct {
	db         *sqlx.DB
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
	txRepo     repository.TransactionRepository
	ledgerRepo repository.LedgerRepository
	cfg        *config.Config
}

func NewTransactionUsecase(
	db *sqlx.DB,
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
	txRepo repository.TransactionRepository,
	ledgerRepo repository.LedgerRepository,
	cfg *config.Config,
) TransactionUsecase {
	return &transactionUsecase{
		db:         db,
		userRepo:   userRepo,
		walletRepo: walletRepo,
		txRepo:     txRepo,
		ledgerRepo: ledgerRepo,
		cfg:        cfg,
	}
}

// DTOs
type TopupRequest struct {
	UserID         uuid.UUID `json:"user_id"`
	Amount         int64     `json:"amount" validate:"required,gt=0"`
	ChannelCode    string    `json:"channel_code" validate:"required"`
	IdempotencyKey string    `json:"idempotency_key" validate:"required"`
}

type TransferRequest struct {
	UserID         uuid.UUID `json:"user_id"`
	ToUserID       uuid.UUID `json:"to_user_id" validate:"required"`
	Amount         int64     `json:"amount" validate:"required,gt=0"`
	Description    string    `json:"description"`
	PIN            string    `json:"pin" validate:"required,len=6"`
	IdempotencyKey string    `json:"idempotency_key" validate:"required"`
}

type TransactionResponse struct {
	TransactionID uuid.UUID                `json:"transaction_id"`
	Type          domain.TransactionType   `json:"type"`
	Amount        int64                    `json:"amount"`
	AmountIDR     string                   `json:"amount_idr"`
	Status        domain.TransactionStatus `json:"status"`
	Description   string                   `json:"description"`
	CreatedAt     time.Time                `json:"created_at"`
}

type TransactionDetail struct {
	ID           uuid.UUID                `json:"id"`
	Type         domain.TransactionType   `json:"type"`
	Amount       int64                    `json:"amount"`
	AmountIDR    string                   `json:"amount_idr"`
	Status       domain.TransactionStatus `json:"status"`
	FromWalletID *uuid.UUID               `json:"from_wallet_id,omitempty"`
	ToWalletID   *uuid.UUID               `json:"to_wallet_id,omitempty"`
	Description  string                   `json:"description"`
	ReferenceID  *string                  `json:"reference_id,omitempty"`
	CreatedAt    time.Time                `json:"created_at"`
	CompletedAt  *time.Time               `json:"completed_at,omitempty"`
}

// Topup handles wallet topup
func (uc *transactionUsecase) Topup(ctx context.Context, req TopupRequest) (*TransactionResponse, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := validator.ValidateAmount(req.Amount); err != nil {
		return nil, err
	}

	// Check idempotency
	existingTx, err := uc.txRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err == nil && existingTx != nil {
		return &TransactionResponse{
			TransactionID: existingTx.ID,
			Type:          existingTx.TransactionType,
			Amount:        existingTx.Amount,
			AmountIDR:     formatCurrency(existingTx.Amount),
			Status:        existingTx.Status,
			Description:   existingTx.Description,
			CreatedAt:     existingTx.CreatedAt,
		}, nil
	}

	wallet, err := uc.walletRepo.GetByUserIDAndType(ctx, req.UserID, domain.WalletTypeMain)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	if !wallet.IsActive() {
		return nil, domain.ErrWalletNotActive
	}

	tx, err := uc.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	wallet, err = uc.walletRepo.LockForUpdate(ctx, tx, wallet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to lock wallet: %w", err)
	}

	now := time.Now()
	metadata, _ := json.Marshal(map[string]interface{}{
		"channel_code": req.ChannelCode,
		"topup_method": "simulation",
	})

	transaction := &domain.Transaction{
		ID:              uuid.New(),
		IdempotencyKey:  req.IdempotencyKey,
		UserID:          req.UserID,
		TransactionType: domain.TransactionTypeTopup,
		Amount:          req.Amount,
		Currency:        uc.cfg.App.Currency,
		Status:          domain.TransactionStatusSuccess,
		ToWalletID:      &wallet.ID,
		ReferenceID:     stringPtr(fmt.Sprintf("TOPUP-%s", uuid.New().String()[:8])),
		Description:     fmt.Sprintf("Topup via %s", req.ChannelCode),
		Metadata:        metadata,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	transaction.MarkSuccess()

	if err := uc.txRepo.Create(ctx, tx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	ledgerEntry := domain.NewCreditEntry(
		transaction.ID,
		wallet.ID,
		req.Amount,
		wallet.Balance,
		transaction.Description,
	)

	if err := uc.ledgerRepo.CreateEntry(ctx, tx, ledgerEntry); err != nil {
		return nil, fmt.Errorf("failed to create ledger entry: %w", err)
	}

	newBalance := wallet.Balance + req.Amount
	if err := uc.walletRepo.UpdateBalance(ctx, tx, wallet.ID, newBalance); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &TransactionResponse{
		TransactionID: transaction.ID,
		Type:          transaction.TransactionType,
		Amount:        transaction.Amount,
		AmountIDR:     formatCurrency(transaction.Amount),
		Status:        transaction.Status,
		Description:   transaction.Description,
		CreatedAt:     transaction.CreatedAt,
	}, nil
}

// Transfer handles transfer between wallets
func (uc *transactionUsecase) Transfer(ctx context.Context, req TransferRequest) (*TransactionResponse, error) {
	// Validate input
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := validator.ValidateAmount(req.Amount); err != nil {
		return nil, err
	}

	// Get user and verify PIN
	user, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if user.PINHash == nil {
		return nil, domain.ErrInvalidPIN
	}

	// Verify PIN using crypto package
	if !crypto.VerifyPIN(req.PIN, *user.PINHash) {
		return nil, domain.ErrInvalidPIN
	}

	if req.UserID == req.ToUserID {
		return nil, domain.ErrSameWallet
	}

	// Check idempotency
	existingTx, err := uc.txRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err == nil && existingTx != nil {
		return &TransactionResponse{
			TransactionID: existingTx.ID,
			Type:          existingTx.TransactionType,
			Amount:        existingTx.Amount,
			AmountIDR:     formatCurrency(existingTx.Amount),
			Status:        existingTx.Status,
			Description:   existingTx.Description,
			CreatedAt:     existingTx.CreatedAt,
		}, nil
	}

	fromWallet, err := uc.walletRepo.GetByUserIDAndType(ctx, req.UserID, domain.WalletTypeMain)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender wallet: %w", err)
	}

	toWallet, err := uc.walletRepo.GetByUserIDAndType(ctx, req.ToUserID, domain.WalletTypeMain)
	if err != nil {
		return nil, fmt.Errorf("failed to get receiver wallet: %w", err)
	}

	if fromWallet.ID == toWallet.ID {
		return nil, domain.ErrSameWallet
	}

	if !fromWallet.IsActive() {
		return nil, domain.ErrWalletNotActive
	}

	if !toWallet.IsActive() {
		return nil, domain.ErrWalletNotActive
	}

	tx, err := uc.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock wallets in order
	var firstWallet, secondWallet *domain.Wallet
	if fromWallet.ID.String() < toWallet.ID.String() {
		firstWallet, err = uc.walletRepo.LockForUpdate(ctx, tx, fromWallet.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to lock first wallet: %w", err)
		}
		secondWallet, err = uc.walletRepo.LockForUpdate(ctx, tx, toWallet.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to lock second wallet: %w", err)
		}

		if firstWallet.ID == fromWallet.ID {
			fromWallet = firstWallet
			toWallet = secondWallet
		} else {
			fromWallet = secondWallet
			toWallet = firstWallet
		}
	} else {
		firstWallet, err = uc.walletRepo.LockForUpdate(ctx, tx, toWallet.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to lock first wallet: %w", err)
		}
		secondWallet, err = uc.walletRepo.LockForUpdate(ctx, tx, fromWallet.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to lock second wallet: %w", err)
		}

		if firstWallet.ID == toWallet.ID {
			toWallet = firstWallet
			fromWallet = secondWallet
		} else {
			toWallet = secondWallet
			fromWallet = firstWallet
		}
	}

	if !fromWallet.HasSufficientBalance(req.Amount) {
		return nil, domain.ErrInsufficientBalance
	}

	now := time.Now()
	description := req.Description
	if description == "" {
		description = "Transfer to user"
	}

	metadata, _ := json.Marshal(map[string]interface{}{
		"from_user_id": req.UserID.String(),
		"to_user_id":   req.ToUserID.String(),
	})

	transaction := &domain.Transaction{
		ID:              uuid.New(),
		IdempotencyKey:  req.IdempotencyKey,
		UserID:          req.UserID,
		TransactionType: domain.TransactionTypeTransfer,
		Amount:          req.Amount,
		Currency:        uc.cfg.App.Currency,
		Status:          domain.TransactionStatusSuccess,
		FromWalletID:    &fromWallet.ID,
		ToWalletID:      &toWallet.ID,
		ReferenceID:     stringPtr(fmt.Sprintf("TRF-%s", uuid.New().String()[:8])),
		Description:     description,
		Metadata:        metadata,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	transaction.MarkSuccess()

	if err := uc.txRepo.Create(ctx, tx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	debitEntry := domain.NewDebitEntry(
		transaction.ID,
		fromWallet.ID,
		req.Amount,
		fromWallet.Balance,
		fmt.Sprintf("Transfer out: %s", description),
	)

	creditEntry := domain.NewCreditEntry(
		transaction.ID,
		toWallet.ID,
		req.Amount,
		toWallet.Balance,
		fmt.Sprintf("Transfer in: %s", description),
	)

	ledgerEntries := []*domain.LedgerEntry{debitEntry, creditEntry}
	if err := uc.ledgerRepo.CreateEntries(ctx, tx, ledgerEntries); err != nil {
		return nil, fmt.Errorf("failed to create ledger entries: %w", err)
	}

	newFromBalance := fromWallet.Balance - req.Amount
	if err := uc.walletRepo.UpdateBalance(ctx, tx, fromWallet.ID, newFromBalance); err != nil {
		return nil, fmt.Errorf("failed to update sender wallet balance: %w", err)
	}

	newToBalance := toWallet.Balance + req.Amount
	if err := uc.walletRepo.UpdateBalance(ctx, tx, toWallet.ID, newToBalance); err != nil {
		return nil, fmt.Errorf("failed to update receiver wallet balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &TransactionResponse{
		TransactionID: transaction.ID,
		Type:          transaction.TransactionType,
		Amount:        transaction.Amount,
		AmountIDR:     formatCurrency(transaction.Amount),
		Status:        transaction.Status,
		Description:   transaction.Description,
		CreatedAt:     transaction.CreatedAt,
	}, nil
}

// GetTransaction returns transaction detail
func (uc *transactionUsecase) GetTransaction(ctx context.Context, transactionID uuid.UUID) (*TransactionDetail, error) {
	transaction, err := uc.txRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	return &TransactionDetail{
		ID:           transaction.ID,
		Type:         transaction.TransactionType,
		Amount:       transaction.Amount,
		AmountIDR:    formatCurrency(transaction.Amount),
		Status:       transaction.Status,
		FromWalletID: transaction.FromWalletID,
		ToWalletID:   transaction.ToWalletID,
		Description:  transaction.Description,
		ReferenceID:  transaction.ReferenceID,
		CreatedAt:    transaction.CreatedAt,
		CompletedAt:  transaction.CompletedAt,
	}, nil
}

// GetUserTransactions returns user transaction history
func (uc *transactionUsecase) GetUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*TransactionDetail, error) {
	transactions, err := uc.txRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}

	details := make([]*TransactionDetail, 0, len(transactions))
	for _, tx := range transactions {
		details = append(details, &TransactionDetail{
			ID:           tx.ID,
			Type:         tx.TransactionType,
			Amount:       tx.Amount,
			AmountIDR:    formatCurrency(tx.Amount),
			Status:       tx.Status,
			FromWalletID: tx.FromWalletID,
			ToWalletID:   tx.ToWalletID,
			Description:  tx.Description,
			ReferenceID:  tx.ReferenceID,
			CreatedAt:    tx.CreatedAt,
			CompletedAt:  tx.CompletedAt,
		})
	}

	return details, nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

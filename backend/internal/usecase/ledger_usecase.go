package usecase

import (
	"context"
	"fmt"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/google/uuid"
)

type LedgerUsecase interface {
	GetLedgerByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*domain.LedgerEntry, error)
	GetLedgerByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*domain.LedgerEntry, error)
}

type ledgerUsecase struct {
	ledgerRepo repository.LedgerRepository
}

func NewLedgerUsecase(ledgerRepo repository.LedgerRepository) LedgerUsecase {
	return &ledgerUsecase{
		ledgerRepo: ledgerRepo,
	}
}

// GetLedgerByTransactionID returns all ledger entries for a transaction
func (uc *ledgerUsecase) GetLedgerByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]*domain.LedgerEntry, error) {
	entries, err := uc.ledgerRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger by transaction: %w", err)
	}

	return entries, nil
}

// GetLedgerByWalletID returns ledger entries for a wallet with pagination
func (uc *ledgerUsecase) GetLedgerByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*domain.LedgerEntry, error) {
	entries, err := uc.ledgerRepo.GetByWalletID(ctx, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger by wallet: %w", err)
	}

	return entries, nil
}

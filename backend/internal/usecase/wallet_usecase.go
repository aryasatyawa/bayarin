package usecase

import (
	"context"
	"fmt"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/google/uuid"
)

type WalletUsecase interface {
	GetWalletBalance(ctx context.Context, userID uuid.UUID, walletType domain.WalletType) (*WalletBalance, error)
	GetAllWallets(ctx context.Context, userID uuid.UUID) ([]*WalletBalance, error)
	GetWalletHistory(ctx context.Context, walletID uuid.UUID, limit, offset int) (*WalletHistory, error)
}

type walletUsecase struct {
	walletRepo repository.WalletRepository
	ledgerRepo repository.LedgerRepository
}

func NewWalletUsecase(
	walletRepo repository.WalletRepository,
	ledgerRepo repository.LedgerRepository,
) WalletUsecase {
	return &walletUsecase{
		walletRepo: walletRepo,
		ledgerRepo: ledgerRepo,
	}
}

// DTOs
type WalletBalance struct {
	WalletID   uuid.UUID           `json:"wallet_id"`
	WalletType domain.WalletType   `json:"wallet_type"`
	Balance    int64               `json:"balance"`     // Integer (minor unit)
	BalanceIDR string              `json:"balance_idr"` // Formatted: "Rp 100.000"
	Currency   string              `json:"currency"`
	Status     domain.WalletStatus `json:"status"`
}

type WalletHistory struct {
	WalletID uuid.UUID             `json:"wallet_id"`
	Entries  []*domain.LedgerEntry `json:"entries"`
	Total    int                   `json:"total"`
	Limit    int                   `json:"limit"`
	Offset   int                   `json:"offset"`
}

// GetWalletBalance returns wallet balance
func (uc *walletUsecase) GetWalletBalance(ctx context.Context, userID uuid.UUID, walletType domain.WalletType) (*WalletBalance, error) {
	wallet, err := uc.walletRepo.GetByUserIDAndType(ctx, userID, walletType)
	if err != nil {
		return nil, err
	}

	return &WalletBalance{
		WalletID:   wallet.ID,
		WalletType: wallet.WalletType,
		Balance:    wallet.Balance,
		BalanceIDR: formatCurrency(wallet.Balance),
		Currency:   wallet.Currency,
		Status:     wallet.Status,
	}, nil
}

// GetAllWallets returns all user wallets
func (uc *walletUsecase) GetAllWallets(ctx context.Context, userID uuid.UUID) ([]*WalletBalance, error) {
	wallets, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	balances := make([]*WalletBalance, 0, len(wallets))
	for _, wallet := range wallets {
		balances = append(balances, &WalletBalance{
			WalletID:   wallet.ID,
			WalletType: wallet.WalletType,
			Balance:    wallet.Balance,
			BalanceIDR: formatCurrency(wallet.Balance),
			Currency:   wallet.Currency,
			Status:     wallet.Status,
		})
	}

	return balances, nil
}

// GetWalletHistory returns wallet transaction history
func (uc *walletUsecase) GetWalletHistory(ctx context.Context, walletID uuid.UUID, limit, offset int) (*WalletHistory, error) {
	entries, err := uc.ledgerRepo.GetByWalletID(ctx, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet history: %w", err)
	}

	return &WalletHistory{
		WalletID: walletID,
		Entries:  entries,
		Total:    len(entries),
		Limit:    limit,
		Offset:   offset,
	}, nil
}

// Helper function to format currency
func formatCurrency(amount int64) string {
	// Convert minor unit to major unit
	// Example: 10000000 (100000 rupiah dalam sen) -> "Rp 100.000"
	majorUnit := amount / 100

	// Format with thousand separator
	result := fmt.Sprintf("Rp %d", majorUnit)

	// Add thousand separator (simple implementation)
	// For production, use library like "golang.org/x/text/number"
	return result
}

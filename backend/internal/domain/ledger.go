package domain

import (
	"time"

	"github.com/google/uuid"
)

type LedgerEntry struct {
	ID            uuid.UUID `db:"id" json:"id"`
	TransactionID uuid.UUID `db:"transaction_id" json:"transaction_id"`
	WalletID      uuid.UUID `db:"wallet_id" json:"wallet_id"`
	EntryType     EntryType `db:"entry_type" json:"entry_type"`
	Amount        int64     `db:"amount" json:"amount"` // WAJIB INTEGER
	BalanceBefore int64     `db:"balance_before" json:"balance_before"`
	BalanceAfter  int64     `db:"balance_after" json:"balance_after"`
	Description   string    `db:"description" json:"description"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type EntryType string

const (
	EntryTypeDebit  EntryType = "debit"
	EntryTypeCredit EntryType = "credit"
)

// NewDebitEntry creates new debit ledger entry
func NewDebitEntry(transactionID, walletID uuid.UUID, amount, balanceBefore int64, description string) *LedgerEntry {
	return &LedgerEntry{
		ID:            uuid.New(),
		TransactionID: transactionID,
		WalletID:      walletID,
		EntryType:     EntryTypeDebit,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceBefore - amount, // Debit mengurangi saldo
		Description:   description,
		CreatedAt:     time.Now(),
	}
}

// NewCreditEntry creates new credit ledger entry
func NewCreditEntry(transactionID, walletID uuid.UUID, amount, balanceBefore int64, description string) *LedgerEntry {
	return &LedgerEntry{
		ID:            uuid.New(),
		TransactionID: transactionID,
		WalletID:      walletID,
		EntryType:     EntryTypeCredit,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceBefore + amount, // Credit menambah saldo
		Description:   description,
		CreatedAt:     time.Now(),
	}
}

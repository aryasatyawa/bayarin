package domain

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID         uuid.UUID    `db:"id" json:"id"`
	UserID     uuid.UUID    `db:"user_id" json:"user_id"`
	WalletType WalletType   `db:"wallet_type" json:"wallet_type"`
	Balance    int64        `db:"balance" json:"balance"` // WAJIB INTEGER (minor unit)
	Currency   string       `db:"currency" json:"currency"`
	Status     WalletStatus `db:"status" json:"status"`
	CreatedAt  time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at" json:"updated_at"`
}

type WalletType string

const (
	WalletTypeMain     WalletType = "main"
	WalletTypeBonus    WalletType = "bonus"
	WalletTypeCashback WalletType = "cashback"
)

type WalletStatus string

const (
	WalletStatusActive WalletStatus = "active"
	WalletStatusFrozen WalletStatus = "frozen"
	WalletStatusClosed WalletStatus = "closed"
)

// IsActive checks if wallet is active
func (w *Wallet) IsActive() bool {
	return w.Status == WalletStatusActive
}

// HasSufficientBalance checks if wallet has enough balance
func (w *Wallet) HasSufficientBalance(amount int64) bool {
	return w.Balance >= amount
}

// CanDebit checks if wallet can be debited
func (w *Wallet) CanDebit(amount int64) error {
	if !w.IsActive() {
		return ErrWalletNotActive
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if !w.HasSufficientBalance(amount) {
		return ErrInsufficientBalance
	}
	return nil
}

// CanCredit checks if wallet can be credited
func (w *Wallet) CanCredit(amount int64) error {
	if !w.IsActive() {
		return ErrWalletNotActive
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

// Debit mengurangi balance (PENTING: tidak langsung update DB, hanya kalkulasi)
func (w *Wallet) Debit(amount int64) error {
	if err := w.CanDebit(amount); err != nil {
		return err
	}
	w.Balance -= amount
	return nil
}

// Credit menambah balance (PENTING: tidak langsung update DB, hanya kalkulasi)
func (w *Wallet) Credit(amount int64) error {
	if err := w.CanCredit(amount); err != nil {
		return err
	}
	w.Balance += amount
	return nil
}

package domain

import "errors"

var (
	// User errors
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exists")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrInvalidPhone     = errors.New("invalid phone")
	ErrInvalidFullName  = errors.New("invalid full name")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidPIN       = errors.New("invalid PIN")
	ErrUserNotActive    = errors.New("user is not active")

	// Wallet errors
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrWalletAlreadyExist  = errors.New("wallet already exists")
	ErrWalletNotActive     = errors.New("wallet is not active")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrSameWallet          = errors.New("cannot transfer to same wallet")

	// Transaction errors
	ErrTransactionNotFound    = errors.New("transaction not found")
	ErrDuplicateTransaction   = errors.New("duplicate transaction")
	ErrTransactionFailed      = errors.New("transaction failed")
	ErrInvalidTransactionType = errors.New("invalid transaction type")

	// General errors
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrDatabaseOperation = errors.New("database operation failed")
)

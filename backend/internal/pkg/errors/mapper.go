package errors

import (
	"errors"
	"net/http"

	"github.com/aryasatyawa/bayarin/internal/domain"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MapError maps domain error to HTTP status code and error response
func MapError(err error) (int, ErrorResponse) {
	// User errors
	if errors.Is(err, domain.ErrUserNotFound) {
		return http.StatusNotFound, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		}
	}
	if errors.Is(err, domain.ErrUserAlreadyExist) {
		return http.StatusConflict, ErrorResponse{
			Code:    "USER_ALREADY_EXISTS",
			Message: "User already exists",
		}
	}
	if errors.Is(err, domain.ErrInvalidEmail) {
		return http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_EMAIL",
			Message: "Invalid email format",
		}
	}
	if errors.Is(err, domain.ErrInvalidPhone) {
		return http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PHONE",
			Message: "Invalid phone number",
		}
	}
	if errors.Is(err, domain.ErrInvalidPassword) {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_PASSWORD",
			Message: "Invalid password",
		}
	}
	if errors.Is(err, domain.ErrInvalidPIN) {
		return http.StatusUnauthorized, ErrorResponse{
			Code:    "INVALID_PIN",
			Message: "Invalid PIN",
		}
	}
	if errors.Is(err, domain.ErrUserNotActive) {
		return http.StatusForbidden, ErrorResponse{
			Code:    "USER_NOT_ACTIVE",
			Message: "User account is not active",
		}
	}

	// Wallet errors
	if errors.Is(err, domain.ErrWalletNotFound) {
		return http.StatusNotFound, ErrorResponse{
			Code:    "WALLET_NOT_FOUND",
			Message: "Wallet not found",
		}
	}
	if errors.Is(err, domain.ErrWalletNotActive) {
		return http.StatusForbidden, ErrorResponse{
			Code:    "WALLET_NOT_ACTIVE",
			Message: "Wallet is not active",
		}
	}
	if errors.Is(err, domain.ErrInsufficientBalance) {
		return http.StatusBadRequest, ErrorResponse{
			Code:    "INSUFFICIENT_BALANCE",
			Message: "Insufficient wallet balance",
		}
	}
	if errors.Is(err, domain.ErrInvalidAmount) {
		return http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_AMOUNT",
			Message: "Invalid transaction amount",
		}
	}
	if errors.Is(err, domain.ErrSameWallet) {
		return http.StatusBadRequest, ErrorResponse{
			Code:    "SAME_WALLET",
			Message: "Cannot transfer to same wallet",
		}
	}

	// Transaction errors
	if errors.Is(err, domain.ErrTransactionNotFound) {
		return http.StatusNotFound, ErrorResponse{
			Code:    "TRANSACTION_NOT_FOUND",
			Message: "Transaction not found",
		}
	}
	if errors.Is(err, domain.ErrDuplicateTransaction) {
		return http.StatusConflict, ErrorResponse{
			Code:    "DUPLICATE_TRANSACTION",
			Message: "Duplicate transaction detected",
		}
	}

	// Default error
	return http.StatusInternalServerError, ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: err.Error(),
	}
}

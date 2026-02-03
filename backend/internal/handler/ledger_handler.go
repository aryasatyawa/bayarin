package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/pkg/errors"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/aryasatyawa/bayarin/internal/usecase"
)

type LedgerHandler struct {
	ledgerViewerUsecase usecase.LedgerViewerUsecase
}

func NewLedgerHandler(ledgerViewerUsecase usecase.LedgerViewerUsecase) *LedgerHandler {
	return &LedgerHandler{
		ledgerViewerUsecase: ledgerViewerUsecase,
	}
}

// ==============================
// Get Ledger Entries (Filtered)
// ==============================
func (h *LedgerHandler) GetLedgerEntries(c *gin.Context) {
	filter := usecase.LedgerFilter{
		Limit:  20,
		Offset: 0,
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			response.BadRequest(c, "Invalid user_id", err.Error())
			return
		}
		filter.UserID = &userID
	}

	if walletIDStr := c.Query("wallet_id"); walletIDStr != "" {
		walletID, err := uuid.Parse(walletIDStr)
		if err != nil {
			response.BadRequest(c, "Invalid wallet_id", err.Error())
			return
		}
		filter.WalletID = &walletID
	}

	if txIDStr := c.Query("transaction_id"); txIDStr != "" {
		txID, err := uuid.Parse(txIDStr)
		if err != nil {
			response.BadRequest(c, "Invalid transaction_id", err.Error())
			return
		}
		filter.TransactionID = &txID
	}

	if entryTypeStr := c.Query("entry_type"); entryTypeStr != "" {
		entryType := domain.EntryType(entryTypeStr)
		filter.EntryType = &entryType
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format (YYYY-MM-DD)", err.Error())
			return
		}
		filter.StartDate = &startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format (YYYY-MM-DD)", err.Error())
			return
		}
		filter.EndDate = &endDate
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	result, err := h.ledgerViewerUsecase.GetLedgerEntries(c.Request.Context(), filter)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Ledger entries retrieved successfully", result)
}

// ==============================
// Get Ledger by Transaction
// ==============================
func (h *LedgerHandler) GetLedgerByTransaction(c *gin.Context) {
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid transaction ID", err.Error())
		return
	}

	entries, err := h.ledgerViewerUsecase.GetLedgerByTransactionID(c.Request.Context(), txID)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Ledger entries retrieved successfully", entries)
}

// ==============================
// Get Ledger by Wallet
// ==============================
func (h *LedgerHandler) GetLedgerByWallet(c *gin.Context) {
	walletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid wallet ID", err.Error())
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result, err := h.ledgerViewerUsecase.GetLedgerByWalletID(
		c.Request.Context(),
		walletID,
		limit,
		offset,
	)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Wallet ledger retrieved successfully", result)
}

// ==============================
// Validate Wallet Balance
// ==============================
func (h *LedgerHandler) ValidateBalance(c *gin.Context) {
	walletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid wallet ID", err.Error())
		return
	}

	result, err := h.ledgerViewerUsecase.ValidateBalance(c.Request.Context(), walletID)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Balance validation completed", result)
}

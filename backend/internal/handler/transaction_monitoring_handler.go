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

type TransactionMonitoringHandler struct {
	txMonitoringUsecase usecase.TransactionMonitoringUsecase
}

func NewTransactionMonitoringHandler(
	txMonitoringUsecase usecase.TransactionMonitoringUsecase,
) *TransactionMonitoringHandler {
	return &TransactionMonitoringHandler{
		txMonitoringUsecase: txMonitoringUsecase,
	}
}

// ==============================
// Get All Transactions
// ==============================
func (h *TransactionMonitoringHandler) GetAllTransactions(c *gin.Context) {
	filter := usecase.TransactionFilter{
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

	if txTypeStr := c.Query("transaction_type"); txTypeStr != "" {
		txType := domain.TransactionType(txTypeStr)
		filter.TransactionType = &txType
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.TransactionStatus(statusStr)
		filter.Status = &status
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

	if minAmountStr := c.Query("min_amount"); minAmountStr != "" {
		if minAmount, err := strconv.ParseInt(minAmountStr, 10, 64); err == nil {
			filter.MinAmount = &minAmount
		}
	}

	if maxAmountStr := c.Query("max_amount"); maxAmountStr != "" {
		if maxAmount, err := strconv.ParseInt(maxAmountStr, 10, 64); err == nil {
			filter.MaxAmount = &maxAmount
		}
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

	result, err := h.txMonitoringUsecase.GetAllTransactions(c.Request.Context(), filter)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Transactions retrieved successfully", result)
}

// ==============================
// Get Transaction Detail
// ==============================
func (h *TransactionMonitoringHandler) GetTransactionDetail(c *gin.Context) {
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid transaction ID", err.Error())
		return
	}

	result, err := h.txMonitoringUsecase.GetTransactionDetail(c.Request.Context(), txID)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Transaction detail retrieved successfully", result)
}

// ==============================
// Pending Transactions
// ==============================
func (h *TransactionMonitoringHandler) GetPendingTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result, err := h.txMonitoringUsecase.GetPendingTransactions(
		c.Request.Context(),
		limit,
		offset,
	)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Pending transactions retrieved successfully", result)
}

// ==============================
// Failed Transactions
// ==============================
func (h *TransactionMonitoringHandler) GetFailedTransactions(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result, err := h.txMonitoringUsecase.GetFailedTransactions(
		c.Request.Context(),
		days,
		limit,
		offset,
	)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Failed transactions retrieved successfully", result)
}

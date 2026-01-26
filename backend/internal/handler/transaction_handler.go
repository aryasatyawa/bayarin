package handler

import (
	"strconv"

	"github.com/aryasatyawa/bayarin/internal/middleware"
	"github.com/aryasatyawa/bayarin/internal/pkg/errors"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/aryasatyawa/bayarin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
	}
}

// Topup godoc
// @Summary Topup wallet
// @Description Topup wallet balance via payment channel
// @Tags transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TopupRequestDTO true "Topup request"
// @Success 200 {object} response.Response{data=usecase.TransactionResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /transaction/topup [post]
func (h *TransactionHandler) Topup(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req TopupRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Build usecase request
	topupReq := usecase.TopupRequest{
		UserID:         userID,
		Amount:         req.Amount,
		ChannelCode:    req.ChannelCode,
		IdempotencyKey: req.IdempotencyKey,
	}

	result, err := h.transactionUsecase.Topup(c.Request.Context(), topupReq)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Topup successful", result)
}

// Transfer godoc
// @Summary Transfer to another user
// @Description Transfer money to another user's wallet
// @Tags transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TransferRequestDTO true "Transfer request"
// @Success 200 {object} response.Response{data=usecase.TransactionResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /transaction/transfer [post]
func (h *TransactionHandler) Transfer(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req TransferRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Parse to_user_id
	toUserID, err := uuid.Parse(req.ToUserID)
	if err != nil {
		response.BadRequest(c, "Invalid to_user_id", err.Error())
		return
	}

	// Build usecase request
	transferReq := usecase.TransferRequest{
		UserID:         userID,
		ToUserID:       toUserID,
		Amount:         req.Amount,
		Description:    req.Description,
		PIN:            req.PIN,
		IdempotencyKey: req.IdempotencyKey,
	}

	result, err := h.transactionUsecase.Transfer(c.Request.Context(), transferReq)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Transfer successful", result)
}

// GetTransaction godoc
// @Summary Get transaction detail
// @Description Get transaction detail by ID
// @Tags transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.Response{data=usecase.TransactionDetail}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /transaction/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	// Get transaction ID from path
	txIDStr := c.Param("id")
	txID, err := uuid.Parse(txIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid transaction ID", err.Error())
		return
	}

	transaction, err := h.transactionUsecase.GetTransaction(c.Request.Context(), txID)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Transaction retrieved successfully", transaction)
}

// GetUserTransactions godoc
// @Summary Get user transaction history
// @Description Get authenticated user's transaction history
// @Tags transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=[]usecase.TransactionDetail}
// @Failure 401 {object} response.Response
// @Router /transaction/history [get]
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get pagination params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.transactionUsecase.GetUserTransactions(c.Request.Context(), userID, limit, offset)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.SuccessWithMeta(c, "Transactions retrieved successfully", transactions, gin.H{
		"limit":  limit,
		"offset": offset,
		"count":  len(transactions),
	})
}

// Request DTOs
type TopupRequestDTO struct {
	Amount         int64  `json:"amount" binding:"required,gt=0"`
	ChannelCode    string `json:"channel_code" binding:"required"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

type TransferRequestDTO struct {
	ToUserID       string `json:"to_user_id" binding:"required"`
	Amount         int64  `json:"amount" binding:"required,gt=0"`
	Description    string `json:"description"`
	PIN            string `json:"pin" binding:"required,len=6"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

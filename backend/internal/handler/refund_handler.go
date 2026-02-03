package handler

import (
	"github.com/aryasatyawa/bayarin/internal/middleware"
	"github.com/aryasatyawa/bayarin/internal/pkg/errors"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/aryasatyawa/bayarin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RefundHandler struct {
	refundUsecase usecase.RefundUsecase
}

func NewRefundHandler(refundUsecase usecase.RefundUsecase) *RefundHandler {
	return &RefundHandler{
		refundUsecase: refundUsecase,
	}
}

// RefundTransaction godoc
// @Summary Refund transaction
// @Description Create refund for a transaction (finance admin only)
// @Tags admin-refund
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecase.RefundRequest true "Refund request"
// @Success 200 {object} response.Response{data=usecase.RefundResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/refund [post]
func (h *RefundHandler) RefundTransaction(c *gin.Context) {
	adminID, err := middleware.GetAdminID(c)
	if err != nil {
		response.Unauthorized(c, "Admin not authenticated")
		return
	}

	var req usecase.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.refundUsecase.RefundTransaction(c.Request.Context(), adminID, req)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Refund processed successfully", result)
}

// ReverseTransaction godoc
// @Summary Reverse transaction
// @Description Reverse a transaction completely (finance admin only)
// @Tags admin-refund
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecase.ReverseRequest true "Reverse request"
// @Success 200 {object} response.Response{data=usecase.RefundResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/refund/reverse [post]
func (h *RefundHandler) ReverseTransaction(c *gin.Context) {
	adminID, err := middleware.GetAdminID(c)
	if err != nil {
		response.Unauthorized(c, "Admin not authenticated")
		return
	}

	var req usecase.ReverseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.refundUsecase.ReverseTransaction(c.Request.Context(), adminID, req)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Transaction reversed successfully", result)
}

// GetRefundHistory godoc
// @Summary Get refund history
// @Description Get refund history for original transaction
// @Tags admin-refund
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Original Transaction ID"
// @Success 200 {object} response.Response{data=[]usecase.RefundHistoryItem}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/refund/history/{id} [get]
func (h *RefundHandler) GetRefundHistory(c *gin.Context) {
	txIDStr := c.Param("id")
	txID, err := uuid.Parse(txIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid transaction ID", err.Error())
		return
	}

	history, err := h.refundUsecase.GetRefundHistory(c.Request.Context(), txID)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Refund history retrieved successfully", history)
}

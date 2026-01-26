package handler

import (
	"strconv"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/middleware"
	"github.com/aryasatyawa/bayarin/internal/pkg/errors"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/aryasatyawa/bayarin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecase
}

func NewWalletHandler(walletUsecase usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{
		walletUsecase: walletUsecase,
	}
}

// GetBalance godoc
// @Summary Get wallet balance
// @Description Get balance for specific wallet type
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Wallet type (main, bonus, cashback)" default(main)
// @Success 200 {object} response.Response{data=usecase.WalletBalance}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /wallet/balance [get]
func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get wallet type from query (default: main)
	walletTypeStr := c.DefaultQuery("type", string(domain.WalletTypeMain))
	walletType := domain.WalletType(walletTypeStr)

	balance, err := h.walletUsecase.GetWalletBalance(c.Request.Context(), userID, walletType)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Balance retrieved successfully", balance)
}

// GetAllWallets godoc
// @Summary Get all wallets
// @Description Get all wallets for authenticated user
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]usecase.WalletBalance}
// @Failure 401 {object} response.Response
// @Router /wallet/all [get]
func (h *WalletHandler) GetAllWallets(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	wallets, err := h.walletUsecase.GetAllWallets(c.Request.Context(), userID)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Wallets retrieved successfully", wallets)
}

// GetHistory godoc
// @Summary Get wallet history
// @Description Get transaction history for specific wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param wallet_id path string true "Wallet ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=usecase.WalletHistory}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /wallet/{wallet_id}/history [get]
func (h *WalletHandler) GetHistory(c *gin.Context) {
	// Get wallet ID from path
	walletIDStr := c.Param("wallet_id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid wallet ID", err.Error())
		return
	}

	// Get pagination params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	history, err := h.walletUsecase.GetWalletHistory(c.Request.Context(), walletID, limit, offset)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "History retrieved successfully", history)
}

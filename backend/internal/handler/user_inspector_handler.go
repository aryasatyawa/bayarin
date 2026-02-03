package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aryasatyawa/bayarin/internal/middleware"
	"github.com/aryasatyawa/bayarin/internal/pkg/errors"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/aryasatyawa/bayarin/internal/usecase"
)

type UserInspectorHandler struct {
	userInspectorUsecase usecase.UserInspectorUsecase
}

func NewUserInspectorHandler(
	userInspectorUsecase usecase.UserInspectorUsecase,
) *UserInspectorHandler {
	return &UserInspectorHandler{
		userInspectorUsecase: userInspectorUsecase,
	}
}

// GetUserDetails godoc
// @Summary Get user details
// @Description Get comprehensive user details including wallets and statistics
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=usecase.UserInspectorDetail}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/users/{id} [get]
func (h *UserInspectorHandler) GetUserDetails(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err.Error())
		return
	}

	result, err := h.userInspectorUsecase.GetUserDetails(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "User details retrieved successfully", result)
}

// SearchUsers godoc
// @Summary Search users
// @Description Search users by email, phone number, or name
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=[]usecase.UserSearchResult}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/users/search [get]
func (h *UserInspectorHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Search query is required", nil)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.userInspectorUsecase.SearchUsers(
		c.Request.Context(),
		query,
		limit,
		offset,
	)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Users retrieved successfully", users)
}

// FreezeWallet godoc
// @Summary Freeze wallet
// @Description Freeze a wallet to prevent transactions (ops admin only)
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Param request body FreezeWalletRequest true "Freeze request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/wallets/{id}/freeze [post]
func (h *UserInspectorHandler) FreezeWallet(c *gin.Context) {
	adminID, err := middleware.GetAdminID(c)
	if err != nil {
		response.Unauthorized(c, "Admin not authenticated")
		return
	}

	walletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid wallet ID", err.Error())
		return
	}

	var req FreezeWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.userInspectorUsecase.FreezeWallet(
		c.Request.Context(),
		adminID,
		walletID,
		req.Reason,
	); err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Wallet frozen successfully", nil)
}

// UnfreezeWallet godoc
// @Summary Unfreeze wallet
// @Description Unfreeze a wallet to allow transactions (ops admin only)
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Param request body UnfreezeWalletRequest true "Unfreeze request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/wallets/{id}/unfreeze [post]
func (h *UserInspectorHandler) UnfreezeWallet(c *gin.Context) {
	adminID, err := middleware.GetAdminID(c)
	if err != nil {
		response.Unauthorized(c, "Admin not authenticated")
		return
	}

	walletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid wallet ID", err.Error())
		return
	}

	var req UnfreezeWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.userInspectorUsecase.UnfreezeWallet(
		c.Request.Context(),
		adminID,
		walletID,
		req.Reason,
	); err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Wallet unfrozen successfully", nil)
}

type FreezeWalletRequest struct {
	Reason string `json:"reason" binding:"required,min=10"`
}

type UnfreezeWalletRequest struct {
	Reason string `json:"reason" binding:"required,min=10"`
}

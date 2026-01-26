package handler

import (
	"github.com/aryasatyawa/bayarin/internal/middleware"
	"github.com/aryasatyawa/bayarin/internal/pkg/errors"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/aryasatyawa/bayarin/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// Register godoc
// @Summary Register new user
// @Description Register new user and create main wallet
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecase.RegisterRequest true "Register request"
// @Success 201 {object} response.Response{data=usecase.RegisterResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req usecase.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.userUsecase.Register(c.Request.Context(), req)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Created(c, "User registered successfully", result)
}

// Login godoc
// @Summary User login
// @Description Login with email/phone and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecase.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=usecase.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.userUsecase.Login(c.Request.Context(), req)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Login successful", result)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user profile with wallets
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=usecase.UserProfile}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	profile, err := h.userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Profile retrieved successfully", profile)
}

// SetPIN godoc
// @Summary Set transaction PIN
// @Description Set 6-digit PIN for transaction authorization
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SetPINRequest true "PIN request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /user/pin [post]
func (h *UserHandler) SetPIN(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req SetPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.userUsecase.SetPIN(c.Request.Context(), userID, req.PIN); err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "PIN set successfully", nil)
}

// VerifyPIN godoc
// @Summary Verify transaction PIN
// @Description Verify user's transaction PIN
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body VerifyPINRequest true "PIN request"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /user/pin/verify [post]
func (h *UserHandler) VerifyPIN(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.userUsecase.VerifyPIN(c.Request.Context(), userID, req.PIN); err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "PIN verified successfully", nil)
}

// Request DTOs
type SetPINRequest struct {
	PIN string `json:"pin" binding:"required,len=6"`
}

type VerifyPINRequest struct {
	PIN string `json:"pin" binding:"required,len=6"`
}

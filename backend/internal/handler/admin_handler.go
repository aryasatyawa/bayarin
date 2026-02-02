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

type AdminHandler struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		adminUsecase: adminUsecase,
	}
}

// Login godoc
// @Summary Admin login
// @Description Authenticate admin and get JWT token
// @Tags admin
// @Accept json
// @Produce json
// @Param request body usecase.AdminLoginRequest true "Login request"
// @Success 200 {object} response.Response{data=usecase.AdminLoginResponse}
// @Failure 401 {object} response.Response
// @Router /admin/auth/login [post]
func (h *AdminHandler) Login(c *gin.Context) {
	var req usecase.AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.adminUsecase.Login(c.Request.Context(), req)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Login successful", result)
}

// CreateAdmin godoc
// @Summary Create new admin
// @Description Create new admin (super admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecase.CreateAdminRequest true "Create admin request"
// @Success 201 {object} response.Response{data=usecase.AdminResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/admins [post]
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	adminID, err := middleware.GetAdminID(c)
	if err != nil {
		response.Unauthorized(c, "Admin not authenticated")
		return
	}

	var req usecase.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.adminUsecase.CreateAdmin(c.Request.Context(), adminID, req)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Created(c, "Admin created successfully", result)
}

// GetAdmin godoc
// @Summary Get admin by ID
// @Description Get admin details by ID
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Admin ID"
// @Success 200 {object} response.Response{data=usecase.AdminResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/admins/{id} [get]
func (h *AdminHandler) GetAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid admin ID", err.Error())
		return
	}

	result, err := h.adminUsecase.GetAdminByID(c.Request.Context(), id)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Admin retrieved successfully", result)
}

// ListAdmins godoc
// @Summary List all admins
// @Description Get list of all admins with pagination
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=[]usecase.AdminResponse}
// @Failure 401 {object} response.Response
// @Router /admin/admins [get]
func (h *AdminHandler) ListAdmins(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	admins, err := h.adminUsecase.ListAdmins(c.Request.Context(), limit, offset)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Admins retrieved successfully", admins)
}

// UpdateAdminStatus godoc
// @Summary Update admin status
// @Description Suspend or activate admin (super admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Admin ID"
// @Param request body UpdateAdminStatusRequest true "Status update request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/admins/{id}/status [patch]
func (h *AdminHandler) UpdateAdminStatus(c *gin.Context) {
	actorID, err := middleware.GetAdminID(c)
	if err != nil {
		response.Unauthorized(c, "Admin not authenticated")
		return
	}

	targetIDStr := c.Param("id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid admin ID", err.Error())
		return
	}

	var req UpdateAdminStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	err = h.adminUsecase.UpdateAdminStatus(c.Request.Context(), actorID, targetID, req.Status)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Admin status updated successfully", nil)
}

// GetAuditLogs godoc
// @Summary Get audit logs
// @Description Get audit logs for current admin
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=[]domain.AuditLog}
// @Failure 401 {object} response.Response
// @Router /admin/audit-logs [get]
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	adminID, err := middleware.GetAdminID(c)
	if err != nil {
		response.Unauthorized(c, "Admin not authenticated")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, err := h.adminUsecase.GetAuditLogs(c.Request.Context(), adminID, limit, offset)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Audit logs retrieved successfully", logs)
}

// Request DTOs
type UpdateAdminStatusRequest struct {
	Status domain.AdminStatus `json:"status" binding:"required"`
}

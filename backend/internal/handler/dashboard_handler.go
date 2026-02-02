package handler

import (
	"time"

	"github.com/aryasatyawa/bayarin/internal/pkg/errors"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/aryasatyawa/bayarin/internal/usecase"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardUsecase usecase.DashboardUsecase
}

func NewDashboardHandler(dashboardUsecase usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{
		dashboardUsecase: dashboardUsecase,
	}
}

// GetOverview godoc
// @Summary Get dashboard overview
// @Description Get dashboard overview with key metrics
// @Tags admin-dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=usecase.DashboardOverview}
// @Failure 401 {object} response.Response
// @Router /admin/dashboard/overview [get]
func (h *DashboardHandler) GetOverview(c *gin.Context) {
	overview, err := h.dashboardUsecase.GetOverview(c.Request.Context())
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Overview retrieved successfully", overview)
}

// GetDailyStats godoc
// @Summary Get daily statistics
// @Description Get transaction statistics for specific date
// @Tags admin-dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string false "Date (YYYY-MM-DD)" default(today)
// @Success 200 {object} response.Response{data=usecase.DailyStats}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/dashboard/daily-stats [get]
func (h *DashboardHandler) GetDailyStats(c *gin.Context) {
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "Invalid date format. Use YYYY-MM-DD", err.Error())
		return
	}

	stats, err := h.dashboardUsecase.GetDailyStats(c.Request.Context(), date)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Daily stats retrieved successfully", stats)
}

// GetTransactionSummary godoc
// @Summary Get transaction summary
// @Description Get transaction summary for date range
// @Tags admin-dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)" default(7 days ago)
// @Param end_date query string false "End date (YYYY-MM-DD)" default(today)
// @Success 200 {object} response.Response{data=usecase.TransactionSummary}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/dashboard/transaction-summary [get]
func (h *DashboardHandler) GetTransactionSummary(c *gin.Context) {
	// Default: last 7 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	if startStr := c.Query("start_date"); startStr != "" {
		parsed, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format. Use YYYY-MM-DD", err.Error())
			return
		}
		startDate = parsed
	}

	if endStr := c.Query("end_date"); endStr != "" {
		parsed, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format. Use YYYY-MM-DD", err.Error())
			return
		}
		endDate = parsed
	}

	// Validate date range
	if startDate.After(endDate) {
		response.BadRequest(c, "start_date must be before end_date", nil)
		return
	}

	summary, err := h.dashboardUsecase.GetTransactionSummary(c.Request.Context(), startDate, endDate)
	if err != nil {
		statusCode, errResp := errors.MapError(err)
		response.Error(c, statusCode, errResp.Message, errResp)
		return
	}

	response.Success(c, "Transaction summary retrieved successfully", summary)
}

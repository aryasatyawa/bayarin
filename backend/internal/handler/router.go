package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/aryasatyawa/bayarin/internal/middleware"
	"github.com/aryasatyawa/bayarin/internal/pkg/jwt"
)

type Router struct {
	engine             *gin.Engine
	userHandler        *UserHandler
	walletHandler      *WalletHandler
	transactionHandler *TransactionHandler
	healthHandler      *HealthHandler

	// Admin
	adminHandler     *AdminHandler
	dashboardHandler *DashboardHandler

	tokenManager *jwt.TokenManager
}

func NewRouter(
	userHandler *UserHandler,
	walletHandler *WalletHandler,
	transactionHandler *TransactionHandler,
	healthHandler *HealthHandler,
	adminHandler *AdminHandler,
	dashboardHandler *DashboardHandler,
	tokenManager *jwt.TokenManager,
) *Router {
	return &Router{
		engine:             gin.Default(),
		userHandler:        userHandler,
		walletHandler:      walletHandler,
		transactionHandler: transactionHandler,
		healthHandler:      healthHandler,
		adminHandler:       adminHandler,
		dashboardHandler:   dashboardHandler,
		tokenManager:       tokenManager,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Global middleware
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.LoggerMiddleware())

	// =====================
	// API v1 - USER
	// =====================
	v1 := r.engine.Group("/api/v1")
	{
		// Health (public)
		v1.GET("/health", r.healthHandler.Health)

		// Auth (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userHandler.Register)
			auth.POST("/login", r.userHandler.Login)
		}

		// Protected user routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(r.tokenManager))
		{
			// User
			user := protected.Group("/user")
			{
				user.GET("/profile", r.userHandler.GetProfile)
				user.POST("/pin", r.userHandler.SetPIN)
				user.POST("/pin/verify", r.userHandler.VerifyPIN)
			}

			// Wallet
			wallet := protected.Group("/wallet")
			{
				wallet.GET("/balance", r.walletHandler.GetBalance)
				wallet.GET("/all", r.walletHandler.GetAllWallets)
				wallet.GET("/:wallet_id/history", r.walletHandler.GetHistory)
			}

			// Transaction
			transaction := protected.Group("/transaction")
			{
				transaction.POST("/topup", r.transactionHandler.Topup)
				transaction.POST("/transfer", r.transactionHandler.Transfer)
				transaction.GET("/:id", r.transactionHandler.GetTransaction)
				transaction.GET("/history", r.transactionHandler.GetUserTransactions)
			}
		}
	}

	// =====================
	// API v1 - ADMIN
	// =====================
	admin := r.engine.Group("/api/v1/admin")
	{
		// Admin auth (public)
		adminAuth := admin.Group("/auth")
		{
			adminAuth.POST("/login", r.adminHandler.Login)
		}

		// Admin protected
		adminProtected := admin.Group("")
		adminProtected.Use(middleware.AdminAuthMiddleware(r.tokenManager))
		{
			// Dashboard
			dashboard := adminProtected.Group("/dashboard")
			{
				dashboard.GET("/overview", r.dashboardHandler.GetOverview)
				dashboard.GET("/daily-stats", r.dashboardHandler.GetDailyStats)
				dashboard.GET("/transaction-summary", r.dashboardHandler.GetTransactionSummary)
			}

			// Admin management (SUPER ADMIN only)
			admins := adminProtected.Group("/admins")
			admins.Use(middleware.RequireSuperAdmin())
			{
				admins.POST("", r.adminHandler.CreateAdmin)
				admins.GET("", r.adminHandler.ListAdmins)
				admins.GET("/:id", r.adminHandler.GetAdmin)
				admins.PATCH("/:id/status", r.adminHandler.UpdateAdminStatus)
			}

			// Audit logs
			adminProtected.GET("/audit-logs", r.adminHandler.GetAuditLogs)
		}
	}

	return r.engine
}

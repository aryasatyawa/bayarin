package handler

import (
	"github.com/aryasatyawa/bayarin/internal/middleware"
	"github.com/aryasatyawa/bayarin/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	// User handlers
	userHandler        *UserHandler
	walletHandler      *WalletHandler
	transactionHandler *TransactionHandler
	healthHandler      *HealthHandler
	// Admin handlers
	adminHandler                 *AdminHandler
	dashboardHandler             *DashboardHandler
	ledgerHandler                *LedgerHandler
	transactionMonitoringHandler *TransactionMonitoringHandler
	refundHandler                *RefundHandler
	userInspectorHandler         *UserInspectorHandler
	tokenManager                 *jwt.TokenManager
}

func NewRouter(
	userHandler *UserHandler,
	walletHandler *WalletHandler,
	transactionHandler *TransactionHandler,
	healthHandler *HealthHandler,
	adminHandler *AdminHandler,
	dashboardHandler *DashboardHandler,
	ledgerHandler *LedgerHandler,
	transactionMonitoringHandler *TransactionMonitoringHandler,
	refundHandler *RefundHandler,
	userInspectorHandler *UserInspectorHandler,
	tokenManager *jwt.TokenManager,
) *Router {
	return &Router{
		engine:                       gin.Default(),
		userHandler:                  userHandler,
		walletHandler:                walletHandler,
		transactionHandler:           transactionHandler,
		healthHandler:                healthHandler,
		adminHandler:                 adminHandler,
		dashboardHandler:             dashboardHandler,
		ledgerHandler:                ledgerHandler,
		transactionMonitoringHandler: transactionMonitoringHandler,
		refundHandler:                refundHandler,
		userInspectorHandler:         userInspectorHandler,
		tokenManager:                 tokenManager,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Global middleware
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.LoggerMiddleware())

	// ============================================
	// USER API v1
	// ============================================
	v1 := r.engine.Group("/api/v1")
	{
		// Health check (public)
		v1.GET("/health", r.healthHandler.Health)

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userHandler.Register)
			auth.POST("/login", r.userHandler.Login)
		}

		// Protected user routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(r.tokenManager))
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", r.userHandler.GetProfile)
				user.POST("/pin", r.userHandler.SetPIN)
				user.POST("/pin/verify", r.userHandler.VerifyPIN)
			}

			// Wallet routes
			wallet := protected.Group("/wallet")
			{
				wallet.GET("/balance", r.walletHandler.GetBalance)
				wallet.GET("/all", r.walletHandler.GetAllWallets)
				wallet.GET("/:wallet_id/history", r.walletHandler.GetHistory)
			}

			// Transaction routes
			transaction := protected.Group("/transaction")
			{
				transaction.POST("/topup", r.transactionHandler.Topup)
				transaction.POST("/transfer", r.transactionHandler.Transfer)
				transaction.GET("/:id", r.transactionHandler.GetTransaction)
				transaction.GET("/history", r.transactionHandler.GetUserTransactions)
			}
		}
	}

	// ============================================
	// ADMIN API v1
	// ============================================
	admin := r.engine.Group("/api/v1/admin")
	{
		// Admin auth (public)
		adminAuth := admin.Group("/auth")
		{
			adminAuth.POST("/login", r.adminHandler.Login)
		}

		// Admin protected routes
		adminProtected := admin.Group("")
		adminProtected.Use(middleware.AdminAuthMiddleware(r.tokenManager))
		{
			// ============================================
			// Dashboard (all admins)
			// ============================================
			dashboard := adminProtected.Group("/dashboard")
			{
				dashboard.GET("/overview", r.dashboardHandler.GetOverview)
				dashboard.GET("/daily-stats", r.dashboardHandler.GetDailyStats)
				dashboard.GET("/transaction-summary", r.dashboardHandler.GetTransactionSummary)
			}

			// ============================================
			// Ledger Viewer (all admins - READ ONLY)
			// ============================================
			ledger := adminProtected.Group("/ledger")
			{
				ledger.GET("", r.ledgerHandler.GetLedgerEntries)
				ledger.GET("/transaction/:id", r.ledgerHandler.GetLedgerByTransaction)
				ledger.GET("/wallet/:id", r.ledgerHandler.GetLedgerByWallet)
				ledger.GET("/wallet/:id/validate", r.ledgerHandler.ValidateBalance)
			}

			// ============================================
			// Transaction Monitoring (all admins)
			// ============================================
			transactions := adminProtected.Group("/transactions")
			{
				transactions.GET("", r.transactionMonitoringHandler.GetAllTransactions)
				transactions.GET("/pending", r.transactionMonitoringHandler.GetPendingTransactions)
				transactions.GET("/failed", r.transactionMonitoringHandler.GetFailedTransactions)
				transactions.GET("/:id", r.transactionMonitoringHandler.GetTransactionDetail)
			}

			// ============================================
			// User Inspector (all admins)
			// ============================================
			users := adminProtected.Group("/users")
			{
				users.GET("/search", r.userInspectorHandler.SearchUsers)
				users.GET("/:id", r.userInspectorHandler.GetUserDetails)
			}

			// ============================================
			// Wallet Management (ops admin + super admin)
			// ============================================
			wallets := adminProtected.Group("/wallets")
			wallets.Use(middleware.RequireOpsAdmin())
			{
				wallets.POST("/:id/freeze", r.userInspectorHandler.FreezeWallet)
				wallets.POST("/:id/unfreeze", r.userInspectorHandler.UnfreezeWallet)
			}

			// ============================================
			// Refund & Reversal (finance admin + super admin)
			// ============================================
			refund := adminProtected.Group("/refund")
			refund.Use(middleware.RequireFinanceAdmin())
			{
				refund.POST("", r.refundHandler.RefundTransaction)
				refund.POST("/reverse", r.refundHandler.ReverseTransaction)
				refund.GET("/history/:id", r.refundHandler.GetRefundHistory)
			}

			// ============================================
			// Admin Management (super admin only)
			// ============================================
			admins := adminProtected.Group("/admins")
			admins.Use(middleware.RequireSuperAdmin())
			{
				admins.POST("", r.adminHandler.CreateAdmin)
				admins.GET("", r.adminHandler.ListAdmins)
				admins.GET("/:id", r.adminHandler.GetAdmin)
				admins.PATCH("/:id/status", r.adminHandler.UpdateAdminStatus)
			}

			// ============================================
			// Audit Logs (all admins)
			// ============================================
			adminProtected.GET("/audit-logs", r.adminHandler.GetAuditLogs)
		}
	}

	return r.engine
}

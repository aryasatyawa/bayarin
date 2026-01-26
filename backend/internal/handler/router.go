package handler

import (
	"github.com/aryasatyawa/bayarin/internal/middleware"
	"github.com/aryasatyawa/bayarin/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine             *gin.Engine
	userHandler        *UserHandler
	walletHandler      *WalletHandler
	transactionHandler *TransactionHandler
	healthHandler      *HealthHandler
	tokenManager       *jwt.TokenManager
}

func NewRouter(
	userHandler *UserHandler,
	walletHandler *WalletHandler,
	transactionHandler *TransactionHandler,
	healthHandler *HealthHandler,
	tokenManager *jwt.TokenManager,
) *Router {
	return &Router{
		engine:             gin.Default(),
		userHandler:        userHandler,
		walletHandler:      walletHandler,
		transactionHandler: transactionHandler,
		healthHandler:      healthHandler,
		tokenManager:       tokenManager,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Global middleware
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.LoggerMiddleware())

	// API v1
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

		// Protected routes
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

	return r.engine
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aryasatyawa/bayarin/internal/config"
	"github.com/aryasatyawa/bayarin/internal/handler"
	"github.com/aryasatyawa/bayarin/internal/pkg/database"
	"github.com/aryasatyawa/bayarin/internal/pkg/jwt"
	"github.com/aryasatyawa/bayarin/internal/pkg/logger"
	"github.com/aryasatyawa/bayarin/internal/pkg/redis"
	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/aryasatyawa/bayarin/internal/usecase"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	logger.Init(cfg.Server.Env)
	log.Info().Msg("🚀 Starting Bayarin API Server")
	log.Info().Str("version", cfg.App.Version).Str("env", cfg.Server.Env).Msg("Configuration loaded")

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		}
	}()
	log.Info().Msg("✅ Database connected")

	// Initialize Redis
	redisClient, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing Redis connection")
		}
	}()
	log.Info().Msg("✅ Redis connected")

	// Initialize JWT token manager
	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)
	log.Info().Msg("✅ JWT token manager initialized")

	// ============================================
	// User Repositories
	// ============================================
	userRepo := repository.NewUserRepository(db.DB)
	walletRepo := repository.NewWalletRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)
	ledgerRepo := repository.NewLedgerRepository(db.DB)
	log.Info().Msg("✅ User repositories initialized")

	// ============================================
	// Admin Repositories
	// ============================================
	adminRepo := repository.NewAdminRepository(db.DB)
	auditLogRepo := repository.NewAuditLogRepository(db.DB)
	log.Info().Msg("✅ Admin repositories initialized")

	// ============================================
	// User Usecases
	// ============================================
	userUsecase := usecase.NewUserUsecase(
		db.DB,
		userRepo,
		walletRepo,
		tokenManager,
		cfg,
	)
	walletUsecase := usecase.NewWalletUsecase(
		walletRepo,
		ledgerRepo,
	)
	transactionUsecase := usecase.NewTransactionUsecase(
		db.DB,
		userRepo,
		walletRepo,
		transactionRepo,
		ledgerRepo,
		cfg,
	)
	log.Info().Msg("✅ User usecases initialized")

	// ============================================
	// Admin Usecases
	// ============================================
	adminUsecase := usecase.NewAdminUsecase(
		db.DB,
		adminRepo,
		auditLogRepo,
		tokenManager,
		cfg,
	)
	dashboardUsecase := usecase.NewDashboardUsecase(
		db.DB,
		userRepo,
		walletRepo,
		transactionRepo,
	)
	ledgerViewerUsecase := usecase.NewLedgerViewerUsecase(
		db.DB,
		ledgerRepo,
		walletRepo,
		transactionRepo,
	)
	transactionMonitoringUsecase := usecase.NewTransactionMonitoringUsecase(
		db.DB,
		transactionRepo,
		ledgerRepo,
		userRepo,
	)
	refundUsecase := usecase.NewRefundUsecase(
		db.DB,
		transactionRepo,
		walletRepo,
		ledgerRepo,
		auditLogRepo,
	)
	userInspectorUsecase := usecase.NewUserInspectorUsecase(
		db.DB,
		userRepo,
		walletRepo,
		transactionRepo,
		auditLogRepo,
	)
	log.Info().Msg("✅ Admin usecases initialized")

	// ============================================
	// User Handlers
	// ============================================
	userHandler := handler.NewUserHandler(userUsecase)
	walletHandler := handler.NewWalletHandler(walletUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)
	healthHandler := handler.NewHealthHandler(db, redisClient)
	log.Info().Msg("✅ User handlers initialized")

	// ============================================
	// Admin Handlers
	// ============================================
	adminHandler := handler.NewAdminHandler(adminUsecase)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)
	ledgerHandler := handler.NewLedgerHandler(ledgerViewerUsecase)
	transactionMonitoringHandler := handler.NewTransactionMonitoringHandler(transactionMonitoringUsecase)
	refundHandler := handler.NewRefundHandler(refundUsecase)
	userInspectorHandler := handler.NewUserInspectorHandler(userInspectorUsecase)
	log.Info().Msg("✅ Admin handlers initialized")

	// ============================================
	// Setup Router
	// ============================================
	router := handler.NewRouter(
		userHandler,
		walletHandler,
		transactionHandler,
		healthHandler,
		adminHandler,
		dashboardHandler,
		ledgerHandler,
		transactionMonitoringHandler,
		refundHandler,
		userInspectorHandler,
		tokenManager,
	)
	engine := router.Setup()
	log.Info().Msg("✅ Router configured")

	// ============================================
	// Setup HTTP Server
	// ============================================
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:           serverAddr,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("address", serverAddr).Msg("🌐 Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	log.Info().Msg("✅ Server started successfully")
	log.Info().Msgf("📍 User API: http://%s/api/v1", serverAddr)
	log.Info().Msgf("🔐 Admin API: http://%s/api/v1/admin", serverAddr)
	log.Info().Msgf("💚 Health check: http://%s/api/v1/health", serverAddr)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("🛑 Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("✅ Server stopped gracefully")
}

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
	// =====================
	// Load configuration
	// =====================
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// =====================
	// Initialize logger
	// =====================
	logger.Init(cfg.Server.Env)
	log.Info().Msg("🚀 Starting Bayarin API Server")
	log.Info().
		Str("version", cfg.App.Version).
		Str("env", cfg.Server.Env).
		Msg("Configuration loaded")

	// =====================
	// Database
	// =====================
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

	// =====================
	// Redis
	// =====================
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

	// =====================
	// JWT
	// =====================
	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)
	log.Info().Msg("✅ JWT token manager initialized")

	// =====================
	// Repositories (USER)
	// =====================
	userRepo := repository.NewUserRepository(db.DB)
	walletRepo := repository.NewWalletRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)
	ledgerRepo := repository.NewLedgerRepository(db.DB)
	log.Info().Msg("✅ User repositories initialized")

	// =====================
	// Repositories (ADMIN) - NEW
	// =====================
	adminRepo := repository.NewAdminRepository(db.DB)
	auditLogRepo := repository.NewAuditLogRepository(db.DB)
	log.Info().Msg("✅ Admin repositories initialized")

	// =====================
	// Usecases (USER)
	// =====================
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

	// =====================
	// Usecases (ADMIN) - NEW
	// =====================
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

	log.Info().Msg("✅ Admin usecases initialized")

	// =====================
	// Handlers (USER)
	// =====================
	userHandler := handler.NewUserHandler(userUsecase)
	walletHandler := handler.NewWalletHandler(walletUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)
	healthHandler := handler.NewHealthHandler(db, redisClient)
	log.Info().Msg("✅ User handlers initialized")

	// =====================
	// Handlers (ADMIN) - NEW
	// =====================
	adminHandler := handler.NewAdminHandler(adminUsecase)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)
	log.Info().Msg("✅ Admin handlers initialized")

	// =====================
	// Router (UPDATED)
	// =====================
	router := handler.NewRouter(
		userHandler,
		walletHandler,
		transactionHandler,
		healthHandler,
		adminHandler,     // NEW
		dashboardHandler, // NEW
		tokenManager,
	)
	engine := router.Setup()
	log.Info().Msg("✅ Router configured")

	// =====================
	// HTTP Server
	// =====================
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:           serverAddr,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server
	go func() {
		log.Info().Str("address", serverAddr).Msg("🌐 Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	log.Info().Msg("✅ Server started successfully")
	log.Info().Msgf("📍 API USER    : http://%s/api/v1", serverAddr)
	log.Info().Msgf("📍 API ADMIN   : http://%s/api/v1/admin", serverAddr)
	log.Info().Msgf("💚 Health     : http://%s/api/v1/health", serverAddr)

	// =====================
	// Graceful Shutdown
	// =====================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("✅ Server stopped gracefully")
}

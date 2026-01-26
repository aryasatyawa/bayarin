package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/aryasatyawa/bayarin/internal/config"
	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/pkg/crypto"
	"github.com/aryasatyawa/bayarin/internal/pkg/jwt"
	"github.com/aryasatyawa/bayarin/internal/pkg/validator"
	"github.com/aryasatyawa/bayarin/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserUsecase interface {
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	SetPIN(ctx context.Context, userID uuid.UUID, pin string) error
	VerifyPIN(ctx context.Context, userID uuid.UUID, pin string) error
}

type userUsecase struct {
	db           *sqlx.DB
	userRepo     repository.UserRepository
	walletRepo   repository.WalletRepository
	tokenManager *jwt.TokenManager
	cfg          *config.Config
}

func NewUserUsecase(
	db *sqlx.DB,
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
	tokenManager *jwt.TokenManager,
	cfg *config.Config,
) UserUsecase {
	return &userUsecase{
		db:           db,
		userRepo:     userRepo,
		walletRepo:   walletRepo,
		tokenManager: tokenManager,
		cfg:          cfg,
	}
}

// DTOs
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,indonesian_phone"`
	FullName string `json:"full_name" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Phone  string    `json:"phone"`
	Token  string    `json:"token"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // email atau phone
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Token  string    `json:"token"`
}

type UserProfile struct {
	ID        uuid.UUID         `json:"id"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	FullName  string            `json:"full_name"`
	Status    domain.UserStatus `json:"status"`
	HasPIN    bool              `json:"has_pin"`
	Wallets   []*domain.Wallet  `json:"wallets"`
	CreatedAt time.Time         `json:"created_at"`
}

// Register creates new user and main wallet
func (uc *userUsecase) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// Validate input
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Normalize phone number
	normalizedPhone := validator.NormalizePhone(req.Phone)

	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExist
	}

	existingUser, _ = uc.userRepo.GetByPhone(ctx, normalizedPhone)
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExist
	}

	// Hash password
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Begin transaction
	tx, err := uc.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create user
	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Phone:        normalizedPhone,
		FullName:     req.FullName,
		PasswordHash: passwordHash,
		Status:       domain.UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create main wallet
	wallet := &domain.Wallet{
		ID:         uuid.New(),
		UserID:     user.ID,
		WalletType: domain.WalletTypeMain,
		Balance:    0, // Start with 0 balance
		Currency:   uc.cfg.App.Currency,
		Status:     domain.WalletStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.walletRepo.CreateWithTx(ctx, tx, wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Generate JWT token
	token, err := uc.tokenManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &RegisterResponse{
		UserID: user.ID,
		Email:  user.Email,
		Phone:  user.Phone,
		Token:  token,
	}, nil
}

// Login authenticates user
func (uc *userUsecase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Validate input
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Try to find user by email or phone
	var user *domain.User
	var err error

	// Check if identifier is email
	if validator.ValidateEmail(req.Identifier) == nil {
		user, err = uc.userRepo.GetByEmail(ctx, req.Identifier)
	} else {
		// Assume it's phone number
		normalizedPhone := validator.NormalizePhone(req.Identifier)
		user, err = uc.userRepo.GetByPhone(ctx, normalizedPhone)
	}

	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Verify password
	if !crypto.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidPassword
	}

	// Check user status
	if !user.IsActive() {
		return nil, domain.ErrUserNotActive
	}

	// Generate JWT token
	token, err := uc.tokenManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		UserID: user.ID,
		Email:  user.Email,
		Token:  token,
	}, nil
}

// GetProfile returns user profile with wallets
func (uc *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get wallets
	wallets, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	return &UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		Status:    user.Status,
		HasPIN:    user.PINHash != nil,
		Wallets:   wallets,
		CreatedAt: user.CreatedAt,
	}, nil
}

// SetPIN sets user PIN for transactions
func (uc *userUsecase) SetPIN(ctx context.Context, userID uuid.UUID, pin string) error {
	// Validate PIN
	if err := validator.ValidatePIN(pin); err != nil {
		return err
	}

	// Hash PIN
	pinHash, err := crypto.HashPIN(pin)
	if err != nil {
		return fmt.Errorf("failed to hash PIN: %w", err)
	}

	// Update PIN
	if err := uc.userRepo.UpdatePIN(ctx, userID, pinHash); err != nil {
		return fmt.Errorf("failed to set PIN: %w", err)
	}

	return nil
}

// VerifyPIN verifies user PIN
func (uc *userUsecase) VerifyPIN(ctx context.Context, userID uuid.UUID, pin string) error {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if PIN is set
	if user.PINHash == nil {
		return domain.ErrInvalidPIN
	}

	// Verify PIN
	if !crypto.VerifyPIN(pin, *user.PINHash) {
		return domain.ErrInvalidPIN
	}

	return nil
}

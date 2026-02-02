package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// User Claims (existing)
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// Admin Claims (NEW)
type AdminClaims struct {
	AdminID  uuid.UUID `json:"admin_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secretKey   string
	expireHours int
}

func NewTokenManager(secretKey string, expireHours int) *TokenManager {
	return &TokenManager{
		secretKey:   secretKey,
		expireHours: expireHours,
	}
}

// GenerateToken generates user JWT token (existing)
func (tm *TokenManager) GenerateToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Hour * time.Duration(tm.expireHours))

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bayarin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateAdminToken generates admin JWT token (NEW)
func (tm *TokenManager) GenerateAdminToken(adminID uuid.UUID, username, role string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Hour * time.Duration(tm.expireHours))

	claims := AdminClaims{
		AdminID:  adminID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bayarin-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign admin token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates user JWT token (existing)
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ValidateAdminToken validates admin JWT token (NEW)
func (tm *TokenManager) ValidateAdminToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse admin token: %w", err)
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid admin token")
	}

	return claims, nil
}

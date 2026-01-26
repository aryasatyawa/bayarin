package middleware

import (
	"strings"

	"github.com/aryasatyawa/bayarin/internal/pkg/jwt"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	AuthorizationHeader = "Authorization"
	BearerSchema        = "Bearer "
	UserIDKey           = "user_id"
	UserEmailKey        = "user_email"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(tokenManager *jwt.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// Check Bearer schema
		if !strings.HasPrefix(authHeader, BearerSchema) {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, BearerSchema)
		if tokenString == "" {
			response.Unauthorized(c, "Missing token")
			c.Abort()
			return
		}

		// Validate token
		claims, err := tokenManager.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

// GetUserID gets user ID from gin context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil, response.ErrUnauthorized
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, response.ErrUnauthorized
	}

	return id, nil
}

// GetUserEmail gets user email from gin context
func GetUserEmail(c *gin.Context) (string, error) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", response.ErrUnauthorized
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", response.ErrUnauthorized
	}

	return emailStr, nil
}

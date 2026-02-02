package middleware

import (
	"strings"

	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/pkg/jwt"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	AdminIDKey       = "admin_id"
	AdminUsernameKey = "admin_username"
	AdminRoleKey     = "admin_role"
)

// AdminAuthMiddleware validates admin JWT token
func AdminAuthMiddleware(tokenManager *jwt.TokenManager) gin.HandlerFunc {
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

		// Validate admin token
		claims, err := tokenManager.ValidateAdminToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired admin token")
			c.Abort()
			return
		}

		// Set admin info in context
		c.Set(AdminIDKey, claims.AdminID)
		c.Set(AdminUsernameKey, claims.Username)
		c.Set(AdminRoleKey, claims.Role)

		c.Next()
	}
}

// GetAdminID gets admin ID from context
func GetAdminID(c *gin.Context) (uuid.UUID, error) {
	adminID, exists := c.Get(AdminIDKey)
	if !exists {
		return uuid.Nil, domain.ErrUnauthorized
	}

	id, ok := adminID.(uuid.UUID)
	if !ok {
		return uuid.Nil, domain.ErrUnauthorized
	}

	return id, nil
}

// GetAdminRole gets admin role from context
func GetAdminRole(c *gin.Context) (domain.AdminRole, error) {
	role, exists := c.Get(AdminRoleKey)
	if !exists {
		return "", domain.ErrUnauthorized
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", domain.ErrUnauthorized
	}

	return domain.AdminRole(roleStr), nil
}

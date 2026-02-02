package middleware

import (
	"github.com/aryasatyawa/bayarin/internal/domain"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// RequireRole middleware checks if admin has required role
func RequireRole(allowedRoles ...domain.AdminRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminRole, err := GetAdminRole(c)
		if err != nil {
			response.Unauthorized(c, "Admin not authenticated")
			c.Abort()
			return
		}

		// Super admin bypasses all role checks
		if adminRole == domain.RoleSuperAdmin {
			c.Next()
			return
		}

		// Check if admin role is in allowed roles
		hasPermission := false
		for _, role := range allowedRoles {
			if adminRole == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin middleware requires super admin role
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole(domain.RoleSuperAdmin)
}

// RequireFinanceAdmin middleware requires finance admin role
func RequireFinanceAdmin() gin.HandlerFunc {
	return RequireRole(domain.RoleSuperAdmin, domain.RoleFinanceAdmin)
}

// RequireOpsAdmin middleware requires ops admin role
func RequireOpsAdmin() gin.HandlerFunc {
	return RequireRole(domain.RoleSuperAdmin, domain.RoleOpsAdmin)
}

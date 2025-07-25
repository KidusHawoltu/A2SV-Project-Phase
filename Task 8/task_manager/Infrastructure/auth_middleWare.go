package infrastructure

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService domain.JwtService
}

func NewAuthMiddleware(jwtService domain.JwtService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// Authenticate is the primary authentication middleware.
// It verifies the token and stores *all* claims in the context.
// It does NOT perform any authorization checks itself.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("AuthMiddleware: Missing or malformed Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization token required"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.jwtService.ParseToken(c.Request.Context(), tokenString)
		if err != nil {
			log.Printf("AuthMiddleware: Token parsing/validation failed: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		// Store relevant claims in context using string literals
		c.Set("userID", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("userRole", claims.Role)

		c.Next() // Proceed to the next handler
	}
}

// AuthorizeAdmin is an authorization middleware.
// It ASSUMES Authenticate() has already run and set claims in context.
func (m *AuthMiddleware) AuthorizeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get claims from context using string literals
		role, exists := c.Get("userRole")
		if !exists {
			log.Println("AuthorizeAdmin: User role not found in context. Authenticate middleware likely missing or failed.")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Authentication context missing or invalid"})
			return
		}

		userRole, ok := role.(domain.UserRole) // Type assert to domain.UserRole
		if !ok {
			log.Printf("AuthorizeAdmin: Invalid user role type in context: %T\n", role)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Invalid user role format"})
			return
		}

		if userRole != domain.RoleAdmin {
			log.Printf("AuthorizeAdmin: User '%s' (ID: %s) attempted to access admin route without Admin role (Role: %s)\n", c.GetString("username"), c.GetString("userID"), userRole)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Access forbidden: Admin role required"})
			return
		}

		c.Next() // User is Admin, proceed to the controller
	}
}

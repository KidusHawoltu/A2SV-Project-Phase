package middleware

import (
	"A2SV_ProjectPhase/Task6/TaskManager/models"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Claims struct to hold the JWT payload
type Claims struct {
	UserID   string          `json:"user_id"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.StandardClaims
}

// AuthMiddleware creates a Gin middleware for JWT authentication.
// It validates the JWT from the Authorization header and sets user info in context.
func AuthMiddleware() gin.HandlerFunc {
	// Retrieve JWT_SECRET from environment variables once during middleware initialization
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Fallback for local development, but in production, this should ideally be fatal
		// as a missing secret is a critical configuration error.
		jwtSecret = "default secret" // Ensure this matches your main.go fallback for local dev
		fmt.Println("WARNING: JWT_SECRET environment variable not set. Using default secret.")
	}
	secretBytes := []byte(jwtSecret) // Convert string secret to byte slice

	return func(c *gin.Context) {
		// 1. Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required", "error": "missing_token"})
			return
		}

		// 2. Check for "Bearer " prefix and split the token string
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format", "error": "invalid_header_format"})
			return
		}

		tokenString := parts[1]

		// 3. Parse and validate the JWT token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretBytes, nil // Return the secret key as a byte slice
		})

		if err != nil {
			// Handle different types of JWT errors for specific HTTP responses
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is malformed", "error": "malformed_token"})
					return
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is expired or not yet valid", "error": "expired_or_invalid_token"})
					return
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Could not process token", "error": err.Error()})
					return
				}
			}
			// Generic invalid token error
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": err.Error()})
			return
		}

		// 4. Check if token is valid and claims can be extracted
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": "token_not_valid"})
			return
		}

		// 5. Store user information from claims in Gin's context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		// Proceed to the next handler in the chain
		c.Next()
	}
}

// AuthorizeRole creates a Gin middleware for role-based authorization.
// It checks if the authenticated user's role matches any of the required roles.
func AuthorizeRole(requiredRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			// This indicates AuthMiddleware wasn't run or failed to set the role
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Authentication context missing role", "error": "authorization_internal_error"})
			return
		}

		currentUserRole, ok := userRole.(models.UserRole)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Invalid user role type in context", "error": "authorization_type_error"})
			return
		}

		isAuthorized := slices.Contains(requiredRoles, currentUserRole)

		if !isAuthorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden: Insufficient role permissions", "error": "access_denied"})
			return
		}

		c.Next()
	}
}

// GetUserIDFromContext is a helper to safely extract the user's ObjectID from Gin context.
// Should be used in protected handlers after AuthMiddleware.
func GetUserIDFromContext(c *gin.Context) (primitive.ObjectID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		// This indicates a middleware configuration issue or handler called directly without middleware
		return primitive.NilObjectID, errors.New("user ID not found in context (authentication middleware not applied correctly)")
	}
	idStr, ok := userID.(string)
	if !ok {
		return primitive.NilObjectID, errors.New("user ID in context is not of type string")
	}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.NilObjectID, errors.New("user ID in context is not of type primitive.ObjectID")
	}
	return id, nil
}

// GetUserRoleFromContext is a helper to safely extract the user's Role from Gin context.
// Should be used in protected handlers after AuthMiddleware.
func GetUserRoleFromContext(c *gin.Context) (models.UserRole, error) {
	userRole, exists := c.Get("user_role")
	if !exists {
		// This indicates a middleware configuration issue or handler called directly without middleware
		return "", errors.New("user role not found in context (authentication middleware not applied correctly)")
	}
	role, ok := userRole.(models.UserRole)
	if !ok {
		return "", errors.New("user role in context is not of type models.UserRole")
	}
	return role, nil
}

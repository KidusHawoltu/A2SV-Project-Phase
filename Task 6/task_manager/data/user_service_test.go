package data

import (
	"A2SV_ProjectPhase/Task6/TaskManager/middleware"
	"A2SV_ProjectPhase/Task6/TaskManager/models"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// setupUserTestCollection cleans the test user collection before each user test.
func setupUserTestCollection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := testUserCollection.Drop(ctx) // Use testUserCollection
	if err != nil && err.Error() != "ns not found" {
		t.Fatalf("Failed to drop test user collection: %v", err)
	}
	t.Logf("Cleaned test user collection '%s' for new test.", userCollectionName) // Use t.Logf for test-specific logs
}

// addSampleUser is a helper to register a user for testing purposes.
func addSampleUser(ctx context.Context, t *testing.T, username, password string, role models.UserRole) models.User {
	user := models.User{
		Username: username,
		Password: password,
		Role:     role,
	}
	createdUser, err := testUserManager.RegisterUser(ctx, user) // Use testUserManager
	assert.NoError(t, err, "Failed to register sample user")
	assert.False(t, createdUser.Id.IsZero(), "Registered user ID should not be zero")
	return createdUser
}

func TestRegisterUser(t *testing.T) {
	setupUserTestCollection(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("should register a new user successfully", func(t *testing.T) {
		userToRegister := models.User{
			Username: "testuser1",
			Password: "password123",
			Role:     models.RoleUser,
		}

		createdUser, err := testUserManager.RegisterUser(ctx, userToRegister)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.False(t, createdUser.Id.IsZero(), "User ID should be generated")
		assert.Equal(t, "testuser1", createdUser.Username)
		assert.Equal(t, models.RoleUser, createdUser.Role)
		assert.Empty(t, createdUser.Password, "Password hash should not be returned in the user model")

		// Verify user exists in DB by trying to retrieve it
		retrievedUser, err := testUserManager.(*UserCollection).GetByUsername(ctx, "testuser1") // Access private helper for verification
		assert.NoError(t, err)
		assert.Equal(t, createdUser.Id, retrievedUser.Id)
		assert.Equal(t, "testuser1", retrievedUser.Username)
		assert.Equal(t, models.RoleUser, retrievedUser.Role)
		assert.NotEmpty(t, retrievedUser.Password, "Hashed password should be stored in DB")
	})

	t.Run("should return ErrUsernameTaken if username already exists", func(t *testing.T) {
		// Arrange: Register the first user
		addSampleUser(ctx, t, "existinguser", "pass", models.RoleUser)

		// Act: Try to register with the same username
		duplicateUser := models.User{
			Username: "existinguser",
			Password: "newpassword",
			Role:     models.RoleUser,
		}
		_, err := testUserManager.RegisterUser(ctx, duplicateUser)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrUsernameTaken), "Error should be ErrUsernameTaken")
	})
}

func TestLogin(t *testing.T) {
	setupUserTestCollection(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Arrange: Register a user for login tests
	registeredUser := addSampleUser(ctx, t, "loginuser", "securepass", models.RoleUser)

	t.Run("should successfully log in with correct credentials", func(t *testing.T) {
		token, err := testUserManager.Login(ctx, "loginuser", "securepass")

		assert.NoError(t, err)
		assert.NotEmpty(t, token, "Should return a non-empty JWT token")

		// Optional: Verify JWT claims (this implicitly tests signing method and claims)
		claims := &middleware.Claims{} // Use the same Claims struct from middleware for validation
		parsedToken, parseErr := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("test_jwt_secret"), nil // Use the same secret as in TestMain
		})
		assert.NoError(t, parseErr, "Failed to parse generated JWT token")
		assert.True(t, parsedToken.Valid, "Generated JWT token should be valid")
		assert.Equal(t, registeredUser.Id.Hex(), claims.UserID, "Token UserID should match registered user's ID hex string")
		assert.Equal(t, registeredUser.Username, claims.Username, "Token Username should match registered user's username")
		assert.Equal(t, registeredUser.Role, claims.Role, "Token Role should match registered user's role")
	})

	t.Run("should return ErrInvalidCredentials for non-existent username", func(t *testing.T) {
		_, err := testUserManager.Login(ctx, "nonexistent", "anypassword")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidCredentials), "Error should be ErrInvalidCredentials")
	})

	t.Run("should return ErrInvalidCredentials for incorrect password", func(t *testing.T) {
		_, err := testUserManager.Login(ctx, "loginuser", "wrongpassword") // Use a correct username but wrong password

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidCredentials), "Error should be ErrInvalidCredentials")
	})
}

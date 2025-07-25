package infrastructure_test

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	"A2SV_ProjectPhase/Task8/TaskManager/Infrastructure"
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//===========================================================================
// JwtService Test Suite
//===========================================================================

type JwtServiceSuite struct {
	suite.Suite
	jwtService         *infrastructure.MyJwtService
	secretKey          string
	differentSecretKey string
}

// TestJwtServiceSuite is the entry point for the test suite
func TestJwtServiceSuite(t *testing.T) {
	suite.Run(t, new(JwtServiceSuite))
}

// SetupTest runs before each test method
func (s *JwtServiceSuite) SetupTest() {
	s.secretKey = "a-very-secure-secret-for-testing"
	s.differentSecretKey = "a-completely-different-secret"
	s.jwtService = infrastructure.NewJwtService(s.secretKey)
}

// TestGetAndParseToken_Success tests the full round-trip of creating and parsing a valid token.
func (s *JwtServiceSuite) TestGetAndParseToken_Success() {
	// --- Setup ---
	user := &domain.User{
		Id:       primitive.NewObjectID(),
		Username: "testuser",
		Role:     domain.RoleUser,
	}

	// --- Execution: Generate Token ---
	tokenString, err := s.jwtService.GetSignedToken(context.Background(), user)

	// --- Assertion: Generation ---
	s.Require().NoError(err, "Token generation should not fail")
	s.NotEmpty(tokenString, "Generated token string should not be empty")

	// --- Execution: Parse Token ---
	claims, err := s.jwtService.ParseToken(context.Background(), tokenString)

	// --- Assertion: Parsing and Claims ---
	s.Require().NoError(err, "Parsing a valid token should not fail")
	s.Require().NotNil(claims)

	s.Equal(user.Id.Hex(), claims.UserId, "User ID in claims should match original")
	s.Equal(user.Username, claims.Username, "Username in claims should match original")
	s.Equal(user.Role, claims.Role, "Role in claims should match original")
	s.Greater(claims.ExpiresAt, time.Now().Unix(), "Token should expire in the future")
}

// TestParseToken_Failure tests various invalid token scenarios.
func (s *JwtServiceSuite) TestParseToken_Failure() {
	s.Run("Invalid Signature", func() {
		// --- Setup ---
		// Create a token with a DIFFERENT secret key
		otherService := infrastructure.NewJwtService(s.differentSecretKey)
		user := &domain.User{Id: primitive.NewObjectID(), Username: "user"}
		tokenString, err := otherService.GetSignedToken(context.Background(), user)
		s.Require().NoError(err)

		// --- Execution & Assertion ---
		// Try to parse it with our main service (which has the original key)
		_, err = s.jwtService.ParseToken(context.Background(), tokenString)
		s.Require().Error(err, "Should fail with an invalid signature")
		s.Equal("invalid token", err.Error(), "Error message should be 'invalid token'")
	})

	s.Run("Expired Token", func() {
		// --- Setup ---
		// Manually create an expired token
		expiredClaims := domain.Claims{
			UserId:   primitive.NewObjectID().Hex(),
			Username: "expiredUser",
			StandardClaims: jwt.StandardClaims{
				// Set expiry to one hour in the past
				ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		expiredTokenString, err := token.SignedString([]byte(s.secretKey))
		s.Require().NoError(err)

		// --- Execution & Assertion ---
		_, err = s.jwtService.ParseToken(context.Background(), expiredTokenString)
		s.Require().Error(err, "Should fail with an expired token error")
		s.Equal("token expired", err.Error(), "Error message should be 'token expired'")
	})

	s.Run("Malformed Token", func() {
		// --- Execution & Assertion ---
		_, err := s.jwtService.ParseToken(context.Background(), "this.is.not.a.valid.jwt")
		s.Require().Error(err, "Should fail with a malformed token")
		s.Equal("invalid token", err.Error(), "Error message should be 'invalid token'")
	})

	s.Run("Invalid Signing Method", func() {
		// --- Setup ---
		// A manually crafted token where the header's `alg` field is set to `ES256`,
		// which is not the `HS256` our service expects.
		// Header: {"alg":"ES256","typ":"JWT"}
		// Payload: {}
		// Signature: (empty)
		badAlgToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.e30.A" // Header alg is ES256

		// --- Execution & Assertion ---
		_, err := s.jwtService.ParseToken(context.Background(), badAlgToken)
		s.Require().Error(err)
		// The underlying error is "unexpected signing method", which we map to "invalid token"
		s.Equal("invalid token", err.Error())
	})
}

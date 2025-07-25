package infrastructure_test

import (
	"A2SV_ProjectPhase/Task8/TaskManager/Domain"
	"A2SV_ProjectPhase/Task8/TaskManager/Infrastructure"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mock JwtService for testing ---
type MockJwtService struct {
	ParseTokenFunc func(c context.Context, token string) (*domain.Claims, error)
}

func (m *MockJwtService) ParseToken(c context.Context, token string) (*domain.Claims, error) {
	if m.ParseTokenFunc != nil {
		return m.ParseTokenFunc(c, token)
	}
	return nil, errors.New("ParseTokenFunc not implemented")
}
func (m *MockJwtService) GetSignedToken(c context.Context, user *domain.User) (string, error) {
	return "", errors.New("GetSignedToken not needed for this test")
}

//===========================================================================
// AuthMiddleware Test Suite
//===========================================================================

type AuthMiddlewareSuite struct {
	suite.Suite
	mockJwtService *MockJwtService
	middleware     *infrastructure.AuthMiddleware
}

// TestAuthMiddlewareSuite is the entry point for the test suite
func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareSuite))
}

// SetupTest runs before each test method
func (s *AuthMiddlewareSuite) SetupTest() {
	// Set Gin to test mode to silence unnecessary logs
	gin.SetMode(gin.TestMode)

	s.mockJwtService = &MockJwtService{}
	s.middleware = infrastructure.NewAuthMiddleware(s.mockJwtService)
}

// Helper function to create a new router, serve a request, and return the recorder
func (s *AuthMiddlewareSuite) serveRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

// --- Tests for Authenticate() middleware ---

func (s *AuthMiddlewareSuite) TestAuthenticate() {
	// Create a new router for this test method
	router := gin.New()
	dummyHandler := func(c *gin.Context) { c.Status(http.StatusOK) }
	router.GET("/test-auth", s.middleware.Authenticate(), dummyHandler)

	s.Run("Success", func() {
		userID := primitive.NewObjectID()
		expectedClaims := &domain.Claims{
			UserId:   userID.Hex(),
			Username: "testuser",
			Role:     domain.RoleUser,
		}
		s.mockJwtService.ParseTokenFunc = func(c context.Context, token string) (*domain.Claims, error) {
			s.Equal("valid-token", token, "The correct token string should be passed to the service")
			return expectedClaims, nil
		}
		req, _ := http.NewRequest(http.MethodGet, "/test-auth", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		recorder := s.serveRequest(router, req)
		s.Equal(http.StatusOK, recorder.Code, "Request should be allowed to proceed")
	})

	s.Run("Failure - No Authorization Header", func() {
		req, _ := http.NewRequest(http.MethodGet, "/test-auth", nil)
		recorder := s.serveRequest(router, req)
		s.Equal(http.StatusUnauthorized, recorder.Code)
		s.Contains(recorder.Body.String(), "Authorization token required")
	})

	s.Run("Failure - Malformed Header (No Bearer prefix)", func() {
		req, _ := http.NewRequest(http.MethodGet, "/test-auth", nil)
		req.Header.Set("Authorization", "invalid-token")
		recorder := s.serveRequest(router, req)
		s.Equal(http.StatusUnauthorized, recorder.Code)
	})

	s.Run("Failure - Token Parsing Error", func() {
		s.mockJwtService.ParseTokenFunc = func(c context.Context, token string) (*domain.Claims, error) {
			return nil, errors.New("invalid signature")
		}
		req, _ := http.NewRequest(http.MethodGet, "/test-auth", nil)
		req.Header.Set("Authorization", "Bearer bad-token")
		recorder := s.serveRequest(router, req)
		s.Equal(http.StatusUnauthorized, recorder.Code)
		s.Contains(recorder.Body.String(), "invalid signature")
	})
}

// --- Tests for AuthorizeAdmin() middleware ---

func (s *AuthMiddlewareSuite) TestAuthorizeAdmin() {
	dummyHandler := func(c *gin.Context) { c.Status(http.StatusOK) }

	// Helper function creates a setup middleware with the desired context values.
	createContextSetter := func(userID, username string, role domain.UserRole) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("userID", userID)
			c.Set("username", username)
			c.Set("userRole", role)
			c.Next()
		}
	}

	s.Run("Success - User is Admin", func() {
		// Setup a new router with the correct middleware chain for this specific test case.
		router := gin.New()
		adminSetupMiddleware := createContextSetter(primitive.NewObjectID().Hex(), "admin", domain.RoleAdmin)
		router.GET("/admin-only", adminSetupMiddleware, s.middleware.AuthorizeAdmin(), dummyHandler)

		req, _ := http.NewRequest(http.MethodGet, "/admin-only", nil)
		recorder := s.serveRequest(router, req)
		s.Equal(http.StatusOK, recorder.Code, "Admin user should be allowed access")
	})

	s.Run("Failure - User is not Admin", func() {
		router := gin.New()
		userSetupMiddleware := createContextSetter(primitive.NewObjectID().Hex(), "user", domain.RoleUser)
		router.GET("/admin-only", userSetupMiddleware, s.middleware.AuthorizeAdmin(), dummyHandler)

		req, _ := http.NewRequest(http.MethodGet, "/admin-only", nil)
		recorder := s.serveRequest(router, req)
		s.Equal(http.StatusForbidden, recorder.Code)
		s.Contains(recorder.Body.String(), "Access forbidden")
	})

	s.Run("Failure - Role not in context", func() {
		// For this case, we explicitly do NOT add our context-setting middleware.
		router := gin.New()
		router.GET("/admin-no-context", s.middleware.AuthorizeAdmin(), dummyHandler)

		req, _ := http.NewRequest(http.MethodGet, "/admin-no-context", nil)
		recorder := s.serveRequest(router, req)
		s.Equal(http.StatusInternalServerError, recorder.Code)
		s.Contains(recorder.Body.String(), "Authentication context missing")
	})
}

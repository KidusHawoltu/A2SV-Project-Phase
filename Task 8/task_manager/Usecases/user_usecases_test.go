package usecases_test

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	usecases "A2SV_ProjectPhase/Task8/TaskManager/Usecases"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mocks (can be kept as is, they are well-defined) ---
type MockUserRepository struct {
	GetUserByUsernameFunc func(c context.Context, username string) (*domain.User, error)
	CreateUserFunc        func(c context.Context, user *domain.User) (*domain.User, error)
}

func (m *MockUserRepository) GetUserByUsername(c context.Context, username string) (*domain.User, error) {
	return m.GetUserByUsernameFunc(c, username)
}
func (m *MockUserRepository) CreateUser(c context.Context, user *domain.User) (*domain.User, error) {
	return m.CreateUserFunc(c, user)
}

type MockPasswordService struct {
	HashFunc    func(c context.Context, password string) (string, error)
	CompareFunc func(c context.Context, password string, hashedPassword string) error
}

func (m *MockPasswordService) Hash(c context.Context, password string) (string, error) {
	return m.HashFunc(c, password)
}
func (m *MockPasswordService) Compare(c context.Context, password string, hash string) error {
	return m.CompareFunc(c, password, hash)
}

type MockJwtService struct {
	GetSignedTokenFunc func(c context.Context, user *domain.User) (string, error)
	ParseTokenFunc     func(c context.Context, token string) (*domain.Claims, error)
}

func (m *MockJwtService) GetSignedToken(c context.Context, user *domain.User) (string, error) {
	return m.GetSignedTokenFunc(c, user)
}
func (m *MockJwtService) ParseToken(c context.Context, token string) (*domain.Claims, error) {
	return m.ParseTokenFunc(c, token)
}

//===========================================================================
// UserUseCase Test Suite
//===========================================================================

type UserUseCaseSuite struct {
	suite.Suite
	mockUserRepo    *MockUserRepository
	mockJwtService  *MockJwtService
	mockPassService *MockPasswordService
	useCase         *usecases.UserUseCase
	ctx             context.Context
}

// TestUserUseCaseSuite is the entry point for the test suite
func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseSuite))
}

// SetupTest runs before each test method. It's the perfect place for initialization.
func (s *UserUseCaseSuite) SetupTest() {
	s.mockUserRepo = &MockUserRepository{}
	s.mockJwtService = &MockJwtService{}
	s.mockPassService = &MockPasswordService{}
	s.useCase = usecases.NewUserUseCase(s.mockUserRepo, s.mockJwtService, s.mockPassService)
	s.ctx = context.Background()
}

// TestRegisterUser contains all sub-tests for the registration logic.
func (s *UserUseCaseSuite) TestRegisterUser() {
	testUsername := "testuser"
	testPassword := "password123"
	hashedPassword := "hashedpassword"

	s.Run("Success", func() {
		// --- Mock Configuration ---
		s.mockUserRepo.GetUserByUsernameFunc = func(c context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound // User is available
		}
		s.mockPassService.HashFunc = func(c context.Context, password string) (string, error) {
			s.Equal(testPassword, password)
			return hashedPassword, nil
		}
		s.mockUserRepo.CreateUserFunc = func(c context.Context, user *domain.User) (*domain.User, error) {
			user.Id = primitive.NewObjectID() // Simulate DB setting the ID
			return user, nil
		}

		// --- Execution ---
		registeredUser, err := s.useCase.RegisterUser(s.ctx, testUsername, testPassword)

		// --- Assertion ---
		s.Require().NoError(err)
		s.Require().NotNil(registeredUser)
		s.Equal(testUsername, registeredUser.Username)
		s.Equal(hashedPassword, registeredUser.PasswordHash)
		s.False(registeredUser.Id.IsZero())
	})

	s.Run("Failure - Username Taken", func() {
		s.mockUserRepo.GetUserByUsernameFunc = func(c context.Context, username string) (*domain.User, error) {
			return &domain.User{}, nil // User exists
		}

		_, err := s.useCase.RegisterUser(s.ctx, testUsername, testPassword)

		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrUsernameTaken)
	})

	s.Run("Failure - Password Hashing Error", func() {
		expectedErr := errors.New("hashing failed")
		s.mockUserRepo.GetUserByUsernameFunc = func(c context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		}
		s.mockPassService.HashFunc = func(c context.Context, password string) (string, error) {
			return "", expectedErr
		}

		_, err := s.useCase.RegisterUser(s.ctx, testUsername, testPassword)

		s.Require().Error(err)
		s.ErrorIs(err, expectedErr)
	})

	s.Run("Failure - Domain Validation Error", func() {
		s.mockPassService.HashFunc = func(c context.Context, password string) (string, error) {
			return "any-hashed-password", nil
		}
		_, err := s.useCase.RegisterUser(s.ctx, "", testPassword) // Empty username

		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrValidationFailed)
	})
}

// TestLogin contains all sub-tests for the login logic.
func (s *UserUseCaseSuite) TestLogin() {
	testUsername := "existinguser"
	testPassword := "password123"
	hashedPassword := "hashedpassword"
	expectedToken := "test.jwt.token"
	mockUser := &domain.User{Id: primitive.NewObjectID(), Username: testUsername, PasswordHash: hashedPassword}

	s.Run("Success", func() {
		// --- Mock Configuration ---
		s.mockUserRepo.GetUserByUsernameFunc = func(c context.Context, username string) (*domain.User, error) {
			return mockUser, nil
		}
		s.mockPassService.CompareFunc = func(c context.Context, password, hash string) error {
			s.Equal(testPassword, password)
			s.Equal(hashedPassword, hash)
			return nil // Passwords match
		}
		s.mockJwtService.GetSignedTokenFunc = func(c context.Context, user *domain.User) (string, error) {
			s.Equal(mockUser.Id, user.Id)
			return expectedToken, nil
		}

		// --- Execution ---
		token, err := s.useCase.Login(s.ctx, testUsername, testPassword)

		// --- Assertion ---
		s.Require().NoError(err)
		s.Equal(expectedToken, token)
	})

	s.Run("Failure - User Not Found", func() {
		s.mockUserRepo.GetUserByUsernameFunc = func(c context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		}

		_, err := s.useCase.Login(s.ctx, testUsername, testPassword)

		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrInvalidCredentials)
	})

	s.Run("Failure - Password Mismatch", func() {
		s.mockUserRepo.GetUserByUsernameFunc = func(c context.Context, username string) (*domain.User, error) {
			return mockUser, nil
		}
		s.mockPassService.CompareFunc = func(c context.Context, password, hash string) error {
			return errors.New("password comparison failed") // Simulate mismatch
		}

		_, err := s.useCase.Login(s.ctx, testUsername, testPassword)

		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrInvalidCredentials)
	})

	s.Run("Failure - Token Generation Error", func() {
		expectedErr := errors.New("token signing failed")
		s.mockUserRepo.GetUserByUsernameFunc = func(c context.Context, username string) (*domain.User, error) {
			return mockUser, nil
		}
		s.mockPassService.CompareFunc = func(c context.Context, password, hash string) error {
			return nil // Passwords match
		}
		s.mockJwtService.GetSignedTokenFunc = func(c context.Context, user *domain.User) (string, error) {
			return "", expectedErr
		}

		_, err := s.useCase.Login(s.ctx, testUsername, testPassword)

		s.Require().Error(err)
		s.ErrorIs(err, expectedErr)
	})
}

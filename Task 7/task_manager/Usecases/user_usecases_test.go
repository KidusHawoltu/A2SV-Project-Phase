package usecases_test

import (
	domain "A2SV_ProjectPhase/Task7/TaskManager/Domain"     // Import domain for interfaces and entities
	usecases "A2SV_ProjectPhase/Task7/TaskManager/Usecases" // Import usecases for the UserUseCase
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mocks for Dependencies ---

// MockUserRepository implements domain.UserRepository
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

// MockPasswordService implements domain.PasswordService
type MockPasswordService struct {
	HashFunc    func(c context.Context, password string) (string, error)
	CompareFunc func(c context.Context, password string, hashedPassword string) error
}

func (m *MockPasswordService) Hash(c context.Context, password string) (string, error) {
	return m.HashFunc(c, password)
}
func (m *MockPasswordService) Compare(c context.Context, password string, hashedPassword string) error {
	return m.CompareFunc(c, password, hashedPassword)
}

// MockJwtService implements domain.JwtService
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

// --- Helper for Test Context ---
func getTestContext() context.Context {
	return context.Background()
}

// --- Tests for RegisterUser ---

func TestRegisterUser_Success(t *testing.T) {
	ctx := getTestContext()
	testUsername := "testuser"
	testPassword := "password123"
	hashedPassword := "hashedpassword"

	mockUser := &domain.User{
		Id:           primitive.NewObjectID(),
		Username:     testUsername,
		PasswordHash: hashedPassword,
		Role:         domain.RoleUser,
	}

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound // User not found, so it can be registered
		},
		CreateUserFunc: func(c context.Context, user *domain.User) (*domain.User, error) {
			user.Id = mockUser.Id // Simulate MongoDB assigning an ID
			return user, nil
		},
	}
	mockPasswordService := &MockPasswordService{
		HashFunc: func(c context.Context, password string) (string, error) {
			return hashedPassword, nil
		},
	}
	mockJwtService := &MockJwtService{} // Not used in RegisterUser

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	registeredUser, err := userUseCase.RegisterUser(ctx, testUsername, testPassword)

	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	if registeredUser == nil {
		t.Fatal("Registered user is nil")
	}
	if registeredUser.Username != testUsername {
		t.Errorf("Username mismatch: got %s, want %s", registeredUser.Username, testUsername)
	}
	if registeredUser.PasswordHash != hashedPassword {
		t.Errorf("PasswordHash mismatch: got %s, want %s", registeredUser.PasswordHash, hashedPassword)
	}
	if registeredUser.Id.IsZero() {
		t.Errorf("User ID was not set")
	}
}

func TestRegisterUser_ErrUsernameTaken(t *testing.T) {
	ctx := getTestContext()
	testUsername := "existinguser"

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return &domain.User{Username: username}, nil // User already exists
		},
	}
	mockPasswordService := &MockPasswordService{}
	mockJwtService := &MockJwtService{}

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	_, err := userUseCase.RegisterUser(ctx, testUsername, "password123")

	if err == nil {
		t.Fatal("RegisterUser did not return error for taken username")
	}
	if !errors.Is(err, domain.ErrUsernameTaken) {
		t.Errorf("Error mismatch: got %v, want %v", err, domain.ErrUsernameTaken)
	}
}

func TestRegisterUser_PasswordHashError(t *testing.T) {
	ctx := getTestContext()
	testUsername := "testuser"
	expectedErr := errors.New("hashing failed")

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}
	mockPasswordService := &MockPasswordService{
		HashFunc: func(c context.Context, password string) (string, error) {
			return "", expectedErr // Simulate hashing error
		},
	}
	mockJwtService := &MockJwtService{}

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	_, err := userUseCase.RegisterUser(ctx, testUsername, "password123")

	if err == nil {
		t.Fatal("RegisterUser did not return error for hashing failure")
	}
	if !errors.Is(err, expectedErr) { // Check if the underlying error is preserved
		t.Errorf("Error mismatch: got %v, want error containing %v", err, expectedErr)
	}
}

func TestRegisterUser_CreateUserError(t *testing.T) {
	ctx := getTestContext()
	testUsername := "testuser"
	hashedPassword := "hashedpassword"
	expectedErr := errors.New("db create error")

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
		CreateUserFunc: func(c context.Context, user *domain.User) (*domain.User, error) {
			return nil, expectedErr // Simulate create error
		},
	}
	mockPasswordService := &MockPasswordService{
		HashFunc: func(c context.Context, password string) (string, error) {
			return hashedPassword, nil
		},
	}
	mockJwtService := &MockJwtService{}

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	_, err := userUseCase.RegisterUser(ctx, testUsername, "password123")

	if err == nil {
		t.Fatal("RegisterUser did not return error for create user failure")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Error mismatch: got %v, want error containing %v", err, expectedErr)
	}
}

func TestRegisterUser_DomainNewUserError(t *testing.T) {
	ctx := getTestContext()
	testUsername := "" // Will trigger domain.NewUser error
	testPassword := "password123"

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}
	mockPasswordService := &MockPasswordService{
		HashFunc: func(c context.Context, password string) (string, error) {
			return "hashed", nil
		},
	}
	mockJwtService := &MockJwtService{}

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	_, err := userUseCase.RegisterUser(ctx, testUsername, testPassword)

	if err == nil {
		t.Fatal("RegisterUser did not return error for domain.NewUser failure")
	}
}

// --- Tests for Login ---

func TestLogin_Success(t *testing.T) {
	ctx := getTestContext()
	testUsername := "existinguser"
	testPassword := "password123"
	hashedPassword := "hashedpassword"
	expectedToken := "test.jwt.token"

	mockUser := &domain.User{
		Id:           primitive.NewObjectID(),
		Username:     testUsername,
		PasswordHash: hashedPassword,
		Role:         domain.RoleUser,
	}

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return mockUser, nil
		},
	}
	mockPasswordService := &MockPasswordService{
		CompareFunc: func(c context.Context, password string, hashedPassword string) error {
			return nil // Passwords match
		},
	}
	mockJwtService := &MockJwtService{
		GetSignedTokenFunc: func(c context.Context, user *domain.User) (string, error) {
			return expectedToken, nil
		},
	}

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	token, err := userUseCase.Login(ctx, testUsername, testPassword)

	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token != expectedToken {
		t.Errorf("Token mismatch: got %s, want %s", token, expectedToken)
	}
}

func TestLogin_ErrInvalidCredentials_UserNotFound(t *testing.T) {
	ctx := getTestContext()
	testUsername := "nonexistent"

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound // User not found
		},
	}
	mockPasswordService := &MockPasswordService{}
	mockJwtService := &MockJwtService{}

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	_, err := userUseCase.Login(ctx, testUsername, "password123")

	if err == nil {
		t.Fatal("Login did not return error for user not found")
	}
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Error mismatch: got %v, want %v", err, domain.ErrInvalidCredentials)
	}
}

func TestLogin_ErrInvalidCredentials_PasswordMismatch(t *testing.T) {
	ctx := getTestContext()
	testUsername := "existinguser"
	testPassword := "wrongpassword"
	hashedPassword := "hashedpassword"                                                                                               // Correct hash
	compareErr := errors.New("password service: password comparison failed: crypto/bcrypt: hashedPassword and password don't match") // bcrypt mismatch error

	mockUser := &domain.User{
		Id:           primitive.NewObjectID(),
		Username:     testUsername,
		PasswordHash: hashedPassword,
		Role:         domain.RoleUser,
	}

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return mockUser, nil
		},
	}
	mockPasswordService := &MockPasswordService{
		CompareFunc: func(c context.Context, password string, hashedPassword string) error {
			return compareErr // Simulate password mismatch
		},
	}
	mockJwtService := &MockJwtService{}

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	_, err := userUseCase.Login(ctx, testUsername, testPassword)

	if err == nil {
		t.Fatal("Login did not return error for password mismatch")
	}
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Error mismatch: got %v, want %v", err, domain.ErrInvalidCredentials)
	}
}

func TestLogin_GetSignedTokenError(t *testing.T) {
	ctx := getTestContext()
	testUsername := "existinguser"
	testPassword := "password123"
	hashedPassword := "hashedpassword"
	expectedErr := errors.New("token signing failed")

	mockUser := &domain.User{
		Id:           primitive.NewObjectID(),
		Username:     testUsername,
		PasswordHash: hashedPassword,
		Role:         domain.RoleUser,
	}

	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(c context.Context, username string) (*domain.User, error) {
			return mockUser, nil
		},
	}
	mockPasswordService := &MockPasswordService{
		CompareFunc: func(c context.Context, password string, hashedPassword string) error {
			return nil
		},
	}
	mockJwtService := &MockJwtService{
		GetSignedTokenFunc: func(c context.Context, user *domain.User) (string, error) {
			return "", expectedErr // Simulate token signing error
		},
	}

	userUseCase := usecases.NewUserUseCase(mockUserRepo, mockJwtService, mockPasswordService)

	_, err := userUseCase.Login(ctx, testUsername, testPassword)

	if err == nil {
		t.Fatal("Login did not return error for token signing failure")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Error mismatch: got %v, want error containing %v", err, expectedErr)
	}
}

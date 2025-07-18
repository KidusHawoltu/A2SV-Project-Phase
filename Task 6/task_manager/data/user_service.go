package data

import (
	"A2SV_ProjectPhase/Task6/TaskManager/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserCollection struct {
	collection *mongo.Collection
	jwtSecret  string // Secret key for JWT signing
}

// Define specific sentinel errors for authentication and user management
var ErrInvalidCredentials = errors.New("incorrect Username or Password")
var ErrUserNotFound = errors.New("user not found")
var ErrUsernameTaken = errors.New("username is already taken")

// UserManager interface defines operations for user management
type UserManager interface {
	RegisterUser(ctx context.Context, user models.User) (models.User, error)
	Login(ctx context.Context, username string, password string) (string, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
}

// NewUserManager creates a new UserManager instance
func NewUserManager(uc *mongo.Collection, secret string) UserManager {
	return &UserCollection{
		collection: uc,
		jwtSecret:  secret,
	}
}

// getByUsername is a private helper to retrieve a user by their username
func (us *UserCollection) GetByUsername(ctx context.Context, username string) (models.User, error) {
	filter := bson.M{"username": username}
	var user models.User
	err := us.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, ErrUserNotFound
		}
		// If it's not a "no documents" error, it's a general database retrieval error
		return models.User{}, fmt.Errorf("failed to retrieve user by username: %w", err)
	}
	return user, nil
}

// RegisterUser handles the registration of a new user
func (uc *UserCollection) RegisterUser(ctx context.Context, user models.User) (models.User, error) {
	// Check if username already exists
	_, err := uc.GetByUsername(ctx, user.Username)
	if err == nil { // If no error, it means a user with this username was found
		return models.User{}, ErrUsernameTaken
	}
	if !errors.Is(err, ErrUserNotFound) { // If it's an error other than "not found", it's a DB issue
		return models.User{}, fmt.Errorf("failed to check for existing user: %w", err)
	}

	// Hash the password before storing
	user.Id = primitive.NilObjectID // Ensure MongoDB generates a new ID
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword) // Store the hashed password

	// Insert the new user into the collection
	inserted, err := uc.collection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to add user: %w", err)
	}

	// Get the generated ObjectID and set it back to the user model
	id, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.User{}, fmt.Errorf("created ID is not a valid ObjectID: %v", inserted.InsertedID)
	}
	user.Id = id
	user.Password = ""
	return user, nil
}

// Login authenticates a user and generates a JWT
func (uc *UserCollection) Login(ctx context.Context, username string, password string) (string, error) {
	// Retrieve user by username
	user, err := uc.GetByUsername(ctx, username)
	if err != nil {
		return "", ErrInvalidCredentials // Return a generic invalid credentials error for security
	}

	// Compare provided password with the stored hashed password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", ErrInvalidCredentials // Password mismatch
	}

	// Create JWT claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.Id.Hex(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	// Sign the token with the secret key
	jwtToken, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return jwtToken, nil
}

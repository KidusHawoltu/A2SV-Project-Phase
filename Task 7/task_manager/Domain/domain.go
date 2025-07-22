package domain

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	Pending    TaskStatus = "Pending"
	InProgress TaskStatus = "In progress"
	Done       TaskStatus = "Done"
)

func (status TaskStatus) IsValid() bool {
	switch status {
	case Pending, InProgress, Done:
		return true
	}
	return false
}

type Task struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"duedate" bson:"duedate"`
	Status      TaskStatus         `json:"status" bson:"status"`
}

func NewTask(title string, description string, dueDate time.Time, status TaskStatus) (*Task, error) {
	if title == "" {
		return nil, errors.New("task title cannot be empty")
	}
	if !status.IsValid() {
		return nil, errors.New("invalid task status")
	}
	if dueDate.IsZero() {
		return nil, errors.New("task due date cannot be empty")
	}
	if dueDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("task due date cannot be in the past")
	}
	return &Task{
		Id:          primitive.NilObjectID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
	}, nil
}

type TaskRepository interface {
	CreateTask(c context.Context, task *Task) (*Task, error)
	GetTaskById(c context.Context, id primitive.ObjectID) (*Task, error)
	GetAllTasks(c context.Context) ([]*Task, error)
	UpdateTask(c context.Context, id primitive.ObjectID, task *Task) (*Task, error)
	DeleteTask(c context.Context, id primitive.ObjectID) error
}

type UserRole string

const (
	RoleAdmin UserRole = "Admin"
	RoleUser  UserRole = "User"
)

func (role UserRole) IsValid() bool {
	switch role {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

type User struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	PasswordHash string             `json:"-" bson:"password"`
	Role         UserRole           `json:"role" bson:"role"`
}

func NewUser(username string, hashedPassword string) (*User, error) {
	if username == "" || hashedPassword == "" {
		return nil, errors.New("missing required user fields for new user")
	}
	return &User{
		Id:           primitive.NilObjectID,
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         RoleUser,
	}, nil
}

type UserRepository interface {
	CreateUser(c context.Context, user *User) (*User, error)
	GetUserByUsername(c context.Context, username string) (*User, error)
	// GetUserById(c context.Context, id primitive.ObjectID) (*User, error)
	// GetAllUsers(c context.Context) ([]*User, error)
	// UpdateUser(c context.Context, id primitive.ObjectID, user *User) (*User, error)
	// DeleteUser(c context.Context, id primitive.ObjectID) error
}

type PasswordService interface {
	Hash(c context.Context, password string) (string, error)
	Compare(c context.Context, password string, hashedPassword string) error
}

type Claims struct {
	UserId   string
	Role     UserRole
	Username string
	jwt.StandardClaims
}

type JwtService interface {
	GetSignedToken(c context.Context, user *User) (string, error)
	ParseToken(c context.Context, token string) (*Claims, error)
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTaskNotFound       = errors.New("task not found")
	ErrValidationFailed   = errors.New("validation failed")
)

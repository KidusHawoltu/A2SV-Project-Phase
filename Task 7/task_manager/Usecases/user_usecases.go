package usecases

import (
	domain "A2SV_ProjectPhase/Task7/TaskManager/Domain"
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo        domain.UserRepository
	jwtService      domain.JwtService
	passwordService domain.PasswordService
}

func NewUserUseCase(userrepo domain.UserRepository, jwtservice domain.JwtService, passwordservice domain.PasswordService) *UserUseCase {
	return &UserUseCase{
		userRepo:        userrepo,
		jwtService:      jwtservice,
		passwordService: passwordservice,
	}
}

func (uc *UserUseCase) RegisterUser(c context.Context, username string, password string) (*domain.User, error) {
	existingUser, err := uc.userRepo.GetUserByUsername(c, username)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("usecase: failed to check exsisting user: %w", err)
	}
	if existingUser != nil {
		return nil, domain.ErrUsernameTaken
	}

	hashedPassword, err := uc.passwordService.Hash(c, password)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to hash password: %w", err)
	}

	newuser, err := domain.NewUser(username, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to create new user: %w", err)
	}

	savedUser, err := uc.userRepo.CreateUser(c, newuser)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to save user: %w", err)
	}

	return savedUser, nil
}

func (uc *UserUseCase) Login(c context.Context, username string, password string) (string, error) {
	existingUser, err := uc.userRepo.GetUserByUsername(c, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", fmt.Errorf("usecase: failed to check exsisting user: %w", err)
	}

	if err := uc.passwordService.Compare(c, password, existingUser.PasswordHash); err != nil {
		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Printf("usecase: failed to verify password for user %q: %v\n", username, err)
		}
		return "", domain.ErrInvalidCredentials
	}

	token, err := uc.jwtService.GetSignedToken(c, existingUser)
	if err != nil {
		return "", fmt.Errorf("usecase: failed to get token: %w", err)
	}
	return token, nil
}

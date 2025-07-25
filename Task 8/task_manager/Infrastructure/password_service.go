package infrastructure

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Ensure BcryptPasswordService implements the domain.PasswordService interface
var _ domain.PasswordService = (*BcryptPasswordService)(nil)

type BcryptPasswordService struct {
	cost int
}

func NewBcryptPasswordService(cost int) *BcryptPasswordService {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptPasswordService{cost: cost}
}

func (service *BcryptPasswordService) Hash(c context.Context, password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), service.cost)
	if err != nil {
		return "", fmt.Errorf("password service: failed to generate hash: %w", err)
	}
	return string(hashedBytes), nil
}

func (service *BcryptPasswordService) Compare(c context.Context, password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password service: password comparison failed: %w", err)
	}
	return nil // Passwords match
}

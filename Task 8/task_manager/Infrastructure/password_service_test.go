package infrastructure_test

import (
	"A2SV_ProjectPhase/Task8/TaskManager/Infrastructure"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

//===========================================================================
// BcryptPasswordService Test Suite
//===========================================================================

type BcryptPasswordServiceSuite struct {
	suite.Suite
	passwordService *infrastructure.BcryptPasswordService
}

// TestBcryptPasswordServiceSuite is the entry point for the test suite
func TestBcryptPasswordServiceSuite(t *testing.T) {
	suite.Run(t, new(BcryptPasswordServiceSuite))
}

// SetupTest runs before each test method
func (s *BcryptPasswordServiceSuite) SetupTest() {
	// Use a low cost for testing to make hashing very fast
	s.passwordService = infrastructure.NewBcryptPasswordService(bcrypt.MinCost)
}

// TestHashAndCompare tests the full round-trip of hashing and comparing.
func (s *BcryptPasswordServiceSuite) TestHashAndCompare() {
	password := "my-s3cr3t-p@ssw0rd"

	s.Run("Success - Correct Password", func() {
		// --- Execution: Hash ---
		hashedPassword, err := s.passwordService.Hash(context.Background(), password)

		// --- Assertion: Hash ---
		s.Require().NoError(err, "Hashing should not produce an error")
		s.NotEmpty(hashedPassword, "Hashed password should not be empty")
		// A bcrypt hash is typically 60 characters long
		s.Len(hashedPassword, 60)
		// The hash should not be the same as the original password
		s.NotEqual(password, hashedPassword)

		// --- Execution: Compare ---
		err = s.passwordService.Compare(context.Background(), password, hashedPassword)

		// --- Assertion: Compare ---
		s.Require().NoError(err, "Comparison with the correct password should succeed")
	})

	s.Run("Failure - Incorrect Password", func() {
		// --- Execution: Hash ---
		hashedPassword, err := s.passwordService.Hash(context.Background(), password)
		s.Require().NoError(err)

		// --- Execution: Compare ---
		incorrectPassword := "my-wrong-password"
		err = s.passwordService.Compare(context.Background(), incorrectPassword, hashedPassword)

		// --- Assertion: Compare ---
		s.Require().Error(err, "Comparison with an incorrect password should fail")
		// We can check if the error is the specific one from bcrypt
		s.ErrorIs(err, bcrypt.ErrMismatchedHashAndPassword, "Error should be a bcrypt mismatch error")
	})
}

// TestNewBcryptPasswordService_DefaultCost ensures the constructor sets a default cost.
func (s *BcryptPasswordServiceSuite) TestNewBcryptPasswordService_DefaultCost() {
	// Execute the constructor with a zero cost
	service := infrastructure.NewBcryptPasswordService(0)
	s.Require().NotNil(service)

	// To verify the private 'cost' field, we have to test its behavior.
	// We'll hash a password and see if it works, which implies a valid cost was set.
	_, err := service.Hash(context.Background(), "test-password")
	s.NoError(err, "Service with default cost should be able to hash passwords")
}

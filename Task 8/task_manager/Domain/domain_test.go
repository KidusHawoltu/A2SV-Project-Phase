package domain_test

import (
	"A2SV_ProjectPhase/Task8/TaskManager/Domain" // Adjust your import path
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

//===========================================================================
// Task Test Suite
//===========================================================================

// TaskSuite defines the test suite for the Task domain object.
type TaskSuite struct {
	suite.Suite
}

// TestTaskSuite is the runner for the TaskSuite. 'go test' will find this function.
func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(TaskSuite))
}

// TestSuccess tests the successful creation of a Task.
func (s *TaskSuite) TestSuccess() {
	title := "Test Task"
	description := "Test Description"
	dueDate := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	status := domain.Pending

	task, err := domain.NewTask(title, description, dueDate, status)

	// Use Require for checks that must pass for the test to be valid.
	s.Require().NoError(err, "NewTask should not return an error on valid input")
	s.Require().NotNil(task, "NewTask should not return a nil task on valid input")

	// Use Equal (from Assert) for checking individual fields.
	s.Equal(title, task.Title)
	s.Equal(description, task.Description)
	s.Equal(dueDate, task.DueDate)
	s.Equal(status, task.Status)
}

// TestValidation consolidates all validation failure tests for NewTask.
func (s *TaskSuite) TestValidation() {
	testCases := []struct {
		name          string
		title         string
		description   string
		dueDate       time.Time
		status        domain.TaskStatus
		expectedError string
	}{
		{
			name:          "Empty Title",
			title:         "",
			description:   "desc",
			dueDate:       time.Now().Add(time.Hour),
			status:        domain.Pending,
			expectedError: "task title cannot be empty",
		},
		{
			name:          "Invalid Status",
			title:         "Title",
			description:   "desc",
			dueDate:       time.Now().Add(time.Hour),
			status:        "InvalidStatus",
			expectedError: "invalid task status",
		},
		{
			name:          "Empty Due Date",
			title:         "Title",
			description:   "desc",
			dueDate:       time.Time{}, // Zero value for time
			status:        domain.Pending,
			expectedError: "task due date cannot be empty",
		},
		{
			name:          "Past Due Date",
			title:         "Title",
			description:   "desc",
			dueDate:       time.Now().Add(-24 * time.Hour),
			status:        domain.Pending,
			expectedError: "task due date cannot be in the past",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := domain.NewTask(tc.title, tc.description, tc.dueDate, tc.status)
			s.Require().Error(err, "Expected an error for invalid input")
			s.Equal(tc.expectedError, err.Error(), "Error message mismatch")
		})
	}
}

// TestStatusIsValid tests the IsValid method for TaskStatus.
func (s *TaskSuite) TestStatusIsValid() {
	testCases := []struct {
		status domain.TaskStatus
		want   bool
	}{
		{domain.Pending, true},
		{domain.InProgress, true},
		{domain.Done, true},
		{"Invalid", false},
		{"", false},
	}

	for _, tc := range testCases {
		s.Run(string(tc.status), func() {
			got := tc.status.IsValid()
			s.Equal(tc.want, got)
		})
	}
}

//===========================================================================
// User Test Suite
//===========================================================================

// UserSuite defines the test suite for the User domain object.
type UserSuite struct {
	suite.Suite
}

// TestUserSuite is the runner for the UserSuite.
func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

// TestSuccess tests the successful creation of a User.
func (s *UserSuite) TestSuccess() {
	username := "testuser"
	hashedPassword := "a-very-secure-hashed-password"

	user, err := domain.NewUser(username, hashedPassword)

	s.Require().NoError(err, "NewUser should not return an error on valid input")
	s.Require().NotNil(user, "NewUser should not return a nil user on valid input")

	s.Equal(username, user.Username)
	s.Equal(hashedPassword, user.PasswordHash)
	s.Equal(domain.RoleUser, user.Role, "Default role should be 'user'")
}

// TestValidation_MissingFields tests user creation with missing fields.
func (s *UserSuite) TestValidation_MissingFields() {
	testCases := []struct {
		name     string
		username string
		password string
	}{
		{"Empty username", "", "pass"},
		{"Empty password", "user", ""},
		{"Both empty", "", ""},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := domain.NewUser(tc.username, tc.password)
			s.Require().Error(err)
			s.ErrorIs(err, domain.ErrValidationFailed)
		})
	}
}

// TestRoleIsValid tests the IsValid method for UserRole.
func (s *UserSuite) TestRoleIsValid() {
	testCases := []struct {
		role domain.UserRole
		want bool
	}{
		{domain.RoleAdmin, true},
		{domain.RoleUser, true},
		{"Guest", false},
		{"", false},
	}

	for _, tc := range testCases {
		s.Run(string(tc.role), func() {
			got := tc.role.IsValid()
			s.Equal(tc.want, got)
		})
	}
}

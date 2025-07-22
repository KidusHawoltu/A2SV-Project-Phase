package domain_test

import (
	"A2SV_ProjectPhase/Task7/TaskManager/Domain" // Import the domain package
	"testing"
	"time"
)

// TestNewTask_Success tests successful Task creation
func TestNewTask_Success(t *testing.T) {
	title := "Test Task"
	description := "Test Description"
	dueDate := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour) // Future date, truncated
	status := domain.Pending

	task, err := domain.NewTask(title, description, dueDate, status)

	if err != nil {
		t.Fatalf("NewTask() returned unexpected error: %v", err)
	}
	if task == nil {
		t.Fatal("NewTask() returned nil task")
	}
	if task.Title != title {
		t.Errorf("NewTask() Title mismatch: got %s, want %s", task.Title, title)
	}
	if task.Description != description {
		t.Errorf("NewTask() Description mismatch: got %s, want %s", task.Description, description)
	}
	if !task.DueDate.Equal(dueDate) {
		t.Errorf("NewTask() DueDate mismatch: got %v, want %v", task.DueDate, dueDate)
	}
	if task.Status != status {
		t.Errorf("NewTask() Status mismatch: got %s, want %s", task.Status, status)
	}
}

// TestNewTask_Failure_EmptyTitle tests creation with an empty title
func TestNewTask_Failure_EmptyTitle(t *testing.T) {
	_, err := domain.NewTask("", "desc", time.Now().Add(time.Hour), domain.Pending)
	if err == nil {
		t.Fatal("NewTask() did not return error for empty title")
	}
	if err.Error() != "task title cannot be empty" {
		t.Errorf("NewTask() error message mismatch: got %q, want %q", err.Error(), "task title cannot be empty")
	}
}

// TestNewTask_Failure_InvalidStatus tests creation with an invalid status
func TestNewTask_Failure_InvalidStatus(t *testing.T) {
	_, err := domain.NewTask("Title", "desc", time.Now().Add(time.Hour), "InvalidStatus")
	if err == nil {
		t.Fatal("NewTask() did not return error for invalid status")
	}
	if err.Error() != "invalid task status" {
		t.Errorf("NewTask() error message mismatch: got %q, want %q", err.Error(), "invalid task status")
	}
}

// TestNewTask_Failure_EmptyDueDate tests creation with a zero/empty due date
func TestNewTask_Failure_EmptyDueDate(t *testing.T) {
	_, err := domain.NewTask("Title", "desc", time.Time{}, domain.Pending)
	if err == nil {
		t.Fatal("NewTask() did not return error for empty due date")
	}
	if err.Error() != "task due date cannot be empty" {
		t.Errorf("NewTask() error message mismatch: got %q, want %q", err.Error(), "task due date cannot be empty")
	}
}

// TestNewTask_Failure_PastDueDate tests creation with a due date in the past
func TestNewTask_Failure_PastDueDate(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour) // 24 hours ago
	_, err := domain.NewTask("Title", "desc", pastDate, domain.Pending)
	if err == nil {
		t.Fatal("NewTask() did not return error for past due date")
	}
	if err.Error() != "task due date cannot be in the past" {
		t.Errorf("NewTask() error message mismatch: got %q, want %q", err.Error(), "task due date cannot be in the past")
	}
}

// TestNewUser_Success tests successful User creation
func TestNewUser_Success(t *testing.T) {
	username := "testuser"
	hashedPassword := "hashedpassword"

	user, err := domain.NewUser(username, hashedPassword)

	if err != nil {
		t.Fatalf("NewUser() returned unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("NewUser() returned nil user")
	}
	if user.Username != username {
		t.Errorf("NewUser() Username mismatch: got %s, want %s", user.Username, username)
	}
	if user.PasswordHash != hashedPassword {
		t.Errorf("NewUser() PasswordHash mismatch: got %s, want %s", user.PasswordHash, hashedPassword)
	}
	if user.Role != domain.RoleUser {
		t.Errorf("NewUser() Role mismatch: got %s, want %s", user.Role, domain.RoleUser)
	}
}

// TestNewUser_Failure_MissingFields tests creation with missing fields
func TestNewUser_Failure_MissingFields(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{"Empty username", "", "pass"},
		{"Empty password", "user", ""},
		{"Both empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewUser(tt.username, tt.password)
			if err == nil {
				t.Fatalf("NewUser() did not return error for %s", tt.name)
			}
			if err.Error() != "missing required user fields for new user" {
				t.Errorf("NewUser() error message mismatch: got %q, want %q", err.Error(), "missing required user fields for new user")
			}
		})
	}
}

// TestTaskStatus_IsValid tests TaskStatus IsValid method
func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		status domain.TaskStatus
		want   bool
	}{
		{domain.Pending, true},
		{domain.InProgress, true},
		{domain.Done, true},
		{"Invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("TaskStatus.IsValid() got %v, want %v for status %q", got, tt.want, tt.status)
			}
		})
	}
}

// TestUserRole_IsValid tests UserRole IsValid method
func TestUserRole_IsValid(t *testing.T) {
	tests := []struct {
		role domain.UserRole
		want bool
	}{
		{domain.RoleAdmin, true},
		{domain.RoleUser, true},
		{"Guest", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("UserRole.IsValid() got %v, want %v for role %q", got, tt.want, tt.role)
			}
		})
	}
}

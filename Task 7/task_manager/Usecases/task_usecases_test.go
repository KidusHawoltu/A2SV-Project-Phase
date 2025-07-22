package usecases_test

import (
	domain "A2SV_ProjectPhase/Task7/TaskManager/Domain"
	usecases "A2SV_ProjectPhase/Task7/TaskManager/Usecases"
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mocks for TaskUseCase Dependencies ---

// MockTaskRepository implements domain.TaskRepository
type MockTaskRepository struct {
	CreateTaskFunc  func(c context.Context, task *domain.Task) (*domain.Task, error)
	GetTaskByIdFunc func(c context.Context, id primitive.ObjectID) (*domain.Task, error)
	GetAllTasksFunc func(c context.Context) ([]*domain.Task, error)
	UpdateTaskFunc  func(c context.Context, id primitive.ObjectID, task *domain.Task) (*domain.Task, error)
	DeleteTaskFunc  func(c context.Context, id primitive.ObjectID) error
}

func (m *MockTaskRepository) CreateTask(c context.Context, task *domain.Task) (*domain.Task, error) {
	return m.CreateTaskFunc(c, task)
}
func (m *MockTaskRepository) GetTaskById(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
	return m.GetTaskByIdFunc(c, id)
}
func (m *MockTaskRepository) GetAllTasks(c context.Context) ([]*domain.Task, error) {
	return m.GetAllTasksFunc(c)
}
func (m *MockTaskRepository) UpdateTask(c context.Context, id primitive.ObjectID, task *domain.Task) (*domain.Task, error) {
	return m.UpdateTaskFunc(c, id, task)
}
func (m *MockTaskRepository) DeleteTask(c context.Context, id primitive.ObjectID) error {
	return m.DeleteTaskFunc(c, id)
}

// --- Tests for TaskUseCase ---

func TestCreateTask_Success(t *testing.T) {
	ctx := getTestContext()
	title := "New Task"
	description := "Description"
	dueDate := time.Now().Add(time.Hour * 24).Truncate(24 * time.Hour)
	status := domain.Pending

	mockTask := &domain.Task{
		Id:          primitive.NewObjectID(),
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      status,
	}

	mockTaskRepo := &MockTaskRepository{
		CreateTaskFunc: func(c context.Context, task *domain.Task) (*domain.Task, error) {
			task.Id = mockTask.Id // Simulate MongoDB assigning ID
			return task, nil
		},
	}

	taskUseCase := usecases.NewTaskUseCase(mockTaskRepo)
	createdTask, err := taskUseCase.CreateTask(ctx, title, description, dueDate, status)

	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if createdTask == nil {
		t.Fatal("Created task is nil")
	}
	if createdTask.Id.IsZero() {
		t.Error("Task ID was not set by repository")
	}
}

func TestCreateTask_ValidationFailed(t *testing.T) {
	ctx := getTestContext()
	// Case 1: Empty title
	_, err := usecases.NewTaskUseCase(&MockTaskRepository{}).CreateTask(ctx, "", "desc", time.Now().Add(time.Hour), domain.Pending)
	if err == nil || !errors.Is(err, domain.ErrValidationFailed) {
		t.Errorf("CreateTask did not return ErrValidationFailed for empty title: %v", err)
	}

	// Case 2: Past due date (checked in domain.NewTask)
	_, err = usecases.NewTaskUseCase(&MockTaskRepository{}).CreateTask(ctx, "Title", "desc", time.Now().Add(-24*time.Hour), domain.Pending)
	if err == nil || !errors.Is(err, domain.ErrValidationFailed) {
		t.Errorf("CreateTask did not return ErrValidationFailed for past due date: %v", err)
	}
}

func TestGetTaskByID_Success(t *testing.T) {
	ctx := getTestContext()
	taskID := primitive.NewObjectID()
	expectedTask := &domain.Task{Id: taskID, Title: "Test Task"}

	mockTaskRepo := &MockTaskRepository{
		GetTaskByIdFunc: func(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
			if id != taskID {
				return nil, errors.New("ID mismatch in mock")
			}
			return expectedTask, nil
		},
	}

	taskUseCase := usecases.NewTaskUseCase(mockTaskRepo)
	retrievedTask, err := taskUseCase.GetTaskByID(ctx, taskID.Hex())

	if err != nil {
		t.Fatalf("GetTaskByID failed: %v", err)
	}
	if retrievedTask == nil {
		t.Fatal("Retrieved task is nil")
	}
	if retrievedTask.Id != taskID {
		t.Errorf("Task ID mismatch: got %s, want %s", retrievedTask.Id.Hex(), taskID)
	}
}

func TestGetTaskByID_NotFound(t *testing.T) {
	ctx := getTestContext()
	taskID := primitive.NewObjectID().Hex()

	mockTaskRepo := &MockTaskRepository{
		GetTaskByIdFunc: func(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
			return nil, domain.ErrTaskNotFound
		},
	}

	taskUseCase := usecases.NewTaskUseCase(mockTaskRepo)
	_, err := taskUseCase.GetTaskByID(ctx, taskID)

	if err == nil || !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("GetTaskByID did not return ErrTaskNotFound: %v", err)
	}
}

func TestGetTaskByID_InvalidIDFormat(t *testing.T) {
	ctx := getTestContext()
	taskID := "invalid_id"

	taskUseCase := usecases.NewTaskUseCase(&MockTaskRepository{}) // Repo not called
	_, err := taskUseCase.GetTaskByID(ctx, taskID)

	if err == nil || !errors.Is(err, domain.ErrValidationFailed) {
		t.Errorf("GetTaskByID did not return ErrValidationFailed for invalid ID format: %v", err)
	}
}

func TestGetAllTasks_Success(t *testing.T) {
	ctx := getTestContext()
	expectedTasks := []*domain.Task{
		{Id: primitive.NewObjectID(), Title: "Task 1"},
		{Id: primitive.NewObjectID(), Title: "Task 2"},
	}

	mockTaskRepo := &MockTaskRepository{
		GetAllTasksFunc: func(c context.Context) ([]*domain.Task, error) {
			return expectedTasks, nil
		},
	}

	taskUseCase := usecases.NewTaskUseCase(mockTaskRepo)
	tasks, err := taskUseCase.GetAllTasks(ctx)

	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}
	if len(tasks) != len(expectedTasks) {
		t.Errorf("Task count mismatch: got %d, want %d", len(tasks), len(expectedTasks))
	}
}

func TestUpdateTask_Success(t *testing.T) {
	ctx := getTestContext()
	taskID := primitive.NewObjectID()
	originalTask := &domain.Task{
		Id: taskID, Title: "Old Title", Description: "Old Desc",
		DueDate: time.Now().Add(time.Hour).Truncate(time.Hour), Status: domain.Pending,
	}

	newTitle := "New Title"
	newDescription := "New Desc"
	newDueDate := time.Now().Add(2 * time.Hour).Truncate(time.Hour)
	newStatus := domain.InProgress

	mockTaskRepo := &MockTaskRepository{
		GetTaskByIdFunc: func(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
			return originalTask, nil
		},
		UpdateTaskFunc: func(c context.Context, id primitive.ObjectID, task *domain.Task) (*domain.Task, error) {
			return task, nil // Return the updated task back
		},
	}

	taskUseCase := usecases.NewTaskUseCase(mockTaskRepo)
	updatedTask, err := taskUseCase.UpdateTask(ctx, taskID.Hex(), &newTitle, &newDescription, &newDueDate, &newStatus)

	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}
	if updatedTask == nil {
		t.Fatal("Updated task is nil")
	}
	if updatedTask.Title != newTitle || updatedTask.Description != newDescription || !updatedTask.DueDate.Equal(newDueDate) || updatedTask.Status != newStatus {
		t.Errorf("Task fields not updated correctly")
	}
}

func TestUpdateTask_StatusChangeFromDone(t *testing.T) {
	ctx := getTestContext()
	taskID := primitive.NewObjectID()
	originalTask := &domain.Task{
		Id: taskID, Title: "Title", Description: "Desc",
		DueDate: time.Now().Add(time.Hour), Status: domain.Done,
	}
	newStatus := domain.Pending // Try to change from Done

	mockTaskRepo := &MockTaskRepository{
		GetTaskByIdFunc: func(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
			return originalTask, nil
		},
	}

	taskUseCase := usecases.NewTaskUseCase(mockTaskRepo)
	_, err := taskUseCase.UpdateTask(ctx, taskID.Hex(), nil, nil, nil, &newStatus)

	if err == nil || !errors.Is(err, domain.ErrValidationFailed) {
		t.Errorf("UpdateTask did not return ErrValidationFailed for status change from Done: %v", err)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	ctx := getTestContext()
	taskID := primitive.NewObjectID().Hex()

	mockTaskRepo := &MockTaskRepository{
		DeleteTaskFunc: func(c context.Context, id primitive.ObjectID) error {
			return nil
		},
	}

	taskUseCase := usecases.NewTaskUseCase(mockTaskRepo)
	err := taskUseCase.DeleteTask(ctx, taskID)

	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	ctx := getTestContext()
	taskID := primitive.NewObjectID().Hex()

	mockTaskRepo := &MockTaskRepository{
		DeleteTaskFunc: func(c context.Context, id primitive.ObjectID) error {
			return domain.ErrTaskNotFound
		},
	}

	taskUseCase := usecases.NewTaskUseCase(mockTaskRepo)
	err := taskUseCase.DeleteTask(ctx, taskID)

	if err == nil || !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("DeleteTask did not return ErrTaskNotFound: %v", err)
	}
}

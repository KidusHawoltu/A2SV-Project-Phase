package usecases_test

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	usecases "A2SV_ProjectPhase/Task8/TaskManager/Usecases"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mock stays the same, as it's a good pattern ---
type MockTaskRepository struct {
	CreateTaskFunc  func(c context.Context, task *domain.Task) (*domain.Task, error)
	GetTaskByIdFunc func(c context.Context, id primitive.ObjectID) (*domain.Task, error)
	GetAllTasksFunc func(c context.Context) ([]*domain.Task, error)
	UpdateTaskFunc  func(c context.Context, id primitive.ObjectID, task *domain.Task) (*domain.Task, error)
	DeleteTaskFunc  func(c context.Context, id primitive.ObjectID) error
}

func (m *MockTaskRepository) CreateTask(c context.Context, task *domain.Task) (*domain.Task, error) {
	if m.CreateTaskFunc != nil {
		return m.CreateTaskFunc(c, task)
	}
	return nil, errors.New("CreateTaskFunc not implemented")
}
func (m *MockTaskRepository) GetTaskById(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
	if m.GetTaskByIdFunc != nil {
		return m.GetTaskByIdFunc(c, id)
	}
	return nil, errors.New("GetTaskByIdFunc not implemented")
}
func (m *MockTaskRepository) GetAllTasks(c context.Context) ([]*domain.Task, error) {
	if m.GetAllTasksFunc != nil {
		return m.GetAllTasksFunc(c)
	}
	return nil, errors.New("GetAllTasksFunc not implemented")
}
func (m *MockTaskRepository) UpdateTask(c context.Context, id primitive.ObjectID, task *domain.Task) (*domain.Task, error) {
	if m.UpdateTaskFunc != nil {
		return m.UpdateTaskFunc(c, id, task)
	}
	return nil, errors.New("UpdateTaskFunc not implemented")
}
func (m *MockTaskRepository) DeleteTask(c context.Context, id primitive.ObjectID) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(c, id)
	}
	return errors.New("DeleteTaskFunc not implemented")
}

//===========================================================================
// TaskUseCase Test Suite
//===========================================================================

type TaskUseCaseSuite struct {
	suite.Suite
	mockRepo *MockTaskRepository
	useCase  *usecases.TaskUseCase
	ctx      context.Context
}

// TestTaskUseCaseSuite is the entry point for the test suite
func TestTaskUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseSuite))
}

// SetupTest runs before each test method in the suite.
// It's the perfect place to initialize mocks and the system under test.
func (s *TaskUseCaseSuite) SetupTest() {
	s.mockRepo = &MockTaskRepository{}
	s.useCase = usecases.NewTaskUseCase(s.mockRepo)
	s.ctx = context.Background() // A basic context is fine for these tests
}

// --- Test Methods for TaskUseCase ---

func (s *TaskUseCaseSuite) TestCreateTask() {
	s.Run("Success", func() {
		// Reset mock for this subtest if needed, although SetupTest already does it.
		s.SetupTest()

		title := "New Task"
		description := "Description"
		dueDate := time.Now().Add(time.Hour * 24).Truncate(24 * time.Hour)
		status := domain.Pending

		// Configure the mock's behavior for this specific test
		s.mockRepo.CreateTaskFunc = func(c context.Context, task *domain.Task) (*domain.Task, error) {
			task.Id = primitive.NewObjectID() // Simulate repository assigning an ID
			return task, nil
		}

		createdTask, err := s.useCase.CreateTask(s.ctx, title, description, dueDate, status)

		s.Require().NoError(err)
		s.Require().NotNil(createdTask)
		s.False(createdTask.Id.IsZero(), "Task ID should be set by the repository")
		s.Equal(title, createdTask.Title)
	})

	s.Run("Validation Failed", func() {
		s.SetupTest()
		// No mock setup needed, as validation should fail before the repo is called.

		_, err := s.useCase.CreateTask(s.ctx, "", "desc", time.Now().Add(time.Hour), domain.Pending)
		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrValidationFailed, "Should return validation error for empty title")
	})
}

func (s *TaskUseCaseSuite) TestGetTaskByID() {
	s.Run("Success", func() {
		s.SetupTest()
		taskID := primitive.NewObjectID()
		expectedTask := &domain.Task{Id: taskID, Title: "Test Task"}

		s.mockRepo.GetTaskByIdFunc = func(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
			s.Equal(taskID, id, "ID passed to repository should match")
			return expectedTask, nil
		}

		retrievedTask, err := s.useCase.GetTaskByID(s.ctx, taskID.Hex())

		s.Require().NoError(err)
		s.Require().NotNil(retrievedTask)
		s.Equal(expectedTask.Id, retrievedTask.Id)
	})

	s.Run("Not Found", func() {
		s.SetupTest()
		taskID := primitive.NewObjectID()

		s.mockRepo.GetTaskByIdFunc = func(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
			return nil, domain.ErrTaskNotFound
		}

		_, err := s.useCase.GetTaskByID(s.ctx, taskID.Hex())
		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrTaskNotFound)
	})

	s.Run("Invalid ID Format", func() {
		s.SetupTest()
		// Repo will not be called, so no mock setup needed.
		_, err := s.useCase.GetTaskByID(s.ctx, "this-is-not-a-valid-hex-id")

		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrValidationFailed)
	})
}

func (s *TaskUseCaseSuite) TestGetAllTasks() {
	s.Run("Success", func() {
		s.SetupTest()
		expectedTasks := []*domain.Task{
			{Id: primitive.NewObjectID(), Title: "Task 1"},
			{Id: primitive.NewObjectID(), Title: "Task 2"},
		}

		s.mockRepo.GetAllTasksFunc = func(c context.Context) ([]*domain.Task, error) {
			return expectedTasks, nil
		}

		tasks, err := s.useCase.GetAllTasks(s.ctx)

		s.Require().NoError(err)
		s.Len(tasks, 2)
		s.Equal(expectedTasks, tasks)
	})
}

func (s *TaskUseCaseSuite) TestUpdateTask() {
	s.Run("Success", func() {
		s.SetupTest()
		taskID := primitive.NewObjectID()
		originalTask := &domain.Task{
			Id: taskID, Title: "Old Title", Status: domain.Pending,
		}
		newTitle := "New Title"
		newStatus := domain.InProgress

		// Mock the Get call first
		s.mockRepo.GetTaskByIdFunc = func(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
			return originalTask, nil
		}
		// Mock the Update call
		s.mockRepo.UpdateTaskFunc = func(c context.Context, id primitive.ObjectID, task *domain.Task) (*domain.Task, error) {
			return task, nil // Echo back the updated task
		}

		updatedTask, err := s.useCase.UpdateTask(s.ctx, taskID.Hex(), &newTitle, nil, nil, &newStatus)

		s.Require().NoError(err)
		s.Require().NotNil(updatedTask)
		s.Equal(newTitle, updatedTask.Title)
		s.Equal(newStatus, updatedTask.Status)
	})

	s.Run("Validation Failed - Cannot Change Status From Done", func() {
		s.SetupTest()
		taskID := primitive.NewObjectID()
		doneTask := &domain.Task{Id: taskID, Status: domain.Done}
		newStatus := domain.InProgress

		s.mockRepo.GetTaskByIdFunc = func(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
			return doneTask, nil
		}

		_, err := s.useCase.UpdateTask(s.ctx, taskID.Hex(), nil, nil, nil, &newStatus)

		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrValidationFailed)
	})
}

func (s *TaskUseCaseSuite) TestDeleteTask() {
	s.Run("Success", func() {
		s.SetupTest()
		taskID := primitive.NewObjectID()
		s.mockRepo.DeleteTaskFunc = func(c context.Context, id primitive.ObjectID) error {
			s.Equal(taskID, id)
			return nil
		}

		err := s.useCase.DeleteTask(s.ctx, taskID.Hex())

		s.Require().NoError(err)
	})

	s.Run("Not Found", func() {
		s.SetupTest()
		taskID := primitive.NewObjectID()
		s.mockRepo.DeleteTaskFunc = func(c context.Context, id primitive.ObjectID) error {
			return domain.ErrTaskNotFound
		}

		err := s.useCase.DeleteTask(s.ctx, taskID.Hex())

		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrTaskNotFound)
	})
}

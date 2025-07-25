package repositories_test

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	"A2SV_ProjectPhase/Task8/TaskManager/Repositories"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//===========================================================================
// TaskRepo Integration Test Suite
//===========================================================================

type TaskRepoSuite struct {
	suite.Suite
	db   *mongo.Database
	coll *mongo.Collection
	repo domain.TaskRepository
}

// TestTaskRepoSuite is the entry point for the test suite
func TestTaskRepoSuite(t *testing.T) {
	// This check ensures that these integration tests only run when a DB is available.
	// It assumes the TestMain function has successfully connected and set the global variable.
	if testMongoClient == nil {
		t.Skip("Skipping integration tests: MongoDB connection not available.")
	}
	suite.Run(t, new(TaskRepoSuite))
}

// SetupSuite runs once for the entire suite.
// It gets the database and collection handles from the global client.
func (s *TaskRepoSuite) SetupSuite() {
	s.db = testMongoClient.Database("test_learning_phase")
	s.coll = s.db.Collection("task8")
}

// SetupTest runs before EACH test method.
// This is crucial for test isolation in a live database environment.
func (s *TaskRepoSuite) SetupTest() {
	// 1. Clean the collection to ensure a pristine state for each test.
	_, err := s.coll.DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err, "Failed to clean task collection before test")

	// 2. Create a new repository instance for the test.
	s.repo = repositories.NewMongoDBTaskRepository(s.coll)
}

// TestCreateTask tests the CreateTask repository method against a live DB.
func (s *TaskRepoSuite) TestCreateTask() {
	taskToCreate := &domain.Task{
		Title:       "Live DB Test",
		Description: "A task to test creation against a live database.",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      domain.Pending,
	}

	createdTask, err := s.repo.CreateTask(context.Background(), taskToCreate)

	s.Require().NoError(err, "CreateTask should not return an error")
	s.Require().NotNil(createdTask)
	s.False(createdTask.Id.IsZero(), "Created task ID should be set by the DB")
	s.Equal(taskToCreate.Title, createdTask.Title)
}

// TestGetTaskById tests finding a task by its ID.
func (s *TaskRepoSuite) TestGetTaskById() {
	s.Run("Success - Task Found", func() {
		// Setup: Seed the database with a task to find.
		taskToFind := &domain.Task{Id: primitive.NewObjectID(), Title: "Find Me"}
		_, err := s.coll.InsertOne(context.Background(), taskToFind)
		s.Require().NoError(err, "Failed to seed database for test")

		// Execution
		foundTask, err := s.repo.GetTaskById(context.Background(), taskToFind.Id)

		// Assertion
		s.Require().NoError(err)
		s.Require().NotNil(foundTask)
		s.Equal(taskToFind.Id, foundTask.Id)
	})

	s.Run("Failure - Task Not Found", func() {
		nonExistentID := primitive.NewObjectID()

		// Execution
		_, err := s.repo.GetTaskById(context.Background(), nonExistentID)

		// Assertion
		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrTaskNotFound)
	})
}

// TestGetAllTasks tests retrieving all tasks.
func (s *TaskRepoSuite) TestGetAllTasks() {
	// Setup: Seed with multiple tasks
	tasksToInsert := []interface{}{
		&domain.Task{Id: primitive.NewObjectID(), Title: "Task 1"},
		&domain.Task{Id: primitive.NewObjectID(), Title: "Task 2"},
	}
	_, err := s.coll.InsertMany(context.Background(), tasksToInsert)
	s.Require().NoError(err)

	// Execution
	allTasks, err := s.repo.GetAllTasks(context.Background())

	// Assertion
	s.Require().NoError(err)
	s.Len(allTasks, 2, "Expected to retrieve 2 tasks")
}

// TestUpdateTask tests the update functionality.
func (s *TaskRepoSuite) TestUpdateTask() {
	// Setup: Seed the database
	originalTask := &domain.Task{
		Id:     primitive.NewObjectID(),
		Title:  "Original Title",
		Status: domain.Pending,
	}
	_, err := s.coll.InsertOne(context.Background(), originalTask)
	s.Require().NoError(err)

	// Execution
	taskWithUpdates := &domain.Task{
		Title:  "Updated Title",
		Status: domain.InProgress,
	}
	updatedTask, err := s.repo.UpdateTask(context.Background(), originalTask.Id, taskWithUpdates)

	// Assertion
	s.Require().NoError(err)
	s.Require().NotNil(updatedTask)
	s.Equal(originalTask.Id, updatedTask.Id)
	s.Equal("Updated Title", updatedTask.Title)
}

// TestDeleteTask tests the deletion functionality.
func (s *TaskRepoSuite) TestDeleteTask() {
	// Setup: Seed a task to delete
	taskToDelete := &domain.Task{Id: primitive.NewObjectID(), Title: "Delete Me"}
	_, err := s.coll.InsertOne(context.Background(), taskToDelete)
	s.Require().NoError(err)

	// Execution
	err = s.repo.DeleteTask(context.Background(), taskToDelete.Id)

	// Assertion
	s.Require().NoError(err)

	// Verification: Ensure it's actually gone by trying to get it
	_, err = s.repo.GetTaskById(context.Background(), taskToDelete.Id)
	s.ErrorIs(err, domain.ErrTaskNotFound, "Task should not be found after deletion")
}

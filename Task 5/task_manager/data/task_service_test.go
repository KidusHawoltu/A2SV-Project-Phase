package data

import (
	"A2SV_ProjectPhase/Task5/TaskManager/models"
	"context"
	"errors"
	"log"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	testClient     *mongo.Client
	testCollection *mongo.Collection
	testManager    TaskManager
	databaseName   = "test_learning_phase" // Use a separate DB for testing!
	collectionName = "test_tasks"          // Use a separate collection for testing!
)

// TestMain runs before all tests in the package.
// This is where you set up and tear down your MongoDB connection.
func TestMain(m *testing.M) {
	// Load .env file for tests.
	// If your .env is in the project root and tests are in a sub-directory,
	// you might need to specify the path like `godotenv.Load("../.env")`.
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found for tests or error loading .env. Assuming environment variables are set directly.")
	}

	// Prioritize MONGO_TEST_URI for tests. Fallback to MONGO_URI if not set.
	mongoURI := os.Getenv("MONGO_TEST_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_URI") // Fallback to main URI
		if mongoURI == "" {
			log.Fatalf("MONGO_TEST_URI or MONGO_URI environment variable not set for tests. Please set one in .env or your environment.")
		}
		log.Println("WARNING: MONGO_TEST_URI not set, using MONGO_URI for tests. Ensure this URI points to a TEST database to prevent data loss!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second) // Increased timeout for tests
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	var err error
	testClient, err = mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB for tests: %v", err)
	}

	// Ping the database to ensure connection
	if err = testClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB for tests: %v", err)
	}
	log.Println("Successfully connected to MongoDB for tests!")

	testCollection = testClient.Database(databaseName).Collection(collectionName)
	testManager = NewTaskManager(testCollection)

	// 2. Run all tests
	code := m.Run()

	// 3. Teardown MongoDB connection
	log.Println("Disconnecting from MongoDB test database...")
	if err := testClient.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting from MongoDB test database: %v", err)
	}
	log.Println("MongoDB test connection closed.")

	os.Exit(code)
}

// setupTestCollection cleans the test collection before each test.
// This ensures test isolation.
func setupTestCollection(t *testing.T) {
	// Pass context to Drop command
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Drop the entire collection to ensure a clean slate
	err := testCollection.Drop(ctx)
	if err != nil && err.Error() != "ns not found" { // "ns not found" is okay if collection doesn't exist yet
		t.Fatalf("Failed to drop test collection: %v", err)
	}
	log.Printf("Cleaned test collection '%s' for new test.", collectionName)
}

// addSampleTasks is a helper to add common data for tests that require initial tasks.
func addSampleTasks(ctx context.Context, t *testing.T) ([]models.Task, error) {
	tasks := []models.Task{
		{
			Title:       "Initial Task 1",
			Description: "First task for setup",
			Status:      models.Pending,
			DueDate:     time.Now().Add(time.Hour),
		},
		{
			Title:       "Initial Task 2",
			Description: "Second task for setup",
			Status:      models.InProgress,
			DueDate:     time.Now().Add(2 * time.Hour),
		},
	}

	// Add tasks one by one to get their generated ObjectIDs
	var insertedTasks []models.Task
	for _, task := range tasks {
		createdTask, err := testManager.AddTask(ctx, task)
		if err != nil {
			t.Fatalf("Failed to add sample task: %v", err)
		}
		insertedTasks = append(insertedTasks, createdTask)
	}
	return insertedTasks, nil
}

// TestAddTask tests the AddTask method.
func TestAddTask(t *testing.T) {
	setupTestCollection(t) // Ensure a clean collection for this test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Arrange
	taskToAdd := models.Task{
		Title:       "A New Task",
		Description: "Details of the new task",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      models.Pending,
	}

	// Act
	createdTask, err := testManager.AddTask(ctx, taskToAdd)

	// Assert
	assert.NoError(t, err, "AddTask should not return an error")
	assert.NotNil(t, createdTask, "Created task should not be nil")
	assert.False(t, createdTask.Id.IsZero(), "ID of the created task should not be empty") // ID is now primitive.ObjectID

	// Verify it was actually stored by trying to retrieve it
	retrievedTask, err := testManager.GetTaskById(ctx, createdTask.Id)
	assert.NoError(t, err)
	// Compare relevant fields, as time.Time might have subtle differences in stored/retrieved
	assert.Equal(t, createdTask.Id, retrievedTask.Id)
	assert.Equal(t, createdTask.Title, retrievedTask.Title)
	assert.Equal(t, createdTask.Description, retrievedTask.Description)
	assert.Equal(t, createdTask.Status, retrievedTask.Status)
	// For DueDate, assert approximately equal or ignore microsecond differences
	assert.WithinDuration(t, createdTask.DueDate, retrievedTask.DueDate, time.Second)
}

// TestGetTasks tests the GetTasks method.
func TestGetTasks(t *testing.T) {
	setupTestCollection(t) // Clean before this test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Arrange: Add some tasks using the helper
	sampleTasks, err := addSampleTasks(ctx, t)
	assert.NoError(t, err)
	assert.Len(t, sampleTasks, 2)

	// Act
	allTasks, err := testManager.GetTasks(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, allTasks, "The returned slice of tasks should not be nil")
	assert.Len(t, allTasks, 2, "Should return a slice containing 2 tasks")

	// Sort the slice by ID (or title for consistent checking) because MongoDB doesn't guarantee order without explicit sort.
	sort.Slice(allTasks, func(i, j int) bool {
		return allTasks[i].Title < allTasks[j].Title // Using title for consistent ordering here
	})

	assert.Equal(t, "Initial Task 1", allTasks[0].Title)
	assert.Equal(t, "Initial Task 2", allTasks[1].Title)
}

// TestGetTaskById tests the GetTaskById method for success and failure cases.
func TestGetTaskById(t *testing.T) {
	setupTestCollection(t) // Clean before this test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Arrange: Add sample tasks to retrieve
	sampleTasks, err := addSampleTasks(ctx, t)
	assert.NoError(t, err)
	assert.Len(t, sampleTasks, 2)

	// Sort by title to ensure we know which ID corresponds to which task
	sort.Slice(sampleTasks, func(i, j int) bool {
		return sampleTasks[i].Title < sampleTasks[j].Title
	})

	task1ID := sampleTasks[0].Id // ID of "Initial Task 1"
	// task2ID := sampleTasks[1].Id // ID of "Initial Task 2"

	t.Run("should find an existing task", func(t *testing.T) {
		// Act
		task, err := testManager.GetTaskById(ctx, task1ID)

		// Assert
		assert.NoError(t, err, "Should not return an error for an existing ID")
		assert.NotNil(t, task, "Returned task should not be nil")
		assert.Equal(t, task1ID, task.Id, "Task ID should match the requested ID")
		assert.Equal(t, "Initial Task 1", task.Title)
	})

	t.Run("should return an error for a non-existent task", func(t *testing.T) {
		// Arrange: Generate a valid but non-existent ObjectID
		nonExistentID := primitive.NewObjectID()

		// Act
		task, err := testManager.GetTaskById(ctx, nonExistentID)

		// Assert
		assert.True(t, errors.Is(err, ErrTaskNotFound), "Error should be ErrTaskNotFound")
		assert.Contains(t, err.Error(), nonExistentID.Hex(), "Error message should contain the ID") // Specific error message
		assert.True(t, task.Id.IsZero(), "Returned task should have zero ID when not found")        // task struct will be zero value
	})
}

// TestUpdateTask tests the UpdateTask method.
func TestUpdateTask(t *testing.T) {
	setupTestCollection(t) // Clean before this test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Arrange: Add a task to update
	sampleTasks, err := addSampleTasks(ctx, t)
	assert.NoError(t, err)
	assert.Len(t, sampleTasks, 2)

	sort.Slice(sampleTasks, func(i, j int) bool {
		return sampleTasks[i].Title < sampleTasks[j].Title
	})
	taskToUpdateID := sampleTasks[0].Id // ID of "Initial Task 1"

	t.Run("should update an existing task", func(t *testing.T) {
		updatePayload := models.Task{
			Title:       "Updated Title",
			Description: "This description has been updated.",
			Status:      models.Done,
			DueDate:     time.Now().Add(48 * time.Hour),
		}

		// Act
		updatedTask, err := testManager.UpdateTask(ctx, taskToUpdateID, updatePayload)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedTask)
		assert.Equal(t, taskToUpdateID, updatedTask.Id, "ID should remain the same")
		assert.Equal(t, "Updated Title", updatedTask.Title, "Title should be updated")
		assert.Equal(t, models.Done, updatedTask.Status, "Status should be updated")
		assert.WithinDuration(t, updatePayload.DueDate, updatedTask.DueDate, time.Second) // Compare dates loosely

		// Verify the change was persisted
		retrievedTask, _ := testManager.GetTaskById(ctx, taskToUpdateID)
		assert.Equal(t, "Updated Title", retrievedTask.Title)
		assert.Equal(t, models.Done, retrievedTask.Status)
	})

	t.Run("should return error when updating a non-existent task", func(t *testing.T) {
		// Arrange: Generate a non-existent ObjectID
		nonExistentID := primitive.NewObjectID()

		// Act
		_, err := testManager.UpdateTask(ctx, nonExistentID, models.Task{Title: "Won't work"})

		// Assert
		assert.True(t, errors.Is(err, ErrTaskNotFound), "Error should be ErrTaskNotFound")
		assert.Contains(t, err.Error(), nonExistentID.Hex(), "Error message should contain the ID")
	})
}

// TestDeleteTask tests the DeleteTask method.
func TestDeleteTask(t *testing.T) {
	setupTestCollection(t) // Clean before this test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Arrange: Add tasks, one of which will be deleted
	sampleTasks, err := addSampleTasks(ctx, t)
	assert.NoError(t, err)
	assert.Len(t, sampleTasks, 2)

	sort.Slice(sampleTasks, func(i, j int) bool {
		return sampleTasks[i].Title < sampleTasks[j].Title
	})
	taskToDeleteID := sampleTasks[0].Id // ID of "Initial Task 1"

	t.Run("should delete an existing task", func(t *testing.T) {
		// Act
		err := testManager.DeleteTask(ctx, taskToDeleteID)

		// Assert
		assert.NoError(t, err, "Deleting an existing task should not produce an error")

		// Verify it's gone
		_, errAfterDelete := testManager.GetTaskById(ctx, taskToDeleteID)
		assert.True(t, errors.Is(errAfterDelete, ErrTaskNotFound), "Error should be ErrTaskNotFound")
		assert.Contains(t, errAfterDelete.Error(), "not found", "Error message should indicate task not found")

		// Verify other tasks are unaffected
		allTasks, err := testManager.GetTasks(ctx)
		assert.NoError(t, err)
		assert.Len(t, allTasks, 1, "The total number of tasks should be 1 after deletion")
		assert.Equal(t, sampleTasks[1].Title, allTasks[0].Title, "The remaining task should be the other sample task")
	})

	t.Run("should return error when deleting a non-existent task", func(t *testing.T) {
		// Arrange: Generate a non-existent ObjectID
		nonExistentID := primitive.NewObjectID()

		// Act
		err := testManager.DeleteTask(ctx, nonExistentID)

		// Assert
		assert.True(t, errors.Is(err, ErrTaskNotFound), "Error should be ErrTaskNotFound")
		assert.Contains(t, err.Error(), nonExistentID.Hex(), "Error message should contain the ID")
		assert.Contains(t, err.Error(), "not found", "Error message should indicate task not found")
	})
}

package data

import (
	"A2SV_ProjectPhase/Task4/TaskManager/models"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// setupService is a helper function to create a new service and populate it with some initial data.
// This helps avoid repetitive code in each test case.
func setupService() TaskManager {
	service := NewTaskManager()
	service.AddTask(models.Task{
		Title:       "Initial Task 1",
		Description: "First task for setup",
		Status:      models.Pending,
		DueDate:     time.Now(),
	})
	service.AddTask(models.Task{
		Title:       "Initial Task 2",
		Description: "Second task for setup",
		Status:      models.InProgress,
		DueDate:     time.Now(),
	})
	return service
}

// TestAddTask tests the AddTask method.
func TestAddTask(t *testing.T) {
	// Arrange
	service := NewTaskManager() // Start with a fresh service
	taskToAdd := models.Task{
		Title:       "A New Task",
		Description: "Details of the new task",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      models.Pending,
	}

	// Act
	createdTask := service.AddTask(taskToAdd)

	// Assert
	assert.NotNil(t, createdTask, "Created task should not be nil")
	assert.Equal(t, 1, createdTask.Id, "ID of the first task should be 1")
	assert.Equal(t, "A New Task", createdTask.Title, "Task title should match the input")
	assert.Equal(t, models.Pending, createdTask.Status, "Task status should match the input")

	// Verify it was actually stored by trying to retrieve it
	retrievedTask, err := service.GetTaskById(1)
	assert.NoError(t, err)
	assert.Equal(t, createdTask, retrievedTask, "The created task and retrieved task should be the same")
}

// TestGetTasks tests the GetTasks method.
func TestGetTasks(t *testing.T) {
	// Arrange
	service := setupService() // Use helper to get a service with 2 tasks

	// Act
	allTasks := service.GetTasks()

	// Sort the slice by ID to ensure a consistent order for testing since we used unordered map.
	sort.Slice(allTasks, func(i, j int) bool {
		return allTasks[i].Id < allTasks[j].Id
	})

	// Assert
	assert.NotNil(t, allTasks, "The returned slice of tasks should not be nil")
	assert.Len(t, allTasks, 2, "Should return a slice containing 2 tasks")
	assert.Equal(t, "Initial Task 1", allTasks[0].Title)
	assert.Equal(t, "Initial Task 2", allTasks[1].Title)
}

// TestGetTaskById tests the GetTaskById method for success and failure cases.
func TestGetTaskById(t *testing.T) {
	// Arrange
	service := setupService()

	t.Run("should find an existing task", func(t *testing.T) {
		// Act
		task, err := service.GetTaskById(1)

		// Assert
		assert.NoError(t, err, "Should not return an error for an existing ID")
		assert.NotNil(t, task, "Returned task should not be nil")
		assert.Equal(t, 1, task.Id, "Task ID should be 1")
		assert.Equal(t, "Initial Task 1", task.Title)
	})

	t.Run("should return an error for a non-existent task", func(t *testing.T) {
		// Act
		task, err := service.GetTaskById(99) // An ID that doesn't exist

		// Assert
		assert.Error(t, err, "Should return an error for a non-existent ID")
		assert.Nil(t, task, "Returned task should be nil when an error occurs")
	})
}

// TestUpdateTask tests the UpdateTask method.
func TestUpdateTask(t *testing.T) {
	// Arrange
	service := setupService()

	t.Run("should update an existing task", func(t *testing.T) {
		updatePayload := models.Task{
			Title:       "Updated Title",
			Description: "This description has been updated.",
			Status:      models.Done,
			DueDate:     time.Now().Add(48 * time.Hour),
		}

		// Act
		updatedTask, err := service.UpdateTask(1, updatePayload)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedTask)
		assert.Equal(t, 1, updatedTask.Id, "ID should remain the same")
		assert.Equal(t, "Updated Title", updatedTask.Title, "Title should be updated")
		assert.Equal(t, models.Done, updatedTask.Status, "Status should be updated")

		// Verify the change was persisted
		retrievedTask, _ := service.GetTaskById(1)
		assert.Equal(t, "Updated Title", retrievedTask.Title)
	})

	t.Run("should return error when updating a non-existent task", func(t *testing.T) {
		// Act
		_, err := service.UpdateTask(99, models.Task{Title: "Won't work"})

		// Assert
		assert.Error(t, err, "Should return an error when trying to update a non-existent task")
	})
}

// TestDeleteTask tests the DeleteTask method.
func TestDeleteTask(t *testing.T) {
	// Arrange
	service := setupService()

	t.Run("should delete an existing task", func(t *testing.T) {
		// Act
		err := service.DeleteTask(1)

		// Assert
		assert.NoError(t, err, "Deleting an existing task should not produce an error")

		// Verify it's gone
		_, errAfterDelete := service.GetTaskById(1)
		assert.Error(t, errAfterDelete, "Should return an error when trying to get a deleted task")

		// Verify other tasks are unaffected
		allTasks := service.GetTasks()
		assert.Len(t, allTasks, 1, "The total number of tasks should be 1 after deletion")
	})

	t.Run("should return error when deleting a non-existent task", func(t *testing.T) {
		// Act
		err := service.DeleteTask(99)

		// Assert
		assert.Error(t, err, "Should return an error when trying to delete a non-existent task")
	})
}

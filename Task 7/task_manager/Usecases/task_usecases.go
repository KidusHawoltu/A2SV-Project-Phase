package usecases

import (
	domain "A2SV_ProjectPhase/Task7/TaskManager/Domain"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUseCase struct {
	taskRepo domain.TaskRepository
}

func NewTaskUseCase(taskRepo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{
		taskRepo: taskRepo,
	}
}

func (uc *TaskUseCase) CreateTask(c context.Context, title, description string, dueDate time.Time, status domain.TaskStatus) (*domain.Task, error) {
	newTask, err := domain.NewTask(title, description, dueDate, status)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create task entity: %s", domain.ErrValidationFailed, err.Error())
	}

	// 2. Persist the task via repository
	savedTask, err := uc.taskRepo.CreateTask(c, newTask)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to save task: %w", err)
	}

	return savedTask, nil
}

// GetTaskByID handles fetching a single task by its ID.
func (uc *TaskUseCase) GetTaskByID(c context.Context, taskID string) (*domain.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid task ID format", domain.ErrValidationFailed)
	}

	task, err := uc.taskRepo.GetTaskById(c, objectID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, domain.ErrTaskNotFound // Propagate task not found
		}
		return nil, fmt.Errorf("usecase: failed to get task by ID: %w", err) // Unexpected repo error
	}
	return task, nil
}

// GetAllTasks handles fetching all tasks.
func (uc *TaskUseCase) GetAllTasks(c context.Context) ([]*domain.Task, error) {
	tasks, err := uc.taskRepo.GetAllTasks(c)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to get all tasks: %w", err)
	}
	return tasks, nil
}

// It takes optional fields using pointers, allowing partial updates.
func (uc *TaskUseCase) UpdateTask(c context.Context, taskID string, title, description *string, dueDate *time.Time, status *domain.TaskStatus) (*domain.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid task ID format", domain.ErrValidationFailed)
	}

	// 1. Fetch existing task to ensure it exists
	existingTask, err := uc.taskRepo.GetTaskById(c, objectID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("usecase: failed to retrieve existing task for update: %w", err)
	}

	// 2. Apply updates to the existing domain entity based on provided non-nil pointers
	if title != nil {
		if *title == "" { // title cannot be empty on update
			return nil, fmt.Errorf("%w: task title cannot be empty on update", domain.ErrValidationFailed)
		}
		existingTask.Title = *title
	}
	if description != nil {
		existingTask.Description = *description
	}
	if dueDate != nil {
		if dueDate.Before(time.Now().Truncate(24 * time.Hour)) { // updated due date not in past
			return nil, fmt.Errorf("%w: updated due date cannot be in the past", domain.ErrValidationFailed)
		}
		existingTask.DueDate = *dueDate
	}
	if status != nil {
		if !status.IsValid() {
			return nil, fmt.Errorf("%w: invalid task status for update", domain.ErrValidationFailed)
		}
		// Cannot change status from Done to anything else
		if existingTask.Status == domain.Done && *status != domain.Done {
			return nil, fmt.Errorf("%w: cannot change status of a completed task", domain.ErrValidationFailed)
		}
		existingTask.Status = *status
	}

	// 3. Persist the updated task
	updatedTaskResult, err := uc.taskRepo.UpdateTask(c, objectID, existingTask)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to update task: %w", err)
	}

	return updatedTaskResult, nil
}

// DeleteTask handles deleting a task by its ID.
func (uc *TaskUseCase) DeleteTask(c context.Context, taskID string) error {
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return fmt.Errorf("%w: invalid task ID format", domain.ErrValidationFailed)
	}

	err = uc.taskRepo.DeleteTask(c, objectID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return domain.ErrTaskNotFound // Propagate task not found
		}
		return fmt.Errorf("usecase: failed to delete task: %w", err)
	}
	return nil
}

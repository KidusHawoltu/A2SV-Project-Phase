package data

import (
	"A2SV_ProjectPhase/Task6/TaskManager/models"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson" // Import bson package
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" // Import options for FindOneAndUpdate
)

type TaskCollection struct {
	collection *mongo.Collection
}

// Define a sentinel error for "task not found"
var ErrTaskNotFound = errors.New("task not found")

// TaskManager interface now includes context.Context for all operations.
type TaskManager interface {
	GetTasks(ctx context.Context) ([]models.Task, error)
	GetTaskById(ctx context.Context, id primitive.ObjectID) (models.Task, error)
	UpdateTask(ctx context.Context, id primitive.ObjectID, updatedTask models.Task) (models.Task, error)
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
	AddTask(ctx context.Context, task models.Task) (models.Task, error)
}

func NewTaskManager(tc *mongo.Collection) TaskManager {
	return &TaskCollection{
		collection: tc,
	}
}

func (tc *TaskCollection) GetTasks(ctx context.Context) ([]models.Task, error) {
	tasks := []models.Task{}
	// Pass context to Find operation
	curr, err := tc.collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}
	defer curr.Close(ctx) // Always close the cursor

	if err = curr.All(ctx, &tasks); err != nil { // Efficiently decode all documents
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}
	return tasks, nil
}

func (tc *TaskCollection) GetTaskById(ctx context.Context, id primitive.ObjectID) (models.Task, error) {
	filter := bson.M{"_id": id} // Use bson.M for simpler filters
	var task models.Task
	// Pass context to FindOne operation
	err := tc.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, fmt.Errorf("%w: task with id '%s'", ErrTaskNotFound, id.Hex())
		}
		return models.Task{}, fmt.Errorf("failed to retrieve task: %w", err)
	}
	return task, nil
}

func (tc *TaskCollection) UpdateTask(ctx context.Context, id primitive.ObjectID, updatedTask models.Task) (models.Task, error) {
	// Use $set operator to update only provided fields
	// The `bson` tags in models.Task will ensure correct field names in MongoDB
	updateDoc := bson.M{"$set": bson.M{
		"title":       updatedTask.Title,
		"description": updatedTask.Description,
		"duedate":     updatedTask.DueDate,
		"status":      updatedTask.Status,
	}}

	filter := bson.M{"_id": id}
	var result models.Task

	// FindOneAndUpdate returns the document after the update
	err := tc.collection.FindOneAndUpdate(
		ctx,
		filter,
		updateDoc,
		options.FindOneAndUpdate().SetReturnDocument(options.After), // Get the document AFTER the update
	).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, fmt.Errorf("%w: task with id '%s'", ErrTaskNotFound, id.Hex())
		}
		return models.Task{}, fmt.Errorf("failed to update task: %w", err)
	}
	return result, nil
}

func (tc *TaskCollection) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	// Pass context to DeleteOne operation
	res, err := tc.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("%w: task with id '%s'", ErrTaskNotFound, id.Hex())
	}
	return nil
}

func (tc *TaskCollection) AddTask(ctx context.Context, task models.Task) (models.Task, error) {
	// MongoDB automatically generates _id if it's not provided or is primitive.NilObjectID
	task.Id = primitive.NilObjectID // Ensure MongoDB generates a new ID

	// Pass context to InsertOne operation
	inserted, err := tc.collection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to add task: %w", err)
	}

	// Cast the inserted ID to primitive.ObjectID
	id, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.Task{}, fmt.Errorf("created ID is not a valid ObjectID: %v", inserted.InsertedID)
	}
	task.Id = id // Set the generated ID back to the task object
	return task, nil
}

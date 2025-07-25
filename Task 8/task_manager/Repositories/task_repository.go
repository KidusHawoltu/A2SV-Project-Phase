package repositories

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure TaskRepo implements the domain.TaskRepository interface
var _ domain.TaskRepository = (*TaskRepo)(nil)

type TaskRepo struct {
	collection *mongo.Collection
}

func NewMongoDBTaskRepository(col *mongo.Collection) *TaskRepo {
	return &TaskRepo{
		collection: col,
	}
}

func (tr *TaskRepo) CreateTask(c context.Context, task *domain.Task) (*domain.Task, error) {
	result, err := tr.collection.InsertOne(c, task)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to insert task: %w", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("repository: inserted ID is not of type ObjectID: %T", result.InsertedID)
	}
	task.Id = insertedID

	return task, nil
}

func (tr *TaskRepo) GetTaskById(c context.Context, id primitive.ObjectID) (*domain.Task, error) {
	var task domain.Task
	filter := bson.M{"_id": id}
	err := tr.collection.FindOne(c, filter).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("repository: failed to find task by ID '%s': %w", id.Hex(), err)
	}
	return &task, nil
}

func (tr *TaskRepo) GetAllTasks(c context.Context) ([]*domain.Task, error) {
	cursor, err := tr.collection.Find(c, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("repository: failed to retrieve tasks cursor: %w", err)
	}
	defer cursor.Close(c)

	tasks := []*domain.Task{}
	if err = cursor.All(c, &tasks); err != nil {
		return nil, fmt.Errorf("repository: failed to decode tasks from cursor: %w", err)
	}
	return tasks, nil
}

func (tr *TaskRepo) UpdateTask(c context.Context, id primitive.ObjectID, updatedTask *domain.Task) (*domain.Task, error) {
	updateDoc := bson.M{"$set": bson.M{
		"title":       updatedTask.Title,
		"description": updatedTask.Description,
		"duedate":     updatedTask.DueDate,
		"status":      updatedTask.Status,
	}}

	filter := bson.M{"_id": id}
	var result domain.Task

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // Get the document AFTER the update
	err := tr.collection.FindOneAndUpdate(c, filter, updateDoc, opts).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("repository: failed to update task by ID '%s': %w", id.Hex(), err)
	}
	return &result, nil
}

func (tr *TaskRepo) DeleteTask(c context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	res, err := tr.collection.DeleteOne(c, filter)
	if err != nil {
		return fmt.Errorf("repository: failed to delete task by ID '%s': %w", id.Hex(), err)
	}
	if res.DeletedCount == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

package repositories

import (
	domain "A2SV_ProjectPhase/Task7/TaskManager/Domain"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewMongoDBUserRepository(col *mongo.Collection) *UserRepo {
	return &UserRepo{
		collection: col,
	}
}

var _ domain.UserRepository = (*UserRepo)(nil)

func (ur *UserRepo) CreateUser(c context.Context, user *domain.User) (*domain.User, error) {
	result, err := ur.collection.InsertOne(c, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrUsernameTaken
		}
		return nil, fmt.Errorf("repository: failed to insert user: %w", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("repository: inserted ID is not of type ObjectID: %T", result.InsertedID)
	}
	user.Id = insertedID

	return user, nil
}

func (ur *UserRepo) GetUserByUsername(c context.Context, username string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{"username": username}
	err := ur.collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("repository: failed to find user by username '%s': %w", username, err)
	}
	return &user, nil
}

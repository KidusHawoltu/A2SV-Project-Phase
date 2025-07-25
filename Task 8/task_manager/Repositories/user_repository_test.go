package repositories_test

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	"A2SV_ProjectPhase/Task8/TaskManager/Repositories"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//===========================================================================
// UserRepo Integration Test Suite
//===========================================================================

type UserRepoSuite struct {
	suite.Suite
	db   *mongo.Database
	coll *mongo.Collection
	repo domain.UserRepository
}

// TestUserRepoSuite is the entry point for the test suite
func TestUserRepoSuite(t *testing.T) {
	if testMongoClient == nil {
		t.Skip("Skipping integration tests: MongoDB connection not available.")
	}
	suite.Run(t, new(UserRepoSuite))
}

// SetupSuite runs once for the entire suite.
func (s *UserRepoSuite) SetupSuite() {
	s.db = testMongoClient.Database("test_learning_phase")
	s.coll = s.db.Collection("user8")
}

// SetupTest runs before EACH test method.
func (s *UserRepoSuite) SetupTest() {
	// 1. Clean the collection to ensure test isolation.
	_, err := s.coll.DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err, "Failed to clean user collection before test")

	// 2. Ensure the unique index on 'username' exists for testing duplicate errors.
	// This is idempotent; if the index already exists, MongoDB handles it gracefully.
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = s.coll.Indexes().CreateOne(context.Background(), indexModel)
	s.Require().NoError(err, "Failed to create unique index on username")

	// 3. Create a new repository instance for the test.
	s.repo = repositories.NewMongoDBUserRepository(s.coll)
}

// TestCreateUser tests the user creation repository logic.
func (s *UserRepoSuite) TestCreateUser() {
	s.Run("Success", func() {
		user := &domain.User{Username: "testuser", PasswordHash: "hash"}
		createdUser, err := s.repo.CreateUser(context.Background(), user)

		s.Require().NoError(err)
		s.Require().NotNil(createdUser)
		s.False(createdUser.Id.IsZero())
		s.Equal("testuser", createdUser.Username)
	})

	s.Run("Failure - Duplicate Username", func() {
		// First user (seeded)
		user1 := &domain.User{Username: "duplicate", PasswordHash: "hash1"}
		_, err := s.repo.CreateUser(context.Background(), user1)
		s.Require().NoError(err, "Seeding the first user should succeed")

		// Attempt to create a second user with the same username
		user2 := &domain.User{Username: "duplicate", PasswordHash: "hash2"}
		_, err = s.repo.CreateUser(context.Background(), user2)

		// Assert that we get the correct domain-level error
		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrUsernameTaken)
	})
}

// TestGetUserByUsername tests finding a user by their username.
func (s *UserRepoSuite) TestGetUserByUsername() {
	s.Run("Success - User Found", func() {
		// Setup: Seed the database with a user to find.
		userToFind := &domain.User{
			Id:       primitive.NewObjectID(),
			Username: "findme",
		}
		// We can use the repo to seed, or insert directly into the collection
		_, err := s.coll.InsertOne(context.Background(), userToFind)
		s.Require().NoError(err, "Failed to seed database for test")

		// Execution
		foundUser, err := s.repo.GetUserByUsername(context.Background(), "findme")

		// Assertion
		s.Require().NoError(err)
		s.Require().NotNil(foundUser)
		s.Equal(userToFind.Id, foundUser.Id)
	})

	s.Run("Failure - User Not Found", func() {
		// No seeding needed, collection is clean from SetupTest.

		// Execution
		_, err := s.repo.GetUserByUsername(context.Background(), "nonexistentuser")

		// Assertion
		s.Require().Error(err)
		s.ErrorIs(err, domain.ErrUserNotFound)
	})
}

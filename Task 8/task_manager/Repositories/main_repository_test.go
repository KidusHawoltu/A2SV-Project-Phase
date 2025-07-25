package repositories_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testMongoClient *mongo.Client

// TestMain is a special function that Go's test runner executes before any other
// tests in this package. It's the perfect place for global setup and teardown.
func TestMain(m *testing.M) {
	// 1. Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		// if not in parent folder, try to load from current folder
		err = godotenv.Load()
		if err != nil {
			log.Printf("Warning: No .env file found or failed to load: %v", err)
		}
	}

	// 2. Get the MongoDB URI, preferring the test-specific URI
	mongoURI := os.Getenv("MONGO_TEST_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_URI")
	}
	if mongoURI == "" {
		log.Fatal("FATAL: Neither MONGO_TEST_URI nor MONGO_URI is set. Cannot run integration tests.")
	}

	// 3. Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to MongoDB for testing: %v", err)
	}

	// 4. Ping to confirm connection
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("FATAL: Failed to ping MongoDB: %v", err)
	}

	log.Println("Test MongoDB connection established.")
	testMongoClient = client // Store the client in a global variable for this package

	// 5. Run all the tests
	exitCode := m.Run()

	// 6. Teardown: Clean up the test database and disconnect
	log.Println("Dropping test database and disconnecting...")
	db := testMongoClient.Database("test_learning_phase")
	if err := db.Drop(context.Background()); err != nil {
		log.Printf("WARNING: Failed to drop test database: %v", err)
	}
	if err := testMongoClient.Disconnect(context.Background()); err != nil {
		log.Printf("WARNING: Failed to disconnect test MongoDB client: %v", err)
	}

	// 7. Exit with the tests' exit code
	os.Exit(exitCode)
}

package main

import (
	"A2SV_ProjectPhase/Task5/TaskManager/controllers"
	"A2SV_ProjectPhase/Task5/TaskManager/data"
	"A2SV_ProjectPhase/Task5/TaskManager/router"
	"context"
	"log" // Changed from time for consistent logging of errors
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Load .env file
	// godotenv.Load() attempts to load .env from the current directory.
	// It's good practice to log if it fails, but not necessarily fatal,
	// as env vars might be set directly in the deployment environment.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env. Assuming environment variables are set directly.")
	}

	// Get MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set. Please set it or ensure your .env file is correctly configured.")
	}

	// Create a context with a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure the context is cancelled when main exits

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	log.Println("Attempting to connect to MongoDB...")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts) // Use the context here
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		log.Println("Closing MongoDB connection...")
		if err = client.Disconnect(ctx); err != nil { // Use the context for disconnect as well
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
		log.Println("MongoDB connection closed.")
	}()

	// Send a ping to confirm a successful connection
	log.Println("Pinging MongoDB deployment...")
	if err := client.Ping(ctx, readpref.Primary()); err != nil { // Use the context here
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB!")

	// Get a handle to the database and collection
	databaseName := "learning_phase" // Can be configured
	collectionName := "tasks"        // Can be configured
	taskCollection := client.Database(databaseName).Collection(collectionName)

	manager := data.NewTaskManager(taskCollection)
	controller := controllers.NewTaskController(manager)
	router := router.NewRouter(controller)
	router.Run(":8080")
}

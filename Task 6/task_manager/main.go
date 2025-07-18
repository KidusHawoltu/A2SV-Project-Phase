package main

import (
	"A2SV_ProjectPhase/Task6/TaskManager/controllers"
	"A2SV_ProjectPhase/Task6/TaskManager/data"
	"A2SV_ProjectPhase/Task6/TaskManager/models"
	"A2SV_ProjectPhase/Task6/TaskManager/router"
	"context"
	"errors"
	"fmt"
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
	taskCollectionName := "task6"    // Can be configured
	userCollectionName := "user6"    // Can be configured
	taskCollection := client.Database(databaseName).Collection(taskCollectionName)
	userCollection := client.Database(databaseName).Collection(userCollectionName)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default secret"
		fmt.Println("WARNING: JWT_SECRET environment variable not set. Using default secret.")
	}

	taskManager := data.NewTaskManager(taskCollection)
	userManager := data.NewUserManager(userCollection, jwtSecret)

	// Automatically add admin from enviroment variable if it doesn't already exsist in database
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
		log.Println("WARNING: ADMIN_USERNAME not set in .env or environment. Using default value.")
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "password"
		log.Println("WARNING: ADMIN_PASSWORD not set in .env or environment. Using default value.")
	}

	log.Printf("Checking for default admin user '%s'...", adminUsername)
	adminUser, err := userManager.GetByUsername(ctx, adminUsername) // Use the context from main
	if err != nil {
		if errors.Is(err, data.ErrUserNotFound) {
			log.Printf("Default admin user '%s' not found. Attempting to create...", adminUsername)
			newAdmin := models.User{
				Username: adminUsername,
				Password: adminPassword,
				Role:     models.RoleAdmin, // Ensure this user is an admin
			}
			_, regErr := userManager.RegisterUser(ctx, newAdmin) // Use the context from main
			if regErr != nil {
				log.Fatalf("Failed to create default admin user '%s': %v", adminUsername, regErr)
			}
			log.Printf("Successfully created default admin user '%s'", adminUsername)
		} else {
			log.Fatalf("Failed to check for default admin user '%s': %v", adminUsername, err)
		}
	} else {
		log.Printf("Default admin user '%s' already exists (ID: %s).", adminUsername, adminUser.Id.Hex())
	}

	controller := controllers.NewTaskController(taskManager, userManager)
	router := router.NewRouter(controller)
	router.Run(":8080")
}

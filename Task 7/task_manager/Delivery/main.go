package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"A2SV_ProjectPhase/Task7/TaskManager/Delivery/controllers"
	"A2SV_ProjectPhase/Task7/TaskManager/Delivery/routers"
	domain "A2SV_ProjectPhase/Task7/TaskManager/Domain"
	infrastructure "A2SV_ProjectPhase/Task7/TaskManager/Infrastructure"
	repositories "A2SV_ProjectPhase/Task7/TaskManager/Repositories"
	usecases "A2SV_ProjectPhase/Task7/TaskManager/Usecases"
	"errors"
)

func main() {
	// --- 0. Load Configuration / Environment Variables ---
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Warning: No .env file found or failed to load: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalf("Fatal: MONGO_URI environment variable not set.")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")
	if jwtSecretKey == "" {
		jwtSecretKey = "default_secret"
		log.Println("WARNING: JWT_SECRET environment variable not set. Using default secret.")
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "adminpassword"
	}

	// --- 1. Initialize External Resources (MongoDB connection) ---
	clientOptions := options.Client().ApplyURI(mongoURI)
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Fatal: Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Warning: Failed to disconnect from MongoDB: %v", err)
		}
	}()
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Fatal: Failed to ping MongoDB: %v", err)
	}
	log.Println("MongoDB connection established.")

	db := mongoClient.Database("learning_phase")

	// --- 2. Instantiate Concrete Infrastructure Services (Needed for bootstrapping too) ---
	// We need passwordService here directly for hashing admin password
	passwordService := infrastructure.NewBcryptPasswordService(bcrypt.DefaultCost)
	jwtService := infrastructure.NewJwtService(jwtSecretKey) // Still needed for JWTs later
	log.Println("Infrastructure services initialized.")

	// --- 3. Instantiate Concrete Repository Implementations (Needed for bootstrapping) ---
	userCollection := db.Collection("user7")
	taskCollection := db.Collection("task7")

	userRepo := repositories.NewMongoDBUserRepository(userCollection) // Needed directly for admin check/create
	taskRepo := repositories.NewMongoDBTaskRepository(taskCollection)
	log.Println("Repositories initialized.")

	// --- 4. Implement Default Admin User Bootstrapping (Directly using Repo and PasswordService) ---
	log.Printf("Checking for default admin user '%s'...", adminUsername)
	adminUserCtx := context.Background() // Use background context for startup task

	// Try to find the admin user
	_, err = userRepo.GetUserByUsername(adminUserCtx, adminUsername)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Printf("Default admin user '%s' not found. Attempting to create...", adminUsername)

			hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
			if hashErr != nil {
				log.Fatalf("Fatal: Failed to hash admin password for bootstrapping: %v", hashErr)
			}
			hashedPassword := string(hashedPasswordBytes)
			newAdminUser := &domain.User{
				Id:           primitive.NilObjectID, // Let MongoDB generate
				Username:     adminUsername,
				PasswordHash: hashedPassword,
				Role:         domain.RoleAdmin, // Explicitly set role to Admin
			}

			// Save the admin user directly via the repository
			createdAdmin, createErr := userRepo.CreateUser(adminUserCtx, newAdminUser)
			if createErr != nil {
				log.Fatalf("Fatal: Failed to create default admin user '%s': %v", adminUsername, createErr)
			}
			log.Printf("Successfully created default admin user '%s' (ID: %s) with role '%s'.\n",
				createdAdmin.Username, createdAdmin.Id.Hex(), createdAdmin.Role)

		} else {
			log.Fatalf("Fatal: Unexpected error during admin user lookup: %v", err)
		}
	} else {
		log.Printf("Default admin user '%s' found. Proceeding.\n", adminUsername)
	}
	log.Println("Admin bootstrapping complete.")

	// --- 5. Instantiate Usecases (Injecting Repositories and Infrastructure Services as Interfaces) ---
	// Note: userUsecase is initialized *after* bootstrapping
	userUsecase := usecases.NewUserUseCase(userRepo, jwtService, passwordService)
	taskUsecase := usecases.NewTaskUseCase(taskRepo)
	log.Println("Usecases initialized.")

	// --- 6. Instantiate Delivery Controllers (Injecting Usecases) ---
	userController := controllers.NewUserController(userUsecase)
	taskController := controllers.NewTaskController(taskUsecase)
	authMiddleware := infrastructure.NewAuthMiddleware(jwtService)
	log.Println("Controllers and middleware initialized.")

	// --- 7. Set Up Delivery Routers ---
	router := gin.Default()
	{
		routers.SetupUserRouters(router, userController)
		routers.SetupTaskRoutes(router, taskController, authMiddleware)
	}

	log.Println("All Routers configured.")

	// --- 8. Start the HTTP Server ---
	log.Printf("Server starting on :8080")
	log.Fatal(router.Run(":8080"))
}

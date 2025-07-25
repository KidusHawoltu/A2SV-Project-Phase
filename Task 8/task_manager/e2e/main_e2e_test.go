package e2e_test

import (
	"A2SV_ProjectPhase/Task8/TaskManager/Delivery/controllers"
	"A2SV_ProjectPhase/Task8/TaskManager/Delivery/routers"
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	infrastructure "A2SV_ProjectPhase/Task8/TaskManager/Infrastructure"
	repositories "A2SV_ProjectPhase/Task8/TaskManager/Repositories"
	usecases "A2SV_ProjectPhase/Task8/TaskManager/Usecases"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// --- Global variables for the test package ---
var (
	testMongoClient *mongo.Client
	mongoURI        string
	jwtSecret       string
	testDBName      = "test_learning_phase"
	userCol         = "user8"
	taskCol         = "task8"
)

// TestMain controls the entire lifecycle for the e2e test package.
func TestMain(m *testing.M) {
	// Load .env file from the project root
	err := godotenv.Load("../.env")
	if err != nil {
		// if not in parent folder, try to load from current folder
		err = godotenv.Load()
		if err != nil {
			log.Printf("Warning: No .env file found or failed to load: %v", err)
		}
	}

	// Get config, preferring E2E-specific variables
	mongoURI = os.Getenv("MONGO_TEST_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_URI")
	}
	if mongoURI == "" {
		log.Fatal("FATAL: Neither MONGO_TEST_URI nor MONGO_URI is set. Cannot run E2E tests.")
	}

	jwtSecret = os.Getenv("JWT_TEST_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if jwtSecret == "" {
		jwtSecret = "default_e2e_secret" // A fallback for safety
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to MongoDB for E2E testing: %v", err)
	}
	testMongoClient = client

	// Run all tests
	exitCode := m.Run()

	// Teardown: drop the entire test database
	log.Printf("Dropping E2E test database: %s", testDBName)
	if err := testMongoClient.Database(testDBName).Drop(context.Background()); err != nil {
		log.Printf("WARNING: Failed to drop test database: %v", err)
	}
	testMongoClient.Disconnect(context.Background())

	os.Exit(exitCode)
}

// setupApplication assembles the entire application stack and returns a usable router.
func setupApplication() *gin.Engine {
	// Use the globally connected client
	db := testMongoClient.Database(testDBName)

	// Instantiate all layers with real implementations
	passwordService := infrastructure.NewBcryptPasswordService(bcrypt.DefaultCost)
	jwtService := infrastructure.NewJwtService(jwtSecret)
	userCollection := db.Collection(userCol)
	taskCollection := db.Collection(taskCol)
	userRepo := repositories.NewMongoDBUserRepository(userCollection)
	taskRepo := repositories.NewMongoDBTaskRepository(taskCollection)
	userUsecase := usecases.NewUserUseCase(userRepo, jwtService, passwordService)
	taskUsecase := usecases.NewTaskUseCase(taskRepo)
	userController := controllers.NewUserController(userUsecase)
	taskController := controllers.NewTaskController(taskUsecase)
	authMiddleware := infrastructure.NewAuthMiddleware(jwtService)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routers.SetupUserRouters(router, userController)
	routers.SetupTaskRoutes(router, taskController, authMiddleware)

	return router
}

//===========================================================================
// Base E2E Test Suite (handles server and DB cleanup)
//===========================================================================

type E2ETestSuite struct {
	suite.Suite
	Router *gin.Engine
	Server *httptest.Server
	DB     *mongo.Database
}

func (s *E2ETestSuite) SetupSuite() {
	s.Router = setupApplication()
	s.Server = httptest.NewServer(s.Router)
	s.DB = testMongoClient.Database(testDBName)
}

func (s *E2ETestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *E2ETestSuite) SetupTest() {
	// Clean all collections before each test method runs
	collections := []string{userCol, taskCol}
	for _, coll := range collections {
		_, err := s.DB.Collection(coll).DeleteMany(context.Background(), bson.D{})
		s.Require().NoError(err)
	}
}

// Helper to make requests to the test server
func (s *E2ETestSuite) makeRequest(method, path, token string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, s.Server.URL+path, body)
	s.Require().NoError(err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	return resp
}

//===========================================================================
// User Endpoints E2E Test Suite
//===========================================================================

type UserE2ETestSuite struct {
	E2ETestSuite // Embed the base suite
}

func TestUserE2E(t *testing.T) {
	suite.Run(t, new(UserE2ETestSuite))
}

func (s *UserE2ETestSuite) TestRegisterAndLogin() {
	// --- 1. Successful Registration ---
	regBody := bytes.NewBufferString(`{"username": "e2e_user", "password": "e2e_password"}`)
	resp := s.makeRequest(http.MethodPost, "/user/register", "", regBody)
	s.Equal(http.StatusCreated, resp.StatusCode)

	// --- 2. Attempt to register the same user again ---
	regBody2 := bytes.NewBufferString(`{"username": "e2e_user", "password": "e2e_password"}`)
	resp2 := s.makeRequest(http.MethodPost, "/user/register", "", regBody2)
	s.Equal(http.StatusConflict, resp2.StatusCode)

	// --- 3. Successful Login ---
	loginBody := bytes.NewBufferString(`{"username": "e2e_user", "password": "e2e_password"}`)
	resp3 := s.makeRequest(http.MethodPost, "/user/login", "", loginBody)
	s.Equal(http.StatusOK, resp3.StatusCode)

	var loginResp map[string]string
	json.NewDecoder(resp3.Body).Decode(&loginResp)
	s.NotEmpty(loginResp["token"])

	// --- 4. Login with wrong password ---
	loginBody4 := bytes.NewBufferString(`{"username": "e2e_user", "password": "wrong_password"}`)
	resp4 := s.makeRequest(http.MethodPost, "/user/login", "", loginBody4)
	s.Equal(http.StatusUnauthorized, resp4.StatusCode)
}

//===========================================================================
// Task Endpoints E2E Test Suite
//===========================================================================

type TaskE2ETestSuite struct {
	E2ETestSuite
	adminToken string
	userToken  string
}

func TestTaskE2E(t *testing.T) {
	suite.Run(t, new(TaskE2ETestSuite))
}

// SetupSuite for tasks needs to create and log in users to get tokens
func (s *TaskE2ETestSuite) SetupSuite() {
	s.E2ETestSuite.SetupSuite() // Call parent setup first

	// This ensures we start with a clean slate, regardless of what other suites did.
	_, err := s.DB.Collection(userCol).DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err, "Failed to clean users collection before task suite setup")

	// Helper to register and login a user, returning their token
	registerAndLogin := func(username, password string) string {
		// Register
		regBody := bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
		resp := s.makeRequest(http.MethodPost, "/user/register", "", regBody)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)

		// Set admin role directly in DB for the admin user
		if username == "e2e_admin" {
			filter := bson.M{"username": "e2e_admin"}
			update := bson.M{"$set": bson.M{"role": domain.RoleAdmin}}
			_, err := s.DB.Collection(userCol).UpdateOne(context.Background(), filter, update)
			s.Require().NoError(err)
		}

		// Login
		loginBody := bytes.NewBufferString(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
		resp = s.makeRequest(http.MethodPost, "/user/login", "", loginBody)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		var loginResp map[string]string
		json.NewDecoder(resp.Body).Decode(&loginResp)
		return loginResp["token"]
	}

	s.adminToken = registerAndLogin("e2e_admin", "admin_pass")
	s.userToken = registerAndLogin("e2e_user", "user_pass")
}

func (s *TaskE2ETestSuite) TestTaskLifecycleAndAuthorization() {
	var createdTaskID string

	// --- 1. Unauthorized user cannot get tasks ---
	s.Run("Unauthenticated Access Fails", func() {
		resp := s.makeRequest(http.MethodGet, "/tasks", "", nil)
		s.Equal(http.StatusUnauthorized, resp.StatusCode)
	})

	// --- 2. Regular user cannot create a task ---
	s.Run("User Cannot Create Task", func() {
		taskBody := bytes.NewBufferString(`{"title": "user task", "duedate": "2099-01-01T15:04:05Z", "status": "Pending"}`)
		resp := s.makeRequest(http.MethodPost, "/tasks", s.userToken, taskBody)
		s.Equal(http.StatusForbidden, resp.StatusCode)
	})

	// --- 3. Admin can create a task ---
	s.Run("Admin Creates Task", func() {
		taskBody := bytes.NewBufferString(`{"title": "admin task", "description": "desc", "duedate": "2099-01-01T15:04:05Z", "status": "Pending"}`)
		resp := s.makeRequest(http.MethodPost, "/tasks", s.adminToken, taskBody)
		s.Equal(http.StatusCreated, resp.StatusCode)

		var createdTask domain.Task
		json.NewDecoder(resp.Body).Decode(&createdTask)
		s.NotEmpty(createdTask.Id)
		s.Equal("admin task", createdTask.Title)
		createdTaskID = createdTask.Id.Hex()
	})

	// --- 4. Regular user can get all tasks (which includes the admin's task) ---
	s.Run("User Gets All Tasks", func() {
		resp := s.makeRequest(http.MethodGet, "/tasks", s.userToken, nil)
		s.Equal(http.StatusOK, resp.StatusCode)

		var tasks []*domain.Task
		json.NewDecoder(resp.Body).Decode(&tasks)
		s.Len(tasks, 1)
		s.Equal(createdTaskID, tasks[0].Id.Hex())
	})

	// --- 5. Admin can update the task ---
	s.Run("Admin Updates Task", func() {
		updateBody := bytes.NewBufferString(`{"title": "updated admin task", "status": "In progress"}`)
		resp := s.makeRequest(http.MethodPut, "/tasks/"+createdTaskID, s.adminToken, updateBody)
		s.Equal(http.StatusOK, resp.StatusCode)

		var updatedTask domain.Task
		json.NewDecoder(resp.Body).Decode(&updatedTask)
		s.Equal("updated admin task", updatedTask.Title)
		s.Equal(domain.InProgress, updatedTask.Status)
	})

	// --- 6. Admin can delete the task ---
	s.Run("Admin Deletes Task", func() {
		resp := s.makeRequest(http.MethodDelete, "/tasks/"+createdTaskID, s.adminToken, nil)
		s.Equal(http.StatusNoContent, resp.StatusCode)
	})

	// --- 7. Task is no longer available ---
	s.Run("Deleted Task Is Not Found", func() {
		resp := s.makeRequest(http.MethodGet, "/tasks/"+createdTaskID, s.adminToken, nil)
		s.Equal(http.StatusNotFound, resp.StatusCode)
	})
}

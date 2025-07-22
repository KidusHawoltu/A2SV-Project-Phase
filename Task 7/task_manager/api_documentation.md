# Task Manager API

A simple RESTful API for managing tasks, built with Go and the Gin framework, backed by MongoDB. This project provides a complete, thread-safe, and well-structured backend service with robust authentication and authorization using JSON Web Tokens (JWT).

---

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [MongoDB & Application Configuration](#mongodb--application-configuration)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
  - [Running Tests](#running-tests)
- [Architecture Overview](#architecture-overview)
  - [Core Principles](#core-principles)
  - [Layer Breakdown and Responsibilities](#layer-breakdown-and-responsibilities)
  - [Guidelines for Future Development](#guidelines-for-future-development)
- [API Endpoint Reference](#api-endpoint-reference)
  - [Task Model](#task-model)
  - [User & Authentication Models](#user--authentication-models)
  - [Common Error Responses](#common-error-responses)
  - [Endpoints](#endpoints)

---

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You will need the following software installed on your system:
-   **Go**: Version 1.20 or later is recommended.
-   **MongoDB**: Access to a MongoDB instance (e.g., a free tier on MongoDB Atlas or a local Docker instance).
-   **Postman or curl**: For testing the API endpoints manually.

### MongoDB & Application Configuration

This application connects to MongoDB using a connection URI and requires additional configuration for JWT and default admin users. It reads these settings from environment variables.

1.  **Set up a MongoDB Database:**
    *   **MongoDB Atlas (Recommended):** Sign up for a free account at [mongodb.com/atlas](https://www.mongodb.com/atlas).
        *   Create a new project and cluster.
        *   Under "Network Access", add your current public IP address (or `0.0.0.0/0` for development, **not recommended for production**) to the IP Access List. This is crucial for your application to be able to connect.
        *   Under "Database Access", create a new database user with a secure password.
        *   Go to your cluster, click "Connect", select "Drivers", choose "Go" and copy the provided connection string. It will look similar to: `mongodb+srv://<username>:<password>@<cluster-url>/<dbname>?retryWrites=true&w=majority&appName=<appname>`.

2.  **Create a `.env` file:**
    In the root directory of the project (`task-manager/`):

    ```dotenv
    # .env
    MONGO_URI="mongodb+srv://kidushawoltu:avkjCvT5PWJyt@first-cluster.fcjxr0t.mongodb.net/?retryWrites=true&w=majority&appName=First-Cluster"
    
    # Optional: use a separate cluster for testing (recommended)
    # The application tests will use this if set, otherwise fallback to MONGO_URI.
    MONGO_TEST_URI="mongodb+srv://kidushawoltu:avkjCvT5PWJyt@first-cluster.fcjxr0t.mongodb.net/?retryWrites=true&w=majority&appName=First-Cluster-Test"
    
    # JWT Secret Key - CRITICAL for security
    # This secret is used to sign and verify JWTs.
    # MUST be a strong, unique, and long random string in production.
    # If not set, a default "default secret" will be used (ONLY for local development/testing).
    JWT_SECRET="your_very_secure_jwt_key_here_e.g._generated_by_a_secret_manager"
    
    # Default Admin User Credentials for automatic bootstrapping
    # If set, the application will check for this user on startup. If not found, it will create them.
    # These should also be strong passwords in production.
    # If not set, default values "admin" and "password" will be used.
    ADMIN_USERNAME="admin"
    ADMIN_PASSWORD="adminpassword"
    ```
    **Important:** Replace the placeholder URIs and secrets with your actual values.

3.  **Add `.env` to `.gitignore`:**
    Make sure your `.env` file is not committed to version control. Add `/.env` to your `.gitignore` file.

### Installation

1.  **Navigate to the project directory:**
    ```bash
    cd task-manager
    ```

2.  **Install dependencies:**
    Go modules will handle the installation of required packages (like Gin, MongoDB driver, godotenv, bcrypt, jwt-go) automatically. Run the following command to ensure all dependencies are downloaded and the `go.sum` file is correct.
    ```bash
    go mod tidy
    ```

### Running the Application

To start the API server, run the following command from the root of the project directory (`task-manager/`):

```bash
go run ./Delivery/main.go
```

The application will perform the following steps on startup:
1.  Load environment variables from the `.env` file.
2.  Connect to MongoDB using `MONGO_URI`.
3.  Check for the default admin user specified by `ADMIN_USERNAME` and `ADMIN_PASSWORD` environment variables.
4.  If the admin user does not exist in the database, it will be automatically created with the `Admin` role.
5.  It will then set up the Gin framework server and start listening for requests.

If successful, you will see output from the MongoDB connection, admin bootstrapping status, and then the Gin framework indicating that the server is running. By default, it listens on port **8080**.

```
2024/07/18 20:39:21 Successfully connected to MongoDB for E2E tests!
...
WARNING: JWT_SECRET environment variable not set. Using default secret.
...
2025/07/18 20:39:24 Checking for default admin user 'admin'...
2025/07/18 20:39:25 Default admin user 'admin' not found. Attempting to create...
2025/07/18 20:39:25 Successfully created default admin user 'admin'
[GIN-debug] Listening and serving HTTP on :8080
```
(Logs may vary slightly based on environment and whether admin already exists)

The API is now live and ready to accept requests at `http://localhost:8080`.

### Running Tests

This project includes a comprehensive test suite for the domain, use cases, and repository layers, demonstrating the testability benefits of the Clean Architecture. The tests will connect to MongoDB using the `MONGO_TEST_URI` (or `MONGO_URI` as fallback) and utilize the `JWT_SECRET` for token operations. They also clean the test database before each test run to ensure isolation.

To run all tests:

1.  Navigate to the project root directory (`task-manager`).
2.  Ensure your `.env` file is configured with `MONGO_TEST_URI` and `JWT_SECRET`.
3.  Run the standard Go test command. It is recommended to use the `-v` flag for verbose output and `-count=1` to disable the test cache for a fresh run.

    ```bash
    go test -v -count=1 ./...
    ```

    -   **`-v`**: Enables verbose mode, which lists each test as it runs.
    -   **`-count=1`**: Disables the test cache, forcing the tests to re-run.
    -   **`./...`**: Tells Go to run tests in the current directory and all sub-directories.

A successful test run will produce extensive output from MongoDB connections, service operations, and Gin requests/responses, culminating in a `PASS` message for all test suites.

---

## Architecture Overview

This Task Management API is designed following **Clean Architecture** principles to enhance maintainability, testability, and scalability. The codebase is organized into distinct layers, each with well-defined responsibilities and clear dependency rules, ensuring a robust and flexible application structure.

### Core Principles

1.  **Separation of Concerns**: Each layer focuses on a specific aspect of the application, isolating business logic from technical details (e.g., database, web framework).
2.  **Dependency Rule**: Dependencies always flow inwards. Higher-level modules (business rules) should not depend on lower-level modules (database, web). Instead, lower-level modules depend on abstractions (interfaces) defined by higher-level modules.
3.  **Independence of Frameworks**: The core business logic (Domain and Use Cases) remains independent of external frameworks or libraries, making it easier to swap out UI, database, or external services without affecting core functionality.
4.  **Testability**: By isolating concerns and inverting dependencies, each layer, especially the core business logic, can be tested independently using mocks for its external dependencies.

### Layer Breakdown and Responsibilities

The project adheres to the following logical and physical layer structure:

```
task-manager/
├── Delivery/
│   ├── main.go
│   ├── controllers/
│   │   └── controller.go
│   └── routers/
│       └── router.go
├── Domain/
│   └── domain.go
├── Infrastructure/
│   ├── auth_middleWare.go
│   ├── jwt_service.go
│   └── password_service.go
├── Repositories/
│   ├── task_repository.go
│   └── user_repository.go
└── Usecases/
    ├── task_usecases.go
    └── user_usecases.go
```

1.  **Domain Layer (`Domain/`)**
    *   **Description**: This is the innermost layer, representing the "Enterprise Business Rules." It contains the core business entities, value objects, and general business logic that would apply regardless of the application type (web, console, etc.). It defines the interfaces (abstractions) for any external services or data access mechanisms it needs.
    *   **Contents**:
        *   `domain.go`: Defines `Task` and `User` structs (entities), `TaskStatus` and `UserRole` (value objects). Contains business rules like `NewTask` and `NewUser` for creating valid entities. It also declares interfaces like `TaskRepository`, `UserRepository`, `PasswordService`, and `JwtService`, which are contracts for data persistence and external services.
    *   **Dependencies**: Has no dependencies on any other application layers. It only imports standard Go libraries and `go.mongodb.org/mongo-driver/bson/primitive` for `ObjectID` (a practical compromise for database IDs).

2.  **Usecases Layer (`Usecases/`)**
    *   **Description**: This layer contains the "Application Business Rules." It defines and orchestrates the specific use cases or workflows of the application. It dictates *what* the application does by invoking domain entities and interacting with the interfaces defined in the Domain layer.
    *   **Contents**:
        *   `task_usecases.go`: Implements business logic related to tasks (e.g., `CreateTask`, `UpdateTask`, `DeleteTask`, `GetTaskByID`, `GetAllTasks`).
        *   `user_usecases.go`: Implements business logic related to users (e.g., `RegisterUser`, `Login`).
    *   **Dependencies**: Depends on the `Domain` layer (e.g., `domain.Task`, `domain.TaskRepository`). It does **not** depend on `Repositories`, `Infrastructure`, or `Delivery` layers directly.

3.  **Repositories Layer (`Repositories/`)**
    *   **Description**: This layer abstracts the data access logic. It contains the concrete implementations of the `*Repository` interfaces defined in the `Domain` layer. It's responsible for interacting with the database (MongoDB in this case).
    *   **Contents**:
        *   `task_repository.go`: Provides a MongoDB-specific implementation for `domain.TaskRepository`.
        *   `user_repository.go`: Provides a MongoDB-specific implementation for `domain.UserRepository`.
    *   **Dependencies**: Depends on the `Domain` layer (to implement its interfaces) and external database drivers (`go.mongodb.org/mongo-driver`).

4.  **Infrastructure Layer (`Infrastructure/`)**
    *   **Description**: This layer implements external services and utilities that the application uses, such as password hashing, JWT token management, and authentication middleware. It contains concrete implementations of the `*Service` interfaces defined in the `Domain` layer.
    *   **Contents**:
        *   `auth_middleWare.go`: Gin middleware for authentication and authorization using JWTs. It relies on `domain.JwtService`.
        *   `jwt_service.go`: Implements `domain.JwtService` for generating and parsing JWT tokens.
        *   `password_service.go`: Implements `domain.PasswordService` for hashing and comparing passwords using `bcrypt`.
    *   **Dependencies**: Depends on the `Domain` layer (to implement its interfaces) and external libraries like `github.com/dgrijalva/jwt-go` and `golang.org/x/crypto/bcrypt`.

5.  **Delivery Layer (`Delivery/`)**
    *   **Description**: This is the outermost layer, responsible for handling external input and output (e.g., HTTP requests and responses). It orchestrates the flow of data from the outside world into the Use Cases and presents the results back. It adapts data models between the network and the application's internal entities.
    *   **Contents**:
        *   `main.go`: The application entry point. It handles dependency injection by initializing concrete implementations of repositories, services, use cases, and controllers, then setting up the HTTP server and routing.
        *   `controllers/controller.go`: Contains HTTP handlers (`UserController`, `TaskController`) that receive requests, validate input (using DTOs and `gin.Context`), invoke appropriate use case methods, and format responses.
        *   `routers/router.go`: Configures the API routes and applies middleware.
    *   **Dependencies**: Depends on `Usecases` (to call business logic), `Infrastructure` (for middleware), and `Domain` (for DTOs and error handling). It is the layer that "knows" about the web framework (`Gin`).

### Guidelines for Future Development

When extending or modifying this application, consider the following guidelines to maintain the Clean Architecture structure:

*   **New Business Logic / Feature**:
    *   If a new core business entity or rule is needed, define it in the `Domain` layer (`domain.go`).
    *   Implement the application-specific business flow for the new feature within a new or existing `Usecase`. The Use Case will coordinate between `Domain` entities and the relevant `Repository` or `Service` interfaces.
*   **Changing Data Storage**:
    *   If you need to change the database (e.g., from MongoDB to PostgreSQL), you will primarily work in the `Repositories` layer. Create a new implementation (e.g., `PostgreSQLUserRepository`) that satisfies the existing `domain.UserRepository` interface. Update `main.go` to inject the new repository. No changes should be needed in `Usecases` or `Domain`.
*   **Adding New External Service Integration**:
    *   If integrating with a new external service (e.g., an email sender, a payment gateway), first define an interface for that service in the `Domain` layer.
    *   Then, create a concrete implementation of that interface within the `Infrastructure` layer.
    *   Finally, inject this new service into the relevant `Usecase` in `main.go`.
*   **Adding a New API Endpoint**:
    *   Define the endpoint in `routers/router.go`.
    *   Create a new method in an existing or new `Controller` (in `Delivery/controllers/`). This method will handle request parsing, error mapping, and calling the appropriate `Usecase` method. Avoid putting business logic directly in the controller.
*   **Error Handling**:
    *   Domain errors (`domain.ErrUserNotFound`, `domain.ErrValidationFailed`, etc.) are defined at the core and propagated upwards. Controllers then map these core errors to appropriate HTTP status codes.
*   **Testing**:
    *   **Domain**: Unit test `Domain` entities and their methods in isolation.
    *   **Usecases**: Unit test `Usecases` by mocking their `Domain` interfaces (`UserRepository`, `TaskRepository`, `PasswordService`, `JwtService`). This is crucial for fast, reliable tests independent of external systems.
    *   **Repositories/Infrastructure**: Integration test these layers with their actual external dependencies (e.g., a test MongoDB instance for `Repositories`, live JWT generation/validation for `JwtService`).
    *   **Delivery**: End-to-end tests or integration tests for controllers and routers to ensure correct HTTP interaction and response formatting.

---

## API Endpoint Reference

This section provides the detailed API contract, including new authentication and authorization endpoints, and updated security rules for existing endpoints.

**Base URL**: `http://localhost:8080`

### Task Model

This is the core data object used in the API for requests and responses.

| Field | Type | Description | Required on Create/Update |
|---|---|---|---|
| `id` | string (ObjectId hex string) | The unique identifier for the task. A 24-character hexadecimal string. Automatically generated by the server. | No |
| `title` | string | The title of the task. | **Yes** |
| `description` | string | A detailed description of the task. | No |
| `duedate` | string (RFC3339) | The due date in RFC3339 format (e.g., `"2024-12-15T17:00:00Z"`). | **Yes** |
| `status` | string | The current status of the task. Must be one of the allowed values listed below. | **Yes** |

#### Allowed Status Values
*   `"Pending"`
*   `"In progress"`
*   `"Done"`

### User & Authentication Models

#### User Model
Represents a user account in the system. Passwords are securely hashed before storage and are not returned in API responses.

| Field | Type | Description |
|---|---|---|
| `id` | string (ObjectId hex string) | Unique identifier for the user. Automatically generated. |
| `username` | string | Unique username for the user. |
| `password` | string (hashed internally) | User's password (never returned in responses). |
| `role` | string | User's assigned role. |

#### Allowed User Roles
*   `"Admin"`
*   `"User"`

#### UserRegisterLogin Model
This structure is used as the request body for both user registration and login endpoints.

| Field | Type | Description | Required |
|---|---|---|---|
| `username` | string | User's chosen username | **Yes** |
| `password` | string | User's chosen password | **Yes** |

### Common Error Responses

These are standardized error responses you can expect from various API endpoints, particularly for authentication and authorization failures.

-   **`400 Bad Request`**: The request body or URL parameter is malformed, missing required fields, or contains invalid data.
    ```json
    {
      "message": "Invalid request body"
    }
    ```
    Or:
    ```json
    {
      "message": "Invalid status provided: 'invalid_status'"
    }
    ```
    Or:
    ```json
    {
      "message": "Invalid Task ID format"
    }
    ```

-   **`401 Unauthorized`**: The request lacks valid authentication credentials (e.g., missing, malformed, expired, or invalid JWT token).
    ```json
    {
      "message": "Authorization token required"
    }
    ```
    Or:
    ```json
    {
      "message": "token expired"
    }
    ```
    Or (for login endpoint specifically):
    ```json
    {
      "message": "invalid credentials"
    }
    ```

-   **`403 Forbidden`**: The request was valid and the client is authenticated, but the authenticated user does not have the necessary permissions (e.g., insufficient role) to access the resource or perform the action.
    ```json
    {
      "message": "Access forbidden: Admin role required"
    }
    ```

-   **`404 Not Found`**: The requested resource could not be found.
    ```json
    {
      "message": "task not found"
    }
    ```

-   **`409 Conflict`**: The request could not be completed due to a conflict with the current state of the resource (e.g., attempting to register a username that already exists).
    ```json
    {
      "message": "username already taken"
    }
    ```

### Endpoints

**Base URL for all endpoints**: `http://localhost:8080`

#### Authentication

##### 1. Register a New User

Creates a new user account with a unique username and a password. Newly registered users are assigned the `"User"` role by default.

-   **Endpoint**: `POST /user/register`
-   **Authorization**: None (Public endpoint)
-   **Request Body**: `UserRegisterLogin` model.
    ```json
    {
      "username": "newuser",
      "password": "securepassword123"
    }
    ```
-   **Responses**:
    -   **`201 Created`**: User account created successfully.
        ```json
        {
          "id": "651f8b1c4e7f3c1a2d3b4e70",
          "username": "newuser",
          "role": "User"
        }
        ```
    -   **`400 Bad Request`**: Invalid request body or validation error (see [Common Error Responses](#common-error-responses)).
    -   **`409 Conflict`**: Username already taken (see [Common Error Responses](#common-error-responses)).

##### 2. User Login

Authenticates a user with provided credentials and issues a JWT upon successful login. This token must be used for subsequent requests to protected endpoints.

-   **Endpoint**: `POST /user/login`
-   **Authorization**: None (Public endpoint)
-   **Request Body**: `UserRegisterLogin` model.
    ```json
    {
      "username": "your_username",
      "password": "your_password"
    }
    ```
-   **Responses**:
    -   **`200 OK`**: Successful login. The response body contains the JWT.
        ```json
        {
          "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjUxZjhiMWM0ZTdmM2MxYTJkM2I0ZTcwIiwidXNlcm5hbWUiOiJuZXd1c2VyIiwicm9sZSI6IlVzZXIiLCJleHAiOjE3NTI5NDc2Mjh9.SomeHashValueHere"
        }
        ```
    -   **`400 Bad Request`**: Invalid request body (see [Common Error Responses](#common-error-responses)).
    -   **`401 Unauthorized`**: Invalid username or password (see [Common Error Responses](#common-error-responses)).

**How to use the JWT for Protected Endpoints:**
Once you receive a JWT from the `/user/login` endpoint, include it in the `Authorization` header of all subsequent requests to protected endpoints, using the `Bearer` scheme.

**Example Header:**
`Authorization: Bearer <your_jwt_token_here>`

#### Task Management (Protected Endpoints)

##### 1. Get All Tasks

Retrieves a list of all tasks. The order of tasks in the returned array is not guaranteed.

-   **Endpoint**: `GET /tasks`
-   **Authorization**: Required. Accessible by any **Authenticated User** (both `"Admin"` and `"User"` roles).
-   **Responses**:
    -   **`200 OK`**: A JSON array of task objects. The array will be empty (`[]`) if no tasks exist.
        ```json
        [
          {
            "id": "651f8b1c4e7f3c1a2d3b4e5f",
            "title": "Create API Documentation",
            "description": "Write a markdown file describing all endpoints.",
            "duedate": "2024-11-10T11:00:00Z",
            "status": "Pending"
          },
          {
            "id": "651f8b1c4e7f3c1a2d3b4e60",
            "title": "Fix Database Connection",
            "description": "Debug MongoDB timeout issues.",
            "duedate": "2024-10-20T09:00:00Z",
            "status": "In progress"
          }
        ]
        ```
        Empty response:
        ```json
        []
        ```
    -   **`401 Unauthorized`**: Missing or invalid token (see [Common Error Responses](#common-error-responses)).

##### 2. Get a Specific Task

Retrieves a single task by its unique ID. The ID must be a valid 24-character hexadecimal string (MongoDB ObjectID).

-   **Endpoint**: `GET /tasks/:id`
-   **Authorization**: Required. Accessible by any **Authenticated User** (both `"Admin"` and `"User"` roles).
-   **Responses**:
    -   **`200 OK`**: The response body contains the requested task object.
        ```json
        {
          "id": "651f8b1c4e7f3c1a2d3b4e5f",
          "title": "Create API Documentation",
          "description": "Write a markdown file describing all endpoints.",
          "duedate": "2024-11-10T11:00:00Z",
          "status": "Pending"
        }
        ```
    -   **`400 Bad Request`**: Invalid ID format or validation error (see [Common Error Responses](#common-error-responses)).
    -   **`401 Unauthorized`**: Missing or invalid token (see [Common Error Responses](#common-error-responses)).
    -   **`404 Not Found`**: No task exists with the specified `id` (see [Common Error Responses](#common-error-responses)).

##### 3. Create a New Task

Creates a new task.

-   **Endpoint**: `POST /tasks`
-   **Authorization**: Required. Accessible only by users with the **`"Admin"` Role**.
-   **Request Body**: A JSON object representing a `Task` (the `id` field is ignored on creation).
    ```json
    {
      "title": "Create API Documentation",
      "description": "Write a markdown file describing all endpoints.",
      "duedate": "2024-11-10T11:00:00Z",
      "status": "Pending"
    }
    ```
-   **Responses**:
    -   **`201 Created`**: The task was created successfully. The response body contains the newly created task object.
        ```json
        {
          "id": "651f8b1c4e7f3c1a2d3b4e5f",
          "title": "Create API Documentation",
          "description": "Write a markdown file describing all endpoints.",
          "duedate": "2024-11-10T11:00:00Z",
          "status": "Pending"
        }
        ```
    -   **`400 Bad Request`**: Invalid request body or validation error (see [Common Error Responses](#common-error-responses)).
    -   **`401 Unauthorized`**: Missing or invalid token (see [Common Error Responses](#common-error-responses)).
    -   **`403 Forbidden`**: Authenticated but not an Admin (see [Common Error Responses](#common-error-responses)).

##### 4. Update a Task

Updates an existing task by its ID. The ID in the URL path must be a valid 24-character hexadecimal string. Allows partial updates using `omitempty` fields in the request DTO.

-   **Endpoint**: `PUT /tasks/:id`
-   **Authorization**: Required. Accessible only by users with the **`"Admin"` Role`.
-   **Request Body**: A JSON object containing the new details for the task. Fields can be omitted if not being updated.
    ```json
    {
      "title": "Update API Documentation",
      "description": "Update the docs to include all response bodies.",
      "duedate": "2024-11-10T12:00:00Z",
      "status": "Done"
    }
    ```
    Example for partial update (only changing status):
    ```json
    {
      "status": "In progress"
    }
    ```
-   **Responses**:
    -   **`200 OK`**: The task was updated successfully. The response body contains the complete, updated task object.
        ```json
        {
          "id": "651f8b1c4e7f3c1a2d3b4e5f",
          "title": "Update API Documentation",
          "description": "Update the docs to include all response bodies.",
          "duedate": "2024-11-10T12:00:00Z",
          "status": "Done"
        }
        ```
    -   **`400 Bad Request`**: Invalid ID, invalid request body, or invalid status (e.g., changing from `Done` status) (see [Common Error Responses](#common-error-responses)).
    -   **`401 Unauthorized`**: Missing or invalid token (see [Common Error Responses](#common-error-responses)).
    -   **`403 Forbidden`**: Authenticated but not an Admin (see [Common Error Responses](#common-error-responses)).
    -   **`404 Not Found`**: No task exists with the specified `id` (see [Common Error Responses](#common-error-responses)).

##### 5. Delete a Task

Deletes a task by its ID. The ID in the URL path must be a valid 24-character hexadecimal string.

-   **Endpoint**: `DELETE /tasks/:id`
-   **Authorization**: Required. Accessible only by users with the **`"Admin"` Role`.
-   **Responses**:
    -   **`204 No Content`**: The task was deleted successfully. The response will have **no body**.
    -   **`400 Bad Request`**: Invalid ID format (see [Common Error Responses](#common-error-responses)).
    -   **`401 Unauthorized`**: Missing or invalid token (see [Common Error Responses](#common-error-responses)).
    -   **`403 Forbidden`**: Authenticated but not an Admin (see [Common Error Responses](#common-error-responses)).
    -   **`404 Not Found`**: No task exists with the specified `id` (see [Common Error Responses](#common-error-responses)).

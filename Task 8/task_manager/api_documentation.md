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
    # Used for running the main application
    MONGO_URI="mongodb+srv://<user>:<password>@<your-dev-cluster>..."
    
    # --- Test-Specific Configuration ---
    # Optional: use a separate, dedicated cluster/database for testing (HIGHLY recommended).
    # The integration and E2E tests will use this URI if set, otherwise they will fall back to MONGO_URI.
    # The test database specified here will be COMPLETELY DROPPED after tests run.
    MONGO_TEST_URI="mongodb+srv://<user>:<password>@<your-test-cluster>..."
    
    # --- JWT Configuration ---
    # This secret is used to sign and verify JWTs.
    # MUST be a strong, unique, and long random string in production.
    JWT_SECRET="your_very_secure_jwt_key_here"
    
    # Optional: A separate secret for tests. Falls back to JWT_SECRET if not set.
    JWT_TEST_SECRET="a_different_secret_just_for_testing"
    
    # --- Default Admin User Credentials for automatic bootstrapping ---
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
    Go modules will handle the installation of required packages. Run the following command to ensure all dependencies are downloaded.
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
3.  Check for and create the default admin user if it doesn't exist.
4.  Set up the Gin framework server and start listening for requests on port **8080**.

### Running Tests

This project includes a comprehensive, multi-layered test suite that validates the application at different levels, ensuring correctness, stability, and confidence in the codebase.

#### Test Strategy Overview

The testing approach follows the "Testing Pyramid" model:

1.  **Unit Tests (`Domain`, `Usecases`, `Infrastructure`):**
    *   **Purpose:** To test individual components and pure business logic in complete isolation.
    *   **Speed:** Very fast.
    *   **Dependencies:** None. These tests use mocks for all external dependencies (like databases or services) and do not require a database connection.

2.  **Integration Tests (`Repositories`):**
    *   **Purpose:** To verify that the repository layer correctly interacts with a real database.
    *   **Speed:** Slower than unit tests.
    *   **Dependencies:** Requires a live MongoDB connection, configured via `MONGO_TEST_URI` (or `MONGO_URI`) in the `.env` file. The tests run against a dedicated test database which is cleaned between tests.

3.  **End-to-End (E2E) Tests (`e2e/`):**
    *   **Purpose:** To test the entire application stack as a whole, from receiving an HTTP request to interacting with the database and returning a response. This provides the highest level of confidence.
    *   **Speed:** Slowest.
    *   **Dependencies:** Requires a live MongoDB connection (`MONGO_TEST_URI`) and a JWT secret (`JWT_TEST_SECRET`). It spins up the entire application in-memory and makes real HTTP calls to it. The test database is completely dropped after the suite runs.

#### How to Run All Tests

1.  Navigate to the project root directory (`task-manager`).
2.  Ensure your `.env` file is correctly configured, especially `MONGO_TEST_URI` and `JWT_TEST_SECRET`, as these are required for the integration and E2E tests.
3.  Run the standard Go test command. It is recommended to use the `-v` flag for verbose output and `-count=1` to disable the test cache for a fresh run.

    ```bash
    go test -v -count=1 ./...
    ```

    -   **`-v`**: Enables verbose mode, which lists each test as it runs.
    -   **`-count=1`**: Disables the test cache, forcing the tests to re-run from scratch.
    -   **`./...`**: Tells Go to run tests in the current directory and all sub-directories.

A successful test run will show PASS status for all packages (`Domain`, `Infrastructure`, `Repositories`, `Usecases`, and `e2e`), culminating in a final `ok` message for each.

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
│   └── routers/
├── Domain/
├── Infrastructure/
├── Repositories/
├── Usecases/
└── e2e/
    └── main_e2e_test.go
```

1.  **Domain Layer (`Domain/`)**: The core of the application. Contains business entities (`Task`, `User`) and the interfaces (`TaskRepository`, `JwtService`, etc.) that define the contracts for external dependencies.
2.  **Usecases Layer (`Usecases/`)**: Orchestrates application-specific workflows by coordinating Domain entities and repository/service interfaces. Contains the application's business logic.
3.  **Repositories Layer (`Repositories/`)**: Implements the data persistence interfaces defined in the Domain layer, interacting directly with MongoDB.
4.  **Infrastructure Layer (`Infrastructure/`)**: Implements other external-facing concerns defined by Domain interfaces, such as JWT handling, password hashing, and authentication middleware.
5.  **Delivery Layer (`Delivery/`)**: The outermost layer. Handles HTTP requests and responses, using the Gin framework. It wires everything together in `main.go`, but the controllers themselves are thin layers that delegate to the Usecases.

### Guidelines for Future Development

When extending or modifying this application, consider the following guidelines to maintain the Clean Architecture structure:

*   **New Business Logic**: Start in the `Domain` layer for new entities/interfaces, then implement the workflow in the `Usecases` layer.
*   **Changing Data Storage**: Create new implementations in the `Repositories` layer that satisfy the existing `Domain` interfaces. Update dependency injection in `main.go`.
*   **Adding a New API Endpoint**: Add the route in `routers/`, create a new handler method in `controllers/`, and ensure it calls the appropriate `Usecase` method.
*   **Error Handling**: Domain errors are defined at the core and propagated upwards. Controllers then map these core errors to appropriate HTTP status codes.
*   **Testing**:
    *   **Domain, Usecases, Infrastructure**: These are primarily covered by **Unit Tests**. Use mocks for all dependencies to ensure tests are fast and isolated.
    *   **Repositories**: These are covered by **Integration Tests**. These tests run against a real test database to verify data persistence logic.
    *   **Delivery (Controllers, Routers)**: The correctness of this layer is validated by **End-to-End (E2E) Tests** located in the `e2e/` directory. These tests spin up the entire application and make real HTTP requests to verify the full flow, from routing and middleware to the database and back.

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
-   **`401 Unauthorized`**: The request lacks valid authentication credentials (e.g., missing, malformed, expired, or invalid JWT token).
-   **`403 Forbidden`**: The client is authenticated, but does not have the necessary permissions (e.g., insufficient role).
-   **`404 Not Found`**: The requested resource could not be found.
-   **`409 Conflict`**: The request could not be completed due to a conflict with the current state of the resource (e.g., duplicate username).

### Endpoints

**Base URL for all endpoints**: `http://localhost:8080`

#### Authentication

##### 1. Register a New User

Creates a new user account with a unique username and a password. Newly registered users are assigned the `"User"` role by default.

-   **Endpoint**: `POST /user/register`
-   **Authorization**: None (Public endpoint)
-   **Responses**: `201 Created`, `400 Bad Request`, `409 Conflict`.

##### 2. User Login

Authenticates a user and issues a JWT.

-   **Endpoint**: `POST /user/login`
-   **Authorization**: None (Public endpoint)
-   **Responses**: `200 OK`, `400 Bad Request`, `401 Unauthorized`.

**How to use the JWT for Protected Endpoints:**
Include the token in the `Authorization` header of all subsequent requests, using the `Bearer` scheme.
**Example Header:** `Authorization: Bearer <your_jwt_token_here>`

#### Task Management (Protected Endpoints)

##### 1. Get All Tasks

Retrieves a list of all tasks.

-   **Endpoint**: `GET /tasks`
-   **Authorization**: **Authenticated User** (`Admin` or `User`).
-   **Responses**: `200 OK`, `401 Unauthorized`.

##### 2. Get a Specific Task

Retrieves a single task by its unique ID.

-   **Endpoint**: `GET /tasks/:id`
-   **Authorization**: **Authenticated User** (`Admin` or `User`).
-   **Responses**: `200 OK`, `400 Bad Request`, `401 Unauthorized`, `404 Not Found`.

##### 3. Create a New Task

Creates a new task.

-   **Endpoint**: `POST /tasks`
-   **Authorization**: **Admin Role** required.
-   **Responses**: `201 Created`, `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`.

##### 4. Update a Task

Updates an existing task by its ID. Allows for partial updates.

-   **Endpoint**: `PUT /tasks/:id`
-   **Authorization**: **Admin Role** required.
-   **Responses**: `200 OK`, `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`.

##### 5. Delete a Task

Deletes a task by its ID.

-   **Endpoint**: `DELETE /tasks/:id`
-   **Authorization**: **Admin Role** required.
-   **Responses**: `204 No Content`, `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`.

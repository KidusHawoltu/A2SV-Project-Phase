# Task Manager API

A simple RESTful API for managing tasks, built with Go and the Gin framework, backed by MongoDB. This project provides a complete, thread-safe, and well-structured backend service for basic CRUD operations.

---

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [MongoDB Configuration](#mongodb-configuration)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
  - [Running Tests](#running-tests)
- [API Endpoint Reference](#api-endpoint-reference)
  - [Task Model](#task-model)
  - [Endpoints](#endpoints)

---

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You will need the following software installed on your system:
- **Go**: Version 1.20 or later is recommended.
- **MongoDB**: Access to a MongoDB instance (e.g., a free tier on MongoDB Atlas or a local Docker instance).
- **Postman or curl**: For testing the API endpoints manually.

### MongoDB Configuration

This application connects to MongoDB using a connection URI. It reads this URI from an environment variable.

1.  **Set up a MongoDB Database:**
    *   **MongoDB Atlas (Recommended):** Sign up for a free account at [mongodb.com/atlas](https://www.mongodb.com/atlas).
        *   Create a new project and cluster.
        *   Under "Network Access", add your current public IP address (or `0.0.0.0/0` for development, **not recommended for production**) to the IP Access List. This is crucial for your application to be able to connect.
        *   Under "Database Access", create a new database user with a secure password.
        *   Go to your cluster, click "Connect", select "Drivers", choose "Go" and copy the provided connection string. It will look similar to: `mongodb+srv://<username>:<password>@<cluster-url>/<dbname>?retryWrites=true&w=majority&appName=<appname>`.

2.  **Create a `.env` file:**
    In the root directory of the project (`task_manager/`), create a file named `.env`. This file will store your MongoDB connection URI.

    ```dotenv
    # .env
    MONGO_URI="mongodb+srv://kidushawoltu:avkjCvT5PWJyt@first-cluster.fcjxr0t.mongodb.net/?retryWrites=true&w=majority&appName=First-Cluster"
    
    # Optional: use a separate cluster for testing
    # Since the app uses different database and collection, this is not necessary
    MONGO_TEST_URI="mongodb+srv://kidushawoltu:avkjCvT5PWJyt@first-cluster.fcjxr0t.mongodb.net/?retryWrites=true&w=majority&appName=First-Cluster-Test"
    ```
    **Important:** Replace the placeholder URIs with your actual MongoDB connection strings. Ensure the `MONGO_URI` is for your main application, and `MONGO_TEST_URI` points to a dedicated test database (e.g., `test_learning_phase`).

3.  **Add `.env` to `.gitignore`:**
    Make sure your `.env` file is not committed to version control. Add `/.env` to your `.gitignore` file.

### Installation

1.  **Navigate to the project directory:**
    ```bash
    cd task_manager
    ```

2.  **Install dependencies:**
    Go modules will handle the installation of required packages (like Gin, MongoDB driver, and godotenv) automatically. Run the following command to ensure all dependencies are downloaded and the `go.sum` file is correct.
    ```bash
    go mod tidy
    ```

### Running the Application

To start the API server, run the following command from the root of the project directory (`task_manager/`):

```bash
go run main.go
```

The application will attempt to connect to MongoDB using the `MONGO_URI` from your `.env` file (or directly from your environment variables if set). If successful, you will see output indicating the MongoDB connection and then the Gin framework server listening. By default, it listens on port **8080**.

```
2024/07/18 11:06:36 Attempting to connect to MongoDB...
2024/07/18 11:06:56 Pinging MongoDB deployment...
2024/07/18 11:06:56 Successfully connected to MongoDB!
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] Listening and serving HTTP on :8080
```

The API is now live and ready to accept requests at `http://localhost:8080`.

### Running Tests

This project includes a comprehensive test suite for the data service layer. The tests will also connect to MongoDB using the `MONGO_TEST_URI` (or `MONGO_URI` as fallback) and run against a separate test database (`test_learning_phase`).

To run all tests:

1.  Navigate to the project root directory.
2.  Ensure your `.env` file is configured with `MONGO_TEST_URI` (or `MONGO_URI` pointing to a test environment).
3.  Run the standard Go test command. It is recommended to use the `-v` flag for verbose output and `-count=1` to disable the test cache for a fresh run.

    ```bash
    go test -v -count=1 ./...
    ```

    - **`-v`**: Enables verbose mode, which lists each test as it runs.
    - **`-count=1`**: Disables the test cache, forcing the tests to re-run.
    - **`./...`**: Tells Go to run tests in the current directory and all sub-directories.

A successful test run will include logs about MongoDB connection and disconnection for the tests, followed by the test results:

```
=== RUN   TestMain
2024/07/18 11:06:36 Successfully connected to MongoDB for tests!
=== RUN   TestAddTask
2024/07/18 11:06:36 Cleaned test collection 'test_tasks' for new test.
--- PASS: TestAddTask (0.05s)
=== RUN   TestGetTasks
2024/07/18 11:06:36 Cleaned test collection 'test_tasks' for new test.
=== RUN   TestGetTasks/should_return_an_empty_slice_when_no_tasks_exist
2024/07/18 11:06:36 Cleaned test collection 'test_tasks' for new test.
--- PASS: TestGetTasks/should_return_an_empty_slice_when_no_tasks_exist (0.01s)
=== RUN   TestGetTasks/should_return_all_tasks_when_tasks_exist
2024/07/18 11:06:36 Cleaned test collection 'test_tasks' for new test.
--- PASS: TestGetTasks/should_return_all_tasks_when_tasks_exist (0.02s)
--- PASS: TestGetTasks (0.03s)
=== RUN   TestGetTaskById
2024/07/18 11:06:36 Cleaned test collection 'test_tasks' for new test.
=== RUN   TestGetTaskById/should_find_an_existing_task
--- PASS: TestGetTaskById/should_find_an_existing_task (0.02s)
=== RUN   TestGetTaskById/should_return_an_error_for_a_non-existent_task
--- PASS: TestGetTaskById/should_return_an_error_for_a_non-existent_task (0.01s)
=== RUN   TestGetTaskById/should_return_an_error_for_an_invalid_ObjectID_format
--- PASS: TestGetTaskById/should_return_an_error_for_an_invalid_ObjectID_format (0.00s)
--- PASS: TestGetTaskById (0.03s)
... (more tests) ...
PASS
2024/07/18 11:06:56 Disconnecting from MongoDB test database...
2024/07/18 11:06:56 MongoDB test connection closed.
ok      A2SV_ProjectPhase/Task5/TaskManager/data    20.015s
```

---

## API Endpoint Reference

This section provides the detailed API contract.

**Base URL**: `http://localhost:8080`

### Task Model

This is the core data object used in the API for requests and responses.

| Field         | Type                          | Description                                                                    | Required on Create/Update |
|---------------|-------------------------------|--------------------------------------------------------------------------------|---------------------------|
| `id`          | string (ObjectId hex string)  | The unique identifier for the task. A 24-character hexadecimal string. Automatically generated by the server. | No                        |
| `title`       | string                        | The title of the task.                                                         | **Yes**                   |
| `description` | string                        | A detailed description of the task.                                            | No                        |
| `duedate`     | string (RFC3339)              | The due date in RFC3339 format (e.g., `"2024-12-15T17:00:00Z"`).                 | **Yes**                   |
| `status`      | string                        | The current status of the task. Must be one of the allowed values listed below. | **Yes**                   |

#### Allowed Status Values
*   `"Pending"`
*   `"In progress"`
*   `"Done"`

### Endpoints

#### 1. Create a New Task

Creates a new task.

- **Endpoint**: `POST /tasks`
- **Request Body**: A JSON object representing a `Task` (the `id` field is ignored on creation).

  **Example Request:**
  ```json
  {
    "title": "Create API Documentation",
    "description": "Write a markdown file describing all endpoints.",
    "duedate": "2024-11-10T11:00:00Z",
    "status": "Pending"
  }
  ```

- **Responses**:
  - **`201 Created`**: The task was created successfully. The response body contains the newly created task object.
    ```json
    {
      "id": "651f8b1c4e7f3c1a2d3b4e5f",
      "title": "Create API Documentation",
      "description": "Write a markdown file describing all endpoints.",
      "duedate": "2024-11-10T11:00:00Z",
      "status": "Pending"
    }
    ```
  - **`400 Bad Request`**: The request body is malformed, missing required fields, or contains an invalid status.
    ```json
    {
      "message": "Invalid request body",
      "error": "json: cannot unmarshal string into Go struct field Task.status of type models.TaskStatus"
    }
    ```
    Or:
    ```json
    {
      "message": "Invalid status provided: 'invalid_status'",
      "error": ""
    }
    ```

#### 2. Get All Tasks

Retrieves a list of all tasks. The order of tasks in the returned array is not guaranteed.

- **Endpoint**: `GET /tasks`

- **Responses**:
  - **`200 OK`**: A JSON array of task objects. The array will be empty (`[]`) if no tasks exist.
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

#### 3. Get a Specific Task

Retrieves a single task by its unique ID. The ID must be a valid 24-character hexadecimal string (MongoDB ObjectID).

- **Endpoint**: `GET /tasks/:id`

- **Responses**:
  - **`200 OK`**: The response body contains the requested task object.
    ```json
    {
      "id": "651f8b1c4e7f3c1a2d3b4e5f",
      "title": "Create API Documentation",
      "description": "Write a markdown file describing all endpoints.",
      "duedate": "2024-11-10T11:00:00Z",
      "status": "Pending"
    }
    ```
  - **`400 Bad Request`**: The provided `id` in the URL is not a valid 24-character hexadecimal string.
    ```json
    {
      "message": "Invalid Task ID format",
      "error": "encoding/hex: invalid byte: U+0061 'a'"
    }
    ```
  - **`404 Not Found`**: No task exists with the specified `id`.
    ```json
    {
      "message": "Task with ID '651f8b1c4e7f3c1a2d3b4e99' not found",
      "error": "task not found: task with id '651f8b1c4e7f3c1a2d3b4e99'"
    }
    ```

#### 4. Update a Task

Updates an existing task by its ID. The ID in the URL path must be a valid 24-character hexadecimal string.

- **Endpoint**: `PUT /tasks/:id`
- **Request Body**: A JSON object containing the new details for the task. All fields in the request body are used for the update.

  **Example Request:**
  ```json
  {
    "title": "Update API Documentation",
    "description": "Update the docs to include all response bodies.",
    "duedate": "2024-11-10T12:00:00Z",
    "status": "Done"
  }
  ```

- **Responses**:
  - **`200 OK`**: The task was updated successfully. The response body contains the complete, updated task object.
    ```json
    {
      "id": "651f8b1c4e7f3c1a2d3b4e5f",
      "title": "Update API Documentation",
      "description": "Update the docs to include all response bodies.",
      "duedate": "2024-11-10T12:00:00Z",
      "status": "Done"
    }
    ```
  - **`400 Bad Request`**: The provided `id` in the URL is invalid, the request body is malformed, or the status is invalid.
    ```json
    {
      "message": "Invalid Task ID format",
      "error": "encoding/hex: invalid byte: U+0061 'a'"
    }
    ```
    Or:
    ```json
    {
      "message": "Invalid request body",
      "error": "json: cannot unmarshal string into Go struct field Task.status of type models.TaskStatus"
    }
    ```
    Or:
    ```json
    {
      "message": "Invalid status provided: 'invalid_status'",
      "error": ""
    }
    ```
  - **`404 Not Found`**: No task exists with the specified `id`.
    ```json
    {
      "message": "Task with ID '651f8b1c4e7f3c1a2d3b4e99' not found",
      "error": "task not found: task with id '651f8b1c4e7f3c1a2d3b4e99'"
    }
    ```

#### 5. Delete a Task

Deletes a task by its ID. The ID in the URL path must be a valid 24-character hexadecimal string.

- **Endpoint**: `DELETE /tasks/:id`

- **Responses**:
  - **`204 No Content`**: The task was deleted successfully. The response will have **no body**.
  - **`400 Bad Request`**: The provided `id` in the URL is not a valid 24-character hexadecimal string.
    ```json
    {
      "message": "Invalid Task ID format",
      "error": "encoding/hex: invalid byte: U+0061 'a'"
    }
    ```
  - **`404 Not Found`**: No task exists with the specified `id`.
    ```json
    {
      "message": "Task with ID '651f8b1c4e7f3c1a2d3b4e99' not found",
      "error": "task not found: task with id '651f8b1c4e7f3c1a2d3b4e99'"
    }
    ```

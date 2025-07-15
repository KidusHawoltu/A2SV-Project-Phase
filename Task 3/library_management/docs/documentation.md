# Console-Based Library Management System

## 1. Overview

This project is a simple, yet robust, console-based application for managing a library's books and members. It is written in the Go programming language and serves as a practical example of key software engineering principles, including:

*   **Separation of Concerns:** The project is divided into distinct layers for models (data structures), services (business logic), and controllers (user interaction).
*   **Dependency Injection:** The service layer is "injected" into the controller, promoting decoupled and highly testable code.
*   **Interface-Based Design:** A `LibraryManager` interface defines the contract for library operations, allowing for flexible and mockable implementations.
*   **Test-Driven Development:** A comprehensive test suite for the service layer ensures the correctness and stability of the core logic.

## 2. Features

The system provides the following functionalities through an interactive console menu:

*   **Member Management:**
    *   Add new members to the library.
    *   Remove existing members.
    *   List all members and the books they have borrowed.
*   **Book Management:**
    *   Add new books to the library.
    *   Remove available books from the library.
*   **Core Library Operations:**
    *   Allow a member to borrow an available book.
    *   Allow a member to return a borrowed book.
    *   List all books currently available for borrowing.
    *   List all books borrowed by a specific member.

## 3. Project Structure

The project follows a clean, layered architecture:
## 3. Project Structure

The project follows a clean, layered architecture:

```text
library_management/
├── main.go                     # Entry point of the application
├── go.mod                      # Go module definition
├── controllers/
│   └── library_controller.go   # Handles user input and console UI
├── models/
│   ├── book.go                 # Defines the Book struct and status
│   └── member.go               # Defines the Member struct
├── services/
│   ├── library_service.go      # Business logic implementation
│   └── library_service_test.go # Unit tests for the service layer
└── docs/
    └── documentation.md        # This documentation file
```

*   **`main.go`**: Initializes the `Library` service and injects it into the controller to start the application.
*   **`controllers/`**: The presentation layer. It is responsible for printing menus, reading user input, and calling the appropriate service methods. It does not contain any business logic.
*   **`models/`**: The data layer. It defines the core data structures of the application, such as `Book` and `Member`.
*   **`services/`**: The business logic layer. It contains the `Library` struct which manages all data and enforces the rules of the system (e.g., a borrowed book cannot be removed).

## 4. Core Components Deep Dive

### Models

*   **`Book`**: Represents a book with an `Id`, `Title`, `Author`, and a `Status`.
*   **`BookStatus`**: A custom type (`Available` or `Borrowed`) to ensure type safety for book statuses.
*   **`Member`**: Represents a library member with an `Id`, `Name`, and a slice of pointers to `BorrowedBooks` (`[]*Book`) to efficiently track borrowed items.

### Services

*   **`LibraryManager` Interface**: This interface defines the contract for all library operations. By having the controller depend on this interface instead of the concrete `Library` struct, we achieve loose coupling.
*   **`Library` Struct**: The concrete implementation of `LibraryManager`. It holds maps of pointers to `Books` and `Members` (`map[int]*models.Book`, `map[int]*models.Member`) for efficient state management.
*   **Auto-Incrementing IDs**: The `Library` service automatically assigns a unique, sequential ID to each new book and member, simplifying the user's interaction with the system.
*   **State Management (Pointer-Based)**: The system has been refactored to use pointers for state management, which is a common and efficient idiom in Go. Instead of storing copies of `Book` and `Member` structs, the library's maps hold pointers to these structs. This means that when an object's field (like a book's status) is modified, the change is reflected immediately across the entire system without needing to write the object back into the map. This approach is memory-efficient and eliminates a class of bugs related to stale data. To ensure pointer stability, new books and members are explicitly copied to a new variable before their address is taken and stored, following Go best practices.

## 5. How to Run the Application

**Prerequisites:**
*   Go (version 1.18 or later) installed and configured.

**Steps:**
1.  Clone or download the project to your local machine.
2.  Open a terminal and navigate to the root directory of the project (`library_management/`).
3.  Run the application with the following command:
    ```bash
    go run main.go
    ```
4.  The application will start, and you can interact with it by entering numbers corresponding to the menu options.

## 6. How to Run Tests

The project includes a comprehensive test suite for the service layer to ensure its logic is correct.

1.  Open a terminal and navigate to the root directory of the project.
2.  Run all tests using the following command:
    ```bash
    go test ./services/
    ```
3.  For a more detailed, verbose output that lists each test case, use the `-v` flag:
    ```bash
    go test -v ./services/
    ```

## 7. Menu Options Guide

| Input | Action                                       |
| :---- | :------------------------------------------- |
| `0`   | Display the menu again.                      |
| `1`   | Add a new member to the library.             |
| `2`   | Remove a member from the library.            |
| `3`   | Add a new book to the library.               |
| `4`   | Remove an existing book from the library.    |
| `5`   | Borrow a book for a member.                  |
| `6`   | Return a borrowed book.                      |
| `7`   | List all members in the library.             |
| `8`   | List all available books.                    |
| `9`   | List the books borrowed by a specific member.|
| Other | Exit the program.                            |
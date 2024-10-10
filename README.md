# User Management API in Go

## Project Overview

This project implements a simple RESTful API for user management, including creating, retrieving, updating, and deleting user records in a MySQL database.

### Features:
- CRUD operations for user management.
- Configurable database connection.
- Integration tests for API routes.
- Structured project layout for maintainability.

## Getting Started

### Prerequisites:
- Go (version 1.16 or newer)
- MySQL
- Git

### Installation:

1. Clone the repository:
    ```
    git clone https://github.com/your-repo/go-crud-api
    cd go-crud-api
    ```

2. Install Go dependencies:
    ```
    go mod tidy
    ```

3. Set up the environment variables by copying `.env.example` to `.env` and modifying the values:
    ```
    cp .env.example .env
    ```

4. Start the MySQL server and create a database:
    ```
    CREATE DATABASE go_crud_api;
    ```

5. Run the application:
    ```
    go run cmd/api/main.go
    ```

## API Endpoints:

- `POST /user` - Create a new user.
- `GET /users` - Retrieve all users.
- `GET /user/{id}` - Retrieve a user by ID.
- `PUT /user/{id}` - Update a user by ID.
- `DELETE /user/{id}` - Delete a user by ID.

## Running Tests:

To run the integration tests, use the following command:

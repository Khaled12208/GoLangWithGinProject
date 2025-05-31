# GolangWithGin API Postman Collection

This directory contains Postman collection and environment files for testing the GolangWithGin API.

## Files

- `golangwithgin.postman_collection.json`: The main Postman collection containing all API endpoints
- `golangwithgin.postman_environment.json`: Environment variables for local testing

## Setup Instructions

1. Import the collection file (`golangwithgin.postman_collection.json`) into Postman
2. Import the environment file (`golangwithgin.postman_environment.json`) into Postman
3. Select the "GolangWithGin Local" environment in Postman

## Available Endpoints

### Auth

- `POST /api/v1/register`: Register a new user
- `POST /api/v1/login`: Login with existing credentials

### Users

- `GET /api/v1/user`: Get current user profile
- `PUT /api/v1/user`: Update current user profile
- `DELETE /api/v1/user/{id}`: Delete a user

### Tasks

- `POST /api/v1/tasks`: Create a new task
- `GET /api/v1/tasks/{id}`: Get a specific task
- `GET /api/v1/tasks`: Get all tasks

## Environment Variables

- `base_url`: The base URL for the API (default: http://localhost:8888)
- `auth_token`: JWT token received after login
- `user_id`: Current user ID
- `task_id`: Task ID for task-specific operations

## Testing Flow

1. Register a new user using the Register endpoint
2. Login with the registered user credentials
3. Copy the received token into the `auth_token` environment variable
4. Use other endpoints with the token in the Authorization header

## Notes

- All protected endpoints require a valid JWT token in the Authorization header
- The token format should be: `Bearer <token>`
- Make sure the API server is running before testing
- The default port is 8888, update the `base_url` if using a different port

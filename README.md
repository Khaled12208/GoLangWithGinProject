# Go Task Management API with Gin Framework

A robust and scalable task management API built with Go and the Gin framework, featuring JWT authentication, Swagger documentation, asynchronous task processing, and database sharding.

## Table of Contents

- [Architecture](#architecture)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Project Structure](#project-structure)

## Architecture

### C4 Model Diagrams

#### System Context Diagram

```mermaid
graph TB
    CA[Client Apps]
    TMS[Task Management System]
    DB[(MySQL Database)]

    CA -->|Uses| TMS
    TMS -->|Stores data| DB

    classDef system fill:#1168bd,stroke:#0b4884,color:white
    classDef external fill:#666,stroke:#333,color:white
    classDef database fill:#2b8a3e,stroke:#1b5a25,color:white

    class TMS system
    class CA external
    class DB database
```

#### Container Diagram

```mermaid
graph TB
    CA[Client Apps]
    API["API Server
(Go/Gin)"]
    TP[Task Processor]
    DB[(MySQL Database)]

    CA -->|Uses| API
    API -->|Processes Tasks| TP
    API -->|Stores Data| DB
    TP -->|Updates| DB

    classDef system fill:#1168bd,stroke:#0b4884,color:white
    classDef external fill:#666,stroke:#333,color:white
    classDef database fill:#2b8a3e,stroke:#1b5a25,color:white

    class API,TP system
    class CA external
    class DB database
```

#### Component Diagram

```mermaid
graph LR
    H[Handlers]
    S[Services]
    R[Repository]
    D[Domain]

    H -->|Uses| S
    S -->|Uses| R
    S -->|Uses| D
    R -->|Uses| D

    classDef component fill:#1168bd,stroke:#0b4884,color:white
    class H,S,R,D component
```

### Sequence Diagram

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant S as Service
    participant SM as Shard Manager
    participant R as Repository

    C->>H: Create Task
    H->>S: Process Task
    S->>SM: Determine Shard
    SM->>R: Route to Shard
    R->>R: Save Task
    H-->>C: 202 Accepted
    Note over S,R: Background Processing

    C->>H: Get Task Status
    H->>S: Get Task
    S->>SM: Get Shard
    SM->>R: Fetch from Shard
    R-->>S: Task Data
    S-->>H: Task Status
    H-->>C: 200 OK
```

### Database Architecture

```mermaid
erDiagram
    Users ||--o{ Tasks : creates
    Users {
        bigint id PK
        varchar username UK
        varchar password
        varchar email UK
        timestamp created_at
        timestamp updated_at
    }
    Tasks {
        bigint id PK
        bigint user_id FK
        varchar title
        text description
        varchar status
        timestamp created_at
        timestamp updated_at
    }
```

The system uses a single MySQL database with the following features:

- **One-to-Many Relationship**: Each user can have multiple tasks
- **Indexing**: Optimized queries with indexes on frequently accessed columns (status, created_at, user_id)
- **Timestamps**: Automatic tracking of creation and update times
- **Status Management**: Pre-defined task states (pending, processing, completed, failed)
- **Data Integrity**: Foreign key constraints and unique constraints where appropriate

## Technologies Used

- **Go 1.24.3**: Main programming language
- **Gin Framework**: HTTP web framework
- **MySQL 8.0**: Database with sharding support
- **JWT**: Authentication
- **Swagger/OpenAPI**: API documentation
- **Docker & Docker Compose**: Containerization
- **GORM**: ORM for database operations
- **Consistent Hashing**: For database sharding
- **Viper**: Configuration management
- **Logrus**: Logging
- **Testify**: Testing framework
- **Docker**: Containerization
- **Shell Scripting**: Automation

## Getting Started

### Prerequisites

- Go 1.24.3 or later
- Docker and Docker Compose
- MySQL 8.0 or later (if running locally)

### Installation

1. Clone the repository:

```bash
git clone git@github.com:Khaled12208/GoLangWithGinProject.git
cd GoLangWithGinProject
```

2. Install dependencies:

```bash
go mod download
```

3. Install Swagger:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Running the Application

### Local Development (Default & Recommended)

Start the application with Swagger UI:

```bash
./run.sh local
```

To stop all services:

```bash
./run.sh stop
```

To clean up resources:

```bash
./run.sh clean
```

Swagger UI will be available at: http://localhost:8888/swagger/index.html

## API Documentation

- Swagger UI: http://localhost:8888/swagger/index.html
- Postman Collection: [Link to Postman Collection](./postman/GoLangWithGin.postman_collection.json)

### Main Endpoints

- POST `/api/v1/register`: Register a new user
- POST `/api/v1/login`: Login and get JWT token
- GET `/api/v1/user`: Get user profile
- PUT `/api/v1/user`: Update user profile
- POST `/api/v1/tasks`: Create a new task
- GET `/api/v1/tasks`: List all tasks
- GET `/api/v1/tasks/:id`: Get task by ID

## Testing

Run different types of tests using the script:

```bash
# Run unit tests
./run.sh test

# Run integration tests
./run.sh integration

# Run all tests with coverage
./run.sh test-all

# Generate coverage report
./run.sh coverage
```

## Project Structure

```
.
├── cmd/                  # Application entrypoints
├── config/              # Configuration files
├── docs/               # Swagger documentation
├── internal/           # Private application code
│   ├── app/           # Application setup
│   ├── domain/        # Domain models
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic
│   └── utils/         # Utilities
├── pkg/               # Public libraries
├── scripts/           # Helper scripts
└── tests/             # Test suites
```

```
 _  __  _   _    _    _     _____ ____
| |/ / | | | |  / \  | |   | ____|  _ \
| ' /  | |_| | / _ \ | |   |  _| | | | |
| . \  |  _  |/ ___ \| |___| |___| |_| |
|_|\_\ |_| |_/_/   \_\_____|_____|____/
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

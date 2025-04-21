# Package Tracking Microservice

A Go microservice for tracking packages with MongoDB backend, following clean architecture principles.

## Features

- RESTful API for package tracking
- MongoDB integration
- Clean Architecture implementation
- Swagger documentation
- Pagination support
- Search functionality
- Unit tests

## Prerequisites

- Go 1.21 or later
- MongoDB 4.4 or later
- Docker (optional)

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/snavarro/microtracker.git
   cd microtracker
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Start MongoDB:
   ```bash
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   ```

4. Run the service:
   ```bash
   go run main.go
   ```

The service will start on `http://localhost:8080`.

## API Documentation

Swagger documentation is available at `http://localhost:8080/swagger/index.html`

## API Endpoints

- `GET /api/v1/packages` - List all packages (with pagination)
- `GET /api/v1/packages/search` - Search packages
- `GET /api/v1/packages/:id` - Get a package by ID
- `POST /api/v1/packages` - Create a new package
- `PUT /api/v1/packages/:id` - Update a package
- `DELETE /api/v1/packages/:id` - Delete a package

## Running Tests

```bash
go test ./...
```

## Package Structure

```
.
├── config/             # Configuration
├── internal/
│   ├── domain/        # Domain models and interfaces
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic
│   └── handler/       # HTTP handlers
├── main.go            # Application entry point
└── README.md
```

## Environment Variables

The service can be configured using the following environment variables:

- `MONGO_URI` - MongoDB connection string (default: "mongodb://localhost:27017")
- `DATABASE_NAME` - MongoDB database name (default: "tracker")
- `SERVER_ADDRESS` - Server address (default: ":8080") 
# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install MySQL client
RUN apk add --no-cache mysql-client

# Copy the binary and necessary files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/config ./config
COPY --from=builder /app/scripts ./scripts

# Make scripts executable
RUN chmod +x /app/scripts/wait-for-it.sh

# Expose port
EXPOSE 8888

# Set the entrypoint to use wait-for-it
ENTRYPOINT ["/app/scripts/wait-for-it.sh", "mysql", "./main"] 
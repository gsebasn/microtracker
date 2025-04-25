# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o microtracker

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/microtracker .

# Copy environment file
COPY .env.development .env.development

# Expose port 9090
EXPOSE 9090

# Set environment variables (as fallbacks)
ENV APP_ENV=development \
    MONGO_URI=mongodb://host.docker.internal:27017 \
    DATABASE_NAME=tracker_dev \
    SERVER_ADDRESS=:9090 \
    RATE_LIMIT_REQUESTS_PER_MINUTE=100 \
    RATE_LIMIT_BURST_SIZE=50 \
    RATE_LIMIT_TTL_MINUTES=5 \
    RATE_LIMIT_LIST_REQUESTS_PER_MINUTE=200 \
    RATE_LIMIT_LIST_BURST_SIZE=100 \
    RATE_LIMIT_LIST_TTL_MINUTES=5 \
    RATE_LIMIT_SEARCH_REQUESTS_PER_MINUTE=150 \
    RATE_LIMIT_SEARCH_BURST_SIZE=75 \
    RATE_LIMIT_SEARCH_TTL_MINUTES=5 \
    RATE_LIMIT_CREATE_REQUESTS_PER_MINUTE=50 \
    RATE_LIMIT_CREATE_BURST_SIZE=25 \
    RATE_LIMIT_CREATE_TTL_MINUTES=5

# Run the application
CMD ["./microtracker"] 
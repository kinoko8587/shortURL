# Step 1: Build stage
FROM golang:1.23-alpine AS build

# Install dependencies for building the Go app
RUN apk add --no-cache git build-base

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download Go dependencies
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go app and output binary as main
RUN go build -o main .

# Step 2: Final stage for runtime (hot reload using air)
FROM alpine:latest

# Install dependencies required for the final runtime
RUN apk add --no-cache bash libc6-compat curl

# Install air for hot reloading
RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b /usr/local/bin

# Set the working directory
WORKDIR /app

# Copy the compiled Go binary from the build stage
COPY --from=build /app/main /app/main

# Copy the source code for air to watch
COPY . .

# Copy the air.toml configuration file
COPY air.toml .

# Expose the application port
EXPOSE 8080

# Use air to hot reload the pre-built binary
CMD ["air"]

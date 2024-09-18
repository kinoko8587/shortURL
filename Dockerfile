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
RUN go install github.com/air-verse/air@latest

# Expose the application port
EXPOSE 8080

# Use air to hot reload the pre-built binary
CMD ["air", "-c", ".air.toml"]

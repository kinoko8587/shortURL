# Development stage with air for hot reload
FROM golang:1.22-alpine

# Install dependencies
RUN apk add --no-cache git build-base

# Set the working directory
WORKDIR /app

# Install air for hot reload - using v1.52.3 which is compatible with Go 1.22
RUN go install github.com/air-verse/air@v1.52.3

# Copy go.mod and go.sum to download Go dependencies
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the entire source code
COPY . .

# Create .air.toml if it doesn't exist
RUN if [ ! -f .air.toml ]; then \
    air init; \
    fi

# Expose the application port
EXPOSE 8080

# Use air to hot reload
CMD ["air", "-c", ".air.toml"]

# Step 1: Build stage
FROM golang:1.23-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Run stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the build stage
COPY --from=build /app/main .

# Expose the port on which the app runs
EXPOSE 8080

# Command to run the application
CMD ["./main"]

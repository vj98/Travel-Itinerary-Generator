# Start with the official Go image as the base
FROM golang:1.19-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o myapp .

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/myapp .

# Copy the environment configuration file
COPY PROD.env .

# Expose the port on which the application will run
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]

# Use the official Golang image as the base image
FROM golang:1.21.0 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o logistics-app ./cmd/main.go

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/logistics-app .

# Create a directory for the web frontend
RUN mkdir -p /app/web/build/static

# Expose the port the app runs on
EXPOSE 6060

# Command to run the application
CMD ["./logistics-app"]
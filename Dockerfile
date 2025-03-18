# Use the official Golang image with Alpine for a small footprint
FROM golang:1.24.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (to cache dependency downloads)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project source code into the container
COPY . .

# Build the application (this will generate a `main` binary)
RUN go build -o main ./cmd

# Expose the application port (if needed)
EXPOSE 8080

# Start the application
CMD ["./main"]

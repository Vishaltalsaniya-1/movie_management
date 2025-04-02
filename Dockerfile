# Start with a Go base image
FROM golang:1.23-alpine


# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the Go source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the application port (adjust if needed)
EXPOSE 8081

# Run the compiled binary
CMD ["./main"]

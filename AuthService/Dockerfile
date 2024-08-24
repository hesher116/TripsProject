# Start with a Golang base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /cmd

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main cmd/main.go

# Expose the port the app runs on
EXPOSE 50051

# Command to run the executable
CMD ["./main"]

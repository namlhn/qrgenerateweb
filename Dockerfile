# Stage 1: Build the Go binary
FROM golang:1.23

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source code
COPY . /app



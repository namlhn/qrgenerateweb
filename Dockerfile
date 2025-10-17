# Stage 1: Build the Go binary
FROM golang:1.24.3

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download


# Make sure GOPATH/bin is in PATH
ENV PATH="$PATH:$(go env GOPATH)/bin"
# Copy the source code
COPY . /app



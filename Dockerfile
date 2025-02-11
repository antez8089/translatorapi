 Use the official Golang image
FROM golang:1.21-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /translatorapi

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o app .

# Expose port for GraphQL API
EXPOSE 8080

# Run the application
CMD ["./app"]

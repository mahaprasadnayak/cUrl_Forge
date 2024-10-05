# Use Golang image for the build stage
FROM golang:1.16-alpine AS builder

WORKDIR /app
# Create and change the directory to /app

COPY go.* ./
# Copy Go module files

RUN go mod download
# Download Go modules

COPY . ./
# Copy local files to the container

RUN go build -o server
# Build the Go application

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Install necessary libraries (if needed)
RUN apk --no-cache add libc6-compat

# Copy the binary from the builder stage
COPY --from=builder /app/server /app/server

# Set the working directory and command to run the application
WORKDIR /app

CMD ["/app/server"] 
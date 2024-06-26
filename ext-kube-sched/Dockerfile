# Use the official Golang image to build the application
FROM golang:1.20 as builder

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o scheduler .

# Use an Alpine base image with glibc
FROM alpine:latest

# Install glibc
RUN apk add --no-cache libc6-compat

# Copy the built application from the builder
COPY --from=builder /app/scheduler /scheduler

# Expose the port on which the application runs
EXPOSE 8080

# Run the application
ENTRYPOINT ["/scheduler"]

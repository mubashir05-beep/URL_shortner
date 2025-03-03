# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /url_Shortener

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o /url-shortener ./main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /url-shortener .

# Ensure .env is copied (if required)
COPY .env .env

# Install dependencies like curl (optional, for debugging)
RUN apk add --no-cache curl

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./url-shortener"]

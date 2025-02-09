# Use Go image for building the API Gateway
FROM golang:1.22 AS builder
WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . ./

# Build the API Gateway
RUN CGO_ENABLED=0 go build -o api-gateway-service ./main.go

# Minimal execution environment
FROM debian:bullseye-slim

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary and .env file from builder stage
COPY --from=builder /app/api-gateway-service /app/api-gateway-service
COPY --from=builder /app/.env /app/.env

# Expose API Gateway port
EXPOSE 8080

# Run the API Gateway
CMD ["/app/api-gateway-service"]

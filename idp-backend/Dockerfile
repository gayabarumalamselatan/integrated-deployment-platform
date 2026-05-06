# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o orchestrator ./cmd/orchestrator/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/orchestrator .
# Copy kubeconfig if it exists in the root
COPY --from=builder /app/kubeconfig ./kubeconfig

# Expose port (adjust if different)
EXPOSE 8080

CMD ["./orchestrator"]

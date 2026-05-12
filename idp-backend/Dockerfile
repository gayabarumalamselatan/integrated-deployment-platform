# Build stage
FROM golang:alpine AS builder
ENV GOTOOLCHAIN=auto

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

# Install certificates for HTTPS requests and git for cloning repos
RUN apk --no-cache add ca-certificates git

# Copy binary from builder
COPY --from=builder /app/orchestrator .
# Copy kubeconfig if it exists in the root
COPY --from=builder /app/kubeconfig ./kubeconfig

# Expose port (adjust if different)
EXPOSE 8080

CMD ["./orchestrator"]

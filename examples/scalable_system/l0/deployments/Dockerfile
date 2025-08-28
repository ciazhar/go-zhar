# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary (path fixed)
RUN go build -o scalable-system ./cmd/scalable_system

# Stage 2: Runtime
FROM alpine:3.19

WORKDIR /app

# Copy binary only
COPY --from=builder /app/scalable-system .

# Copy configs (if needed inside container)
COPY configs ./configs

# Expose port
EXPOSE 8080

# Run app
CMD ["./scalable-system"]

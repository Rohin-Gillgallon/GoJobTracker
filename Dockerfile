# Stage 1: Builder
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy mod files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/server

# Stage 2: Final minimal image
FROM alpine:latest

WORKDIR /app

# Add non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/bin/server .

# Use non-root user
USER appuser

EXPOSE 8080

CMD ["./server"]
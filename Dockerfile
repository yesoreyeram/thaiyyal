# Build stage
FROM golang:1.24.7-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./backend/

# Build the server
RUN cd backend/cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /build/thaiyyal-server .

# Runtime stage
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 thaiyyal && \
    adduser -D -u 1000 -G thaiyyal thaiyyal

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/thaiyyal-server /app/thaiyyal-server

# Change ownership
RUN chown -R thaiyyal:thaiyyal /app

# Switch to non-root user
USER thaiyyal

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health/live || exit 1

# Run the server
ENTRYPOINT ["/app/thaiyyal-server"]
CMD ["-addr", ":8080"]

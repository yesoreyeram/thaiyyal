# Multi-stage build for Thaiyyal workflow engine
# Stage 1: Build frontend
FROM node:alpine AS frontend-builder

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY next.config.ts ./
COPY tsconfig.json ./
COPY postcss.config.mjs ./
COPY eslint.config.mjs ./

# Install dependencies
RUN npm ci

# Copy source files
COPY src ./src
COPY public ./public

# Build Next.js application (export mode only, don't copy to backend yet)
RUN npm run build:frontend

# Stage 2: Build backend
FROM golang:alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source
COPY backend ./backend

# Create static directory and copy frontend build output
RUN mkdir -p backend/pkg/server/static

# Copy frontend build output to backend static directory
# Next.js export mode creates an 'out' directory with all static files
COPY --from=frontend-builder /app/out ./backend/pkg/server/static/

# Build Go binary
RUN cd backend/cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/thaiyyal-server .

# Stage 3: Final runtime image
FROM alpine:3.19

# Install ca-certificates for HTTPS and wget for health checks
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -g 1000 thaiyyal && \
    adduser -D -u 1000 -G thaiyyal thaiyyal

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /app/thaiyyal-server /app/thaiyyal-server

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

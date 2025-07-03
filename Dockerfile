# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for building
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with security flags and optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -buildvcs=false \
    -trimpath \
    -ldflags "-s -w -extldflags '-static'" \
    -o reservio ./cmd/main.go

# Final stage - use distroless for minimal attack surface
FROM gcr.io/distroless/static-debian12:nonroot

# Copy SSL certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/reservio /app/reservio

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Use non-root user (already set in distroless)
USER 65532:65532

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/reservio", "-health-check"] || exit 1

# Run the binary
ENTRYPOINT ["/app/reservio"] 
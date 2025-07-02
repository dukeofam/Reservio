# syntax=docker/dockerfile:1
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o reservio ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/reservio .
COPY --from=builder /app/.env.example ./
COPY --from=builder /app/README.md ./
COPY --from=builder /app/setup_db.sh ./
COPY --from=builder /app/setup_test_db.sh ./
COPY --from=builder /app/run_tests.sh ./
EXPOSE 8080
# Copy .env.example to .env if .env is missing (for dev)
CMD ["/bin/sh", "-c", "[ -f .env ] || cp .env.example .env; ./reservio"] 
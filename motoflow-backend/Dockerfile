# ==============================================================================
# Multi-stage Dockerfile for MotoFlow Backend (Golang)
# ==============================================================================

# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile statically linked binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o motoflow-server .

# Production Image
FROM alpine:3.19

WORKDIR /app

# Install security certificates & set Jakarta Timezone
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" > /etc/timezone

# Create non-root user
RUN adduser -D -u 10001 appuser
USER appuser

# Copy compiled binary from builder
COPY --from=builder /app/motoflow-server /app/motoflow-server

EXPOSE 8080

ENTRYPOINT ["/app/motoflow-server"]

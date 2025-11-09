# Multi-stage Dockerfile for kan-internal-services
# Builder stage
FROM golang:1.21 AS builder
WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the sources and build
COPY . .
# Build without forcing CGO_ENABLED=0 (let Go use system defaults)
RUN go build -ldflags="-s -w" -o /kan-internal-services

# Final stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /kan-internal-services /app/kan-internal-services

# Hint for exposed port
EXPOSE 8080
# Default environment variable for port
ENV PORT=8080

ENTRYPOINT ["/app/kan-internal-services"]

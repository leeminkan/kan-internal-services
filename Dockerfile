# Multi-stage Dockerfile for kan-internal-services


# Builder stage
ARG BUILDARCH=amd64
FROM --platform=linux/${BUILDARCH} golang:1.21-bullseye AS builder
WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the sources and build
COPY . .
# Force static build for Alpine compatibility and allow arch switching
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${BUILDARCH} go build -ldflags="-s -w" -o /kan-internal-services

# Final stage
FROM --platform=$BUILDPLATFORM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /kan-internal-services /app/kan-internal-services

# Hint for exposed port
EXPOSE 8080
# Default environment variable for port
ENV PORT=8080

ENTRYPOINT ["/app/kan-internal-services"]

# syntax=docker/dockerfile:1

# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and embedded assets
COPY cmd/ cmd/
COPY internal/ internal/
COPY VERSION ./

# Build optimized static binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X 'github.com/codewithsorabh/adbot/internal/version.Version=$(cat VERSION)'" \
    -o /bin/adbot \
    ./cmd/adbot

# Stage 2: Minimal production runtime
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/adbot /usr/local/bin/adbot

ENTRYPOINT ["/usr/local/bin/adbot"]

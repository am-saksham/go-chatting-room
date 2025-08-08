# -------- Stage 1: Build --------
FROM golang:1.24.6-alpine AS builder

# Install build dependencies
RUN apk update && apk upgrade --no-cache && apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o server

# -------- Stage 2: Run --------
FROM alpine:3.19

RUN apk update && apk upgrade --no-cache

WORKDIR /app

COPY --from=builder /app/server .

# Create non-root user
RUN adduser -D appuser
USER appuser

EXPOSE 8080
CMD ["./server"]
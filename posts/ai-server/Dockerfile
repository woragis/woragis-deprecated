# Use Go 1.21 alpine as builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o posts-ai-service cmd/server/main.go

# Final stage - minimal runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/posts-ai-service .

# Expose port
EXPOSE 3014

# Run the binary
CMD ["./posts-ai-service"]

# Copy .env if exists
COPY .env* ./

EXPOSE 3014

CMD ["./posts-ai-service"]

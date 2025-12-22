FROM golang:1.23-alpine AS builder

WORKDIR /build

# Enable automatic toolchain download if needed
ENV GOTOOLCHAIN=auto

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/bin/translation-worker ./cmd/translation-worker

FROM alpine:latest AS runtime

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/bin/translation-worker /usr/local/bin/translation-worker

ENV APP_ENV=production

CMD ["translation-worker"]

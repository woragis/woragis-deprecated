FROM golang:1.24-bookworm AS builder

WORKDIR /build

# Enable automatic toolchain download if needed
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./

RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/bin/server ./app/cmd/server

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/bin/server /usr/local/bin/server

ENV APP_NAME=woragis-server \
    APP_ENV=production \
    APP_PORT=8080

EXPOSE 8080

CMD ["server"]


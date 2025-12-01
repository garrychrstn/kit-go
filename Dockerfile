# Stage 1: Build
FROM golang:latest AS builder
WORKDIR /build
COPY go.mod go.sum docker-compose.yaml ./
RUN go mod download
COPY . .
RUN go build -o main src/main.go

# Stage 2: Runtime
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /build/main .
COPY --from=builder /build/db/migrations ./db/migrations
COPY --from=builder /build/docker-compose.yaml ./docker-compose.yaml
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

EXPOSE 7001
CMD ["/app/main"]

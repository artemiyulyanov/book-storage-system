# --- Этап сборки ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# common — сосед book-service, копируем целиком
COPY common ./common

# go.mod/go.sum именно book-service, путь от корня контекста
COPY book-service/go.mod book-service/go.sum* ./book-service/

WORKDIR /app/book-service
RUN go mod download

# Остальные исходники book-service
COPY book-service ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

# --- Финальный этап ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/book-service/server .
COPY --from=builder /app/book-service/migrations ./migrations

EXPOSE 8081

CMD ["./server"]
# --- Этап сборки ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Сначала копируем только go.mod/go.sum, чтобы кэшировать слой с зависимостями
COPY go.mod go.sum* ./
RUN go mod download

# Копируем исходники и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# --- Финальный этап ---
FROM alpine:3.19

# certs нужны, если приложение ходит наружу по https (например к внешним API)
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
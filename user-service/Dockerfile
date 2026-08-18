FROM golang:1.26-alpine AS builder

WORKDIR /app


COPY common ./common


COPY user-service/go.mod user-service/go.sum* ./user-service/

WORKDIR /app/user-service
RUN go mod download


COPY user-service ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd


FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/user-service/server .
COPY --from=builder /app/user-service/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
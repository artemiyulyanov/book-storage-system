FROM golang:1.26-alpine AS builder

WORKDIR /app


COPY common ./common


COPY event-service/go.mod event-service/go.sum* ./event-service/

WORKDIR /app/event-service
RUN go mod download


COPY event-service ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd


FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/event-service/server .
COPY --from=builder /app/event-service/migrations ./migrations

CMD ["./server"]
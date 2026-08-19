FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY common ./common

COPY api-gateway/go.mod api-gateway/go.sum* ./api-gateway/

WORKDIR /app/api-gateway
RUN go mod download


COPY api-gateway ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd


FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/api-gateway/server .

EXPOSE 8080

CMD ["./server"]
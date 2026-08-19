FROM golang:1.26-alpine AS builder

WORKDIR /app


COPY common ./common


COPY auth-service/go.mod auth-service/go.sum* ./auth-service/

WORKDIR /app/auth-service
RUN go mod download


COPY auth-service ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd


FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/auth-service/server .

EXPOSE 8083

CMD ["./server"]
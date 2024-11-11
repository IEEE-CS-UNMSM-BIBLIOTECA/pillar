FROM golang:1.23.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go build -o server ./cmd/serve

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server /app/server

WORKDIR /app

COPY .env .env

EXPOSE 6969

CMD ["./server"]
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o server main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server /app/server

WORKDIR /app
EXPOSE 6969

CMD ["./server"]
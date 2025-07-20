# syntax=docker/dockerfile:1.4
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache openssl

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN chmod +x setup.sh && ./setup.sh

RUN CGO_ENABLED=0 go build -o server ./main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/cert.pem .
COPY --from=builder /app/key.pem .
COPY --from=builder /app/config.yaml .

EXPOSE 4431/udp

ENTRYPOINT ["./server"]

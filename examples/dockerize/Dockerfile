# syntax=docker/dockerfile:1.4
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache openssl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sh setup.sh

RUN CGO_ENABLED=0 go build -o server main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/cert.pem .
COPY --from=builder /app/key.pem .
COPY --from=builder /app/config.yaml .

RUN sysctl -w net.core.rmem_max=7500000 && sysctl -w net.core.wmem_max=7500000 || true

EXPOSE 4431

ENTRYPOINT ["/bin/sh", "-c", "./server & tail -f /dev/null"]

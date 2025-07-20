# syntax=docker/dockerfile:1.4
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy your local module files
COPY go.mod go.sum ./
RUN go mod download

# Copy everything in the dockerize/ dir (including main.go, setup.sh, config.yaml)
COPY . .

# Generate self-signed certs
RUN chmod +x setup.sh && ./setup.sh

# Build the server
RUN CGO_ENABLED=0 go build -o server ./main.go

# Runtime image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/cert.pem .
COPY --from=builder /app/key.pem .
COPY --from=builder /app/config.yaml .

EXPOSE 4431/udp

ENTRYPOINT ["./server"]

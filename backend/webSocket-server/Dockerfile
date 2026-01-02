FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o socket-server .

FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /app/socket-server .
EXPOSE 8081
CMD ["./socket-server"]
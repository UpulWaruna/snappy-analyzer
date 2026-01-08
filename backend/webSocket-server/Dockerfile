FROM golang:1.25-alpine AS builder 
WORKDIR /app

# 1. Copy go.mod files - Added 'backend/' to common/logger
COPY backend/webSocket-server/go.mod backend/webSocket-server/go.sum ./backend/webSocket-server/
COPY backend/common/logger/go.mod ./backend/common/logger/

# 2. Download dependencies
WORKDIR /app/backend/webSocket-server
RUN go mod download

# 3. Copy source code - Added 'backend/' to common/logger
WORKDIR /app
COPY backend/webSocket-server/ ./backend/webSocket-server/
COPY backend/common/logger/ ./backend/common/logger/

# 4. Build the binary
WORKDIR /app/backend/webSocket-server
RUN go build -o socket-server ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/backend/webSocket-server/socket-server .
EXPOSE 8081
CMD ["./socket-server"]
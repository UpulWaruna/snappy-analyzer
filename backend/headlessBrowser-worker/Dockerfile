# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY backend/headlessBrowser-worker/go.mod backend/headlessBrowser-worker/go.sum ./backend/headlessBrowser-worker/
COPY backend/common/logger/go.mod ./backend/common/logger/

WORKDIR /app/backend/headlessBrowser-worker
RUN go mod download

# 3. Copy source code - Added 'backend/' to common/logger
WORKDIR /app
COPY backend/headlessBrowser-worker/ ./backend/headlessBrowser-worker/
COPY backend/common/logger/ ./backend/common/logger/

# 4. Build the binary
WORKDIR /app/backend/headlessBrowser-worker
RUN go build -o worker ./cmd/main.go

# Stage 2: Create the runtime image with Chrome
FROM debian:bullseye-slim

WORKDIR /root/

# Install Google Chrome and dependencies
RUN apt-get update && apt-get install -y \
    wget \
    gnupg \
    ca-certificates \
    --no-install-recommends \
    && wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list' \
    && apt-get update \
    && apt-get install -y google-chrome-stable \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary from the builder stage
COPY --from=builder /app/backend/headlessBrowser-worker/worker .

# Export port
EXPOSE 8080

# Command to run
CMD ["./worker"]
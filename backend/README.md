
# Snappy-Analyzer

A distributed, highly efficient web analyzer built with **Go (Golang)** and **React**. This project uses **Domain-Driven Design (DDD)** to separate concerns and **WebSockets** for real-time data broadcasting.

## 1. High-Level Architecture & Workflow

The application is no longer a monolithic backend. It is split into two specialized services and a shared logging library:

1. **React Frontend**: Sends an analysis request to the Worker.
2. **HeadlessBrowser-Worker**:
* Receives the URL.
* Spawns a **Headless Chrome** instance via `ChromeDP` to render dynamic JavaScript content.
* Parses the rendered HTML and concurrently checks link accessibility.
* Sends the final result to the WebSocket Server.


3. **WebSocket-Server**:
* Maintains a pool of active client connections.
* Receives results from the worker and broadcasts them instantly to the React UI.


4. **Common Logger**: A shared package that provides contextual, scoped logging (user/request groups) across all services.

---

## 2. Refactored Project Structure

Following **Clean Architecture** principles, the code is organized into layers:

```text
/snappy-analyzer
├── /backend
│   ├── /common/logger      # Shared Scoped Logger module
│   ├── /headless-worker    # Headless Chrome & Analysis Service
│   │   ├── /cmd            # Main entry point
│   │   ├── /api/http       # HTTP Handlers (Adapters)
│   │   ├── /application    # Use Cases (Orchestration)
│   │   ├── /domain         # Logic (Parser, Link Checker, Models)
│   │   └── /adapter        # External systems (ChromeDP, API Clients)
│   └── /socket-server      # WebSocket Hub & Broadcaster
│       ├── /cmd            # Main entry point
│       ├── /api/http       # WS and Publish Handlers
│       ├── /application    # Broadcast Use Cases
│       └── /domain/model   # Hub and Client state models
├── /frontend               # React Application
└── docker-compose.yml      # Orchestration for all services

```
## Unit Testing Added for Domain layer 
  /backend/healess-worker/domain/analysis_service_test.go
---

## 3. Key Technical Decisions & Refactors

### 📡 Scoped Contextual Logging

We moved away from basic logs to a **Context-Aware Logging Module**.

* **Minimal Params**: By using `logger.Scoped()`, we initialize a logger once per request with `userGroup` and `requestGroup`. Every subsequent log line automatically includes these IDs without passing them as parameters.
* **Level-Based Source Tracking**: To keep production logs clean, file/line source info is only added to `DEBUG` and `ERROR` logs.

### 🌐 Headless Rendering (ChromeDP)

Standard `http.Get` fails on modern React/SPA sites.

* **Decision**: Switched to `ChromeDP`.
* **Workflow**: The app navigates to the URL, waits for the `body` tag to be visible, and sleeps for 5 seconds to ensure all dynamic elements (like lazy-loaded links) are rendered before parsing.

### 🧵 Semaphore-Controlled Concurrency

Link checking is the primary bottleneck.

* **Semaphore**: We use a buffered channel as a semaphore to limit concurrent link pings to **10**. This prevents the worker from being flagged as a DDoS attack while remaining significantly faster than sequential checking.
* **User-Agent Spoofing**: Added a realistic Chrome User-Agent to prevent `403 Forbidden` responses from sites that block basic Go scrapers.

---

## 4. Setup & Installation

### For New Users

1. **Clone the Repository**:
```bash
git clone https://github.com/UpulWaruna/snappy-analyzer.git

```


2. **Environment Sync**: Ensure your Docker Engine is running. If you encounter "Internal Server Error" from Docker, restart Docker Desktop and run:
```bash
docker builder prune -f

```


3. **Launch with Docker Compose**:
```bash
docker compose up --build --remove-orphans

```


* **Frontend**: `http://localhost:3000`
* **Worker**: `http://localhost:8080`
* **Socket**: `http://localhost:8081`



---

## 5. Important Tunings

* **Docker Networking**: The worker communicates with the socket server using the internal Docker DNS: `http://socket-service:8081/publish`.
* **Sensitive Data**: Use the `logger.Sensitive` type for logging tokens or passwords. It automatically redacts values to `REDACTED` in the JSON output.
* **ChromeDP in Docker**: The Dockerfile uses `debian:bullseye-slim` and installs `google-chrome-stable`. The Go code uses `--no-sandbox` and `--disable-dev-shm-usage` flags, which are mandatory for running Chrome inside a container.

---

## 6. Development Decisions (DDD)

* **Decoupled Domain**: The `domain/service` in the worker contains the `ParseHTML` and `ProcessLinks` logic. It has zero dependencies on HTTP or ChromeDP, making it highly testable.
* **The Hub**: The WebSocket Server uses a centralized `Hub` struct. The `Broadcast` use case runs in a separate goroutine to ensure that one slow client does not block the entire message pipeline.

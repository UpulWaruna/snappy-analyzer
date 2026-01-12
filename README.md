# Web Page Analyzer

A full-stack application built with **Go (Golang)** and **React** that analyzes a given URL for HTML version, headings, link status, and login forms.

## Features
- **Concurrent Link Checking**: Uses Go routines to verify internal/external links simultaneously, significantly reducing processing time.
- **DOM Traversal**: Uses the `golang.org/x/net/html` library to parse the document tree efficiently.
- **Responsive UI**: A modern React dashboard with error handling and loading states.

## Getting Started

### Prerequisites
- Go 1.18+ 
- Node.js 16+

### Running the Backend
1. Navigate to the backend folder: `cd backend`
2. Install dependencies: `go mod tidy`
3. Start the server: `go run .`
4. The server will start on `http://localhost:8080`
5. ### 📂 Detailed Service Documentation
For technical deep-dives, API specifications, and service-level architecture:
👉 **[View Backend Technical Guide](./backend/README.md)**

### Running the Frontend
1. Navigate to the frontend folder: `cd frontend`
2. Install dependencies: `npm install`
3. Start the app: `npm start`
4. The UI will be available at `http://localhost:3000`

## Implementation Details & Decisions

### Concurrency
Checking links is an I/O bound task. I implemented a **Worker Pool pattern using Goroutines** and a `sync.WaitGroup`. This allows the application to check dozens of links in the time it would normally take to check one.

### Login Form Detection
The application identifies a login form by looking for a `<form>` element that contains an `<input type="password">`. This is the most reliable heuristic for detecting authentication points without high-level scrapers.

### Link Resolution
Relative links (e.g., `/about`) are automatically resolved to absolute URLs (e.g., `https://example.com/about`) using Go's `url.ResolveReference` before being checked for accessibility.

## Potential Improvements
1. **Caching**: Implement Redis to store analysis results for frequently requested URLs.
2. **Rate Limiting**: Add a rate limiter to the backend to prevent abuse of the analysis endpoint.
3. **SEO Deep-Dive**: Add checks for OpenGraph tags, Meta descriptions, and Image ALT attributes.
4. **Export**: Allow users to download the analysis report as a PDF or CSV.




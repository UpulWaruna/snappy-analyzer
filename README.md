# snappy-analyzer

Go (Golang) and React,  build a highly efficient, "snappy" analyzer.

Here is a structured roadmap and architectural plan to develop this application.

---

## 1. High-Level Architecture

The application will consist of a React frontend that communicates with a Go REST API. The Go backend will handle the heavy lifting: fetching the HTML, parsing the DOM, and verifying link status.

## 2. Backend Development (Go)

Go is perfect for this because its standard library and concurrency model (goroutines) make checking multiple links extremely fast.

### Core Libraries

* **`net/http`**: For creating the API server and fetching the target URL.
* **`golang.org/x/net/html`**: The standard library for parsing HTML. It’s more robust than regex for finding tags.
* **`goquery`** (Optional): A popular library that provides a jQuery-like syntax for Go, making it much easier to select elements like `<h1>` or `<a>`.

### Logic Steps

1. **URL Validation**: Ensure the user input is a valid URL before attempting to fetch it.
2. **HTML Version**: Check the `<!DOCTYPE>` declaration. (e.g., `<!DOCTYPE html>` signifies HTML5).
3. **Content Analysis**: Use a crawler to count headings (`h1`-`h6`) and find the `<title>`.
4. **Link Analysis**:
* **Internal vs. External**: Check if the `href` starts with `/` or the base domain (Internal) or a different domain (External).
* **Inaccessibility Check**: **This is the bottleneck.** Use **Goroutines** to ping all discovered links concurrently. If you check them one-by-one, the user will wait a long time.


5. **Login Form Detection**: Look for `<form>` tags containing `<input type="password">`.

---

## 3. Frontend Development (React)

The frontend should be a clean, single-page interface.

### Key Components

* **Input Form**: A simple controlled input and a submit button. Disable the button while "loading" to prevent duplicate requests.
* **Results Dashboard**: Use cards or a table to display the metrics.
* **Error Handling**: A clear alert box that displays the HTTP Status Code (e.g., 404 Not Found, 503 Service Unavailable) if the target URL fails.

---

## 4. Suggested Project Structure

Keeping Code organized is vital 

```text
/web-analyzer
├── /backend
│   ├── main.go          # Entry point & Routes
│   ├── parser.go        # Logic for HTML analysis
│   ├── checker.go       # Goroutines for link checking
│   └── models.go        # Structs for JSON response
├── /frontend
│   ├── /src
│   │   ├── App.js       # Main logic and State
│   │   └── /components  # Form.js, Results.js, Error.js
└── README.md            # Your build/deploy instructions

```

---

## 5. Development Decisions & Assumptions

consider including these points:

* **Timeout Policy**: I assumed a 5-second timeout for checking inaccessible links to ensure the application stays responsive.
* **Link Depth**: I only analyzed links found on the immediate page (no deep crawling).
* **Login Detection**: I defined a "Login Form" as any form containing a password field.

## 6. Possible Improvements for the README

suggest these future enhancements:

* **Caching**: Use Redis to cache results for a specific URL for 10 minutes to save bandwidth.
* **SEO Analysis**: Add checks for meta descriptions, image `alt` tags, and mobile responsiveness.
* **Live Updates**: Use WebSockets to stream link-check results to the UI in real-time as they finish, rather than waiting for the whole batch.


* **Models.go**
Why this structure?
HeadingCounts as a Map: Using map[string]int is more flexible than listing H1Count, H2Count, etc., individually. It allows your frontend to simply loop through the keys to display the counts.

Pointer for ErrorDetail: By using a pointer (*ErrorDetail) and the omitempty tag, the JSON will completely hide the error field if the analysis is successful, keeping the API response clean.

LinkStats Sub-struct: Grouping link data makes the JSON easier to read and manage on the React side (e.g., data.links.internal_count).

* **parsar.go**
Key Technical Decisions
Recursive Traversal: The traverse function visits every node in the HTML tree exactly once ($O(n)$ complexity). This is efficient for memory and speed.

Doctype Detection: HTML5 is simply <!DOCTYPE html>. Older versions have long strings (e.g., PUBLIC "-//W3C//DTD HTML 4.01//EN"...).

Login Form Heuristic: The most reliable way to identify a login form without complex AI is checking for the existence of an <input type="password"> inside a <form> element.

Decoupling: Notice that ParseHTML takes an io.Reader. This is a Go best practice because it doesn't care if the HTML comes from a live website, a local file, or a hardcoded string during testing.

* **checker.go**
Key Technical Details

sync.WaitGroup: This acts as a counter. We Add(1) for every link we check and call Done() when finished. wg.Wait() blocks the main execution until the counter hits zero.

sync.Mutex: Since multiple goroutines might try to increment inaccessibleCount at the exact same millisecond, we use a Mutex (Mutual Exclusion) to "lock" the variable during the update to prevent data races.

HEAD vs. GET: We try a HEAD request first. It only fetches the headers, which is much faster than downloading the entire page content. If the server doesn't support HEAD, we fallback to a standard GET.

URL Resolution: We use base.ResolveReference(u) to handle relative links (e.g., <a href="/about">) by turning them into absolute URLs (e.g., https://example.com/about).



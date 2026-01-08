import React, { useState, useEffect, useRef } from 'react';
import './App.css';

function App() {
  const [url, setUrl] = useState('');
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const socketRef = useRef(null);

  // 1. Initialize WebSocket Connection on Mount
  useEffect(() => {
    // Replace with your socket-service URL
    const socket = new WebSocket('ws://localhost:8081/ws');
    socketRef.current = socket;

    socket.onopen = () => console.log('Connected to WebSocket Gateway');
    
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log('Received data via Socket:', data);

      if (data.error) {
        setError(data.error.message);
        setLoading(false);
      } else {
        setResults(data);
        setLoading(false); // Stop loading once the socket pushes data
      }
    };

    socket.onclose = () => console.log('WebSocket Disconnected');
    socket.onerror = (err) => console.error('WebSocket Error:', err);

    return () => {
      if (socketRef.current) socketRef.current.close();
    };
  }, []);

  const handleAnalyze = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResults(null);

    try {
      // Trigger the worker to start analysis
      const response = await fetch('http://localhost:8080/analyze', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error ? data.error.message : "Backend worker error");
      }
      
      // Note: We don't setResults(data) here anymore. 
      // We wait for the WebSocket to push the message.
      console.log("Analysis triggered. Waiting for socket broadcast...");

    } catch (err) {
      setError(err.message || "Could not connect to the backend server.");
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h1>Web Page Analyzer <span className="badge">Real-time</span></h1>
      <form onSubmit={handleAnalyze} className="search-box">
        <input 
          type="url" 
          placeholder="https://example.com" 
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required 
        />
        <button type="submit" disabled={loading}>
          {loading ? "Analyzing..." : "Analyze"}
        </button>
      </form>

      {error && <div className="error-card">⚠️ {error}</div>}

      {/* Progress feedback for the user */}
      {loading && !results && (
        <div className="loading-spinner">
          <p>Browser rendering in progress... Waiting for broadcast.</p>
        </div>
      )}

      {results && (
        <div className="results-grid fade-in">
          <section className="card">
            <h3>General Info</h3>
            <p><strong>URL:</strong> {results.url}</p>
            <p><strong>Title:</strong> {results.page_title || "N/A"}</p>
            <p><strong>Version:</strong> {results.html_version}</p>
            <p><strong>Login Form:</strong> {results.has_login_form ? "✅ Detected" : "❌ Not Found"}</p>
          </section>

          <section className="card">
            <h3>Headings</h3>
            {results.heading_counts && Object.keys(results.heading_counts).length > 0 ? (
              <ul>
                {Object.entries(results.heading_counts).map(([tag, count]) => (
                  <li key={tag}><strong>{tag.toUpperCase()}:</strong> {count}</li>
                ))}
              </ul>
            ) : <p>No headings found.</p>}
          </section>

          <section className="card">
            <h3>Links Analysis</h3>
            <p><strong>Internal:</strong> {results.links.internal_count}</p>
            <p><strong>External:</strong> {results.links.external_count}</p>
            <p className="danger"><strong>Inaccessible:</strong> {results.links.inaccessible}</p>
          </section>
        </div>
      )}
    </div>
  );
}

export default App;
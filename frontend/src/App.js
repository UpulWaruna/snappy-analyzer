import React, { useState } from 'react';
import './App.css';

function App() {
  const [url, setUrl] = useState('');
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleAnalyze = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResults(null);

    try {
      const response = await fetch('http://localhost:8080/analyze', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      });

      const data = await response.json();

      if (!response.ok || data.error) {
        setError(data.error ? data.error.message : `Error: ${response.status}`);
      } else {
        setResults(data);
      }
    } catch (err) {
      setError("Could not connect to the backend server.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h1>Web Page Analyzer</h1>
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

      {error && <div className="error-card">{error}</div>}

      {results && (
        <div className="results-grid">
          <section className="card">
            <h3>General Info</h3>
            <p><strong>Title:</strong> {results.page_title || "N/A"}</p>
            <p><strong>Version:</strong> {results.html_version}</p>
            <p><strong>Login Form:</strong> {results.has_login_form ? "✅ Detected" : "❌ Not Found"}</p>
          </section>

          <section className="card">
            <h3>Headings</h3>
            {Object.keys(results.heading_counts).length > 0 ? (
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

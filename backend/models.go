package main

// AnalysisResult represents the final report sent to the React frontend
type AnalysisResult struct {
	URL           string         `json:"url"`
	HTMLVersion   string         `json:"html_version"`
	PageTitle     string         `json:"page_title"`
	HeadingCounts map[string]int `json:"heading_counts"` // e.g., {"h1": 1, "h2": 5}
	Links         LinkStats      `json:"links"`
	HasLoginForm  bool           `json:"has_login_form"`
	Error         *ErrorDetail   `json:"error,omitempty"`
}

// LinkStats breaks down the internal, external, and broken links
type LinkStats struct {
	InternalCount int `json:"internal_count"`
	ExternalCount int `json:"external_count"`
	Inaccessible  int `json:"inaccessible"`
}

// ErrorDetail provides the specific feedback required if a URL fails
type ErrorDetail struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

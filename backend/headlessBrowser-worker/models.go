package main

// AnalysisResult holds the final data sent to the React frontend
type AnalysisResult struct {
	URL           string         `json:"url"`
	HTMLVersion   string         `json:"html_version"`
	PageTitle     string         `json:"page_title"`
	HeadingCounts map[string]int `json:"heading_counts"`
	Links         LinkStats      `json:"links"`
	HasLoginForm  bool           `json:"has_login_form"`
	Error         *ErrorDetail   `json:"error,omitempty"`
	// unexported field used during processing
	discoveredLinks []string
}

// LinkStats holds the summarized counts
type LinkStats struct {
	InternalCount int `json:"internal_count"`
	ExternalCount int `json:"external_count"`
	Inaccessible  int `json:"inaccessible"`
}

// LinkInfo is what your checker.go specifically is missing
type LinkInfo struct {
	Address    string
	IsExternal bool
	Accessible bool
}

// ErrorDetail for API error responses
type ErrorDetail struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

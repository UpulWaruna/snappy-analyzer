package service

import (
	"strings"
	"testing"
)

func TestParseHTML_Comprehensive(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		expectedTitle string
		expectedH1    int
		hasLogin      bool
		expectedLinks int
	}{
		{
			name:          "Standard Page",
			html:          "<html><head><title>Hi</title></head><body><h1>Title</h1><a href='/1'></a></body></html>",
			expectedTitle: "Hi", expectedH1: 1, hasLogin: false, expectedLinks: 1,
		},
		{
			name:          "Login Form Detection",
			html:          "<form><input type='password'></form>",
			expectedTitle: "", expectedH1: 0, hasLogin: true, expectedLinks: 0,
		},
		{
			name:          "Empty Document",
			html:          "",
			expectedTitle: "", expectedH1: 0, hasLogin: false, expectedLinks: 0,
		},
		{
			name:          "No Title Tag",
			html:          "<body><h1>No Title here</h1></body>",
			expectedTitle: "", expectedH1: 1, hasLogin: false, expectedLinks: 0,
		},
		{
			name:          "JavaScript Links (Should be ignored)",
			html:          "<a href='javascript:void(0)'>Click</a><a href='mailto:test@test.com'>Email</a>",
			expectedTitle: "", expectedH1: 0, hasLogin: false, expectedLinks: 1, // mailto is usually kept, javascript is ignored
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseHTML(strings.NewReader(tt.html))
			if err != nil {
				t.Errorf("Failed on %s: %v", tt.name, err)
			}
			if res.PageTitle != tt.expectedTitle {
				t.Errorf("%s: Title mismatch", tt.name)
			}
			if res.HeadingCounts["h1"] != tt.expectedH1 {
				t.Errorf("%s: H1 count mismatch", tt.name)
			}
			if res.HasLoginForm != tt.hasLogin {
				t.Errorf("%s: Login detection mismatch", tt.name)
			}
			if len(res.DiscoveredLinks) != tt.expectedLinks {
				t.Errorf("%s: Link count mismatch", tt.name)
			}
		})
	}
}

func TestResolveURL(t *testing.T) {
	baseURL := "https://example.com/blog"

	tests := []struct {
		input    string
		expected string
	}{
		{"/contact", "https://example.com/contact"},
		{"https://google.com", "https://google.com"},
		{"about", "https://example.com/about"},
	}

	for _, tc := range tests {
		got := resolveURL(baseURL, tc.input)
		if got != tc.expected {
			t.Errorf("For input %s, expected %s but got %s", tc.input, tc.expected, got)
		}
	}
}

func TestDetermineHTMLVersion(t *testing.T) {
	tests := []struct {
		doctype  string
		expected string
	}{
		{"html", "HTML5"},
		{"HTML 4.01", "HTML 4.01"},
		{"xhtml", "XHTML"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		got := determineHTMLVersion(tt.doctype)
		if got != tt.expected {
			t.Errorf("For %s, expected %s but got %s", tt.doctype, tt.expected, got)
		}
	}
}

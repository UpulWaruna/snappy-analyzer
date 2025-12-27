package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// ParseHTML processes the reader and populates the AnalysisResult
func ParseHTML(body io.Reader) (*AnalysisResult, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	result := &AnalysisResult{
		HeadingCounts: make(map[string]int),
	}

	// Traverse the DOM tree starting from the root
	traverse(doc, result)

	return result, nil
}

func traverse(n *html.Node, res *AnalysisResult) {
	if n.Type == html.DoctypeNode {
		res.HTMLVersion = determineHTMLVersion(n.Data)
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				res.PageTitle = n.FirstChild.Data
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			res.HeadingCounts[n.Data]++
		case "form":
			if isLoginForm(n) {
				res.HasLoginForm = true
			}
		case "a":
			// We will handle links in the next step as they require
			// the base URL to distinguish internal vs external.
		}
	}

	// Recursively visit all children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, res)
	}
}

// determineHTMLVersion checks the doctype string
func determineHTMLVersion(doctype string) string {
	d := strings.ToLower(doctype)
	if d == "html" {
		return "HTML5"
	}
	if strings.Contains(d, "html 4.01") {
		return "HTML 4.01"
	}
	if strings.Contains(d, "xhtml") {
		return "XHTML"
	}
	return "Unknown/Other"
}

// isLoginForm looks for a password input inside a form
func isLoginForm(n *html.Node) bool {
	// Simple logic: If a form contains an <input type="password">, it's a login form
	return hasPasswordInput(n)
}

func hasPasswordInput(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "input" {
		for _, attr := range n.Attr {
			if attr.Key == "type" && attr.Val == "password" {
				return true
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if hasPasswordInput(c) {
			return true
		}
	}
	return false
}

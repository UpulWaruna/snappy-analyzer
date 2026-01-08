package service

import (
	"io"
	"strings"

	"headlessBrowser-worker/domain/model"

	"golang.org/x/net/html"
)

// ParseHTML processes the reader and populates the AnalysisResult
func ParseHTML(body io.Reader) (*model.AnalysisResult, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	result := &model.AnalysisResult{
		HTMLVersion:   "HTML5", // Default fallback for ChromeDP rendered HTML
		HeadingCounts: make(map[string]int),
	}

	// Traverse the DOM tree starting from the root
	traverse(doc, result)

	return result, nil
}

func traverse(n *html.Node, res *model.AnalysisResult) {
	// Handle Doctype Detection
	if n.Type == html.DoctypeNode {
		version := determineHTMLVersion(n.Data)
		if version != "" {
			res.HTMLVersion = version
		}
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				// Clean up tabs and newlines from title
				title := strings.ReplaceAll(n.FirstChild.Data, "\n", "")
				title = strings.ReplaceAll(title, "\t", "")
				res.PageTitle = strings.TrimSpace(n.FirstChild.Data)
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			res.HeadingCounts[n.Data]++
		case "form":
			if isLoginForm(n) {
				res.HasLoginForm = true
			}
		case "a":
			// Extract Links
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := strings.TrimSpace(attr.Val)
					if link != "" && !strings.HasPrefix(link, "javascript:") {
						res.DiscoveredLinks = append(res.DiscoveredLinks, link)
					}
				}
			}
			// Heuristic: If we haven't found a form yet, check if this link looks like a login button
			if !res.HasLoginForm && isLoginLink(n) {
				res.HasLoginForm = true
			}
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
	return ""
}

// isLoginForm looks for a password input inside a form
func isLoginForm(n *html.Node) bool {
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

// isLoginLink checks if an anchor tag text contains "login" (useful for SPAs)
func isLoginLink(n *html.Node) bool {
	if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
		text := strings.ToLower(n.FirstChild.Data)
		if strings.Contains(text, "login") || strings.Contains(text, "sign in") {
			return true
		}
	}
	return false
}

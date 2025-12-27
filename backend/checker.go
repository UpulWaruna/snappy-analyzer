package main

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ProcessLinks analyzes the links found in the HTML
func ProcessLinks(targetURL string, rawLinks []string) LinkStats {
	stats := LinkStats{}
	parsedBase, _ := url.Parse(targetURL)

	uniqueLinks := make(map[string]bool)
	var linksToCheck []string

	for _, link := range rawLinks {
		resolved := resolveLink(parsedBase, link)
		if resolved == "" || uniqueLinks[resolved] {
			continue
		}
		uniqueLinks[resolved] = true
		linksToCheck = append(linksToCheck, resolved)

		// Categorize
		if isExternal(parsedBase, resolved) {
			stats.ExternalCount++
		} else {
			stats.InternalCount++
		}
	}

	// Concurrently check accessibility
	stats.Inaccessible = checkAccessibility(linksToCheck)

	return stats
}

func resolveLink(base *url.URL, href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return base.ResolveReference(u).String()
}

func isExternal(base *url.URL, link string) bool {
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}
	return linkURL.Host != base.Host && linkURL.Host != ""
}

func checkAccessibility(links []string) int {
	var wg sync.WaitGroup
	inaccessibleCount := 0

	// Mutex to safely increment inaccessibleCount from different goroutines
	var mu sync.Mutex

	// Create a client with a timeout so we don't wait forever
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, link := range links {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			// Use HEAD request instead of GET to save bandwidth
			resp, err := client.Head(url)
			if err != nil || resp.StatusCode >= 400 {
				// If HEAD fails, some servers require GET
				resp, err = client.Get(url)
			}

			if err != nil || resp.StatusCode >= 400 {
				mu.Lock()
				inaccessibleCount++
				mu.Unlock()
			}
			if resp != nil {
				resp.Body.Close()
			}
		}(link)
	}

	wg.Wait()
	return inaccessibleCount
}

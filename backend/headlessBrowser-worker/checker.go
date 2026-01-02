package main

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func ProcessLinks(baseURL string, rawLinks []string) []LinkInfo {
	// 1. DEDUPLICATION: Use a map to keep only unique URLs
	uniqueMap := make(map[string]bool)
	var uniqueLinks []string
	for _, l := range rawLinks {
		resolved := resolveURL(baseURL, l)
		if resolved != "" && !uniqueMap[resolved] {
			uniqueMap[resolved] = true
			uniqueLinks = append(uniqueLinks, resolved)
		}
	}

	var wg sync.WaitGroup
	linksChan := make(chan LinkInfo, len(uniqueLinks))
	semaphore := make(chan struct{}, 10) // Reduced to 10 to avoid 403/429 errors

	targetParsed, _ := url.Parse(baseURL)

	for _, linkAddr := range uniqueLinks {
		wg.Add(1)
		go func(resolvedURL string) {
			defer wg.Done()

			// 2. EXTERNAL LOGIC FIX
			isExt := false
			linkParsed, err := url.Parse(resolvedURL)
			if err == nil && linkParsed.Host != "" {
				// Compare hosts. strings.Contains handles elakiri.com vs www.elakiri.com
				isExt = !strings.Contains(linkParsed.Host, targetParsed.Host)
			}

			// 3. CHECK ACCESSIBILITY
			semaphore <- struct{}{}
			accessible := checkLink(resolvedURL)
			<-semaphore

			linksChan <- LinkInfo{
				Address:    resolvedURL,
				IsExternal: isExt,
				Accessible: accessible,
			}
		}(linkAddr)
	}

	go func() {
		wg.Wait()
		close(linksChan)
	}()

	var results []LinkInfo
	for l := range linksChan {
		results = append(results, l)
	}
	return results
}

func checkLink(link string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", link, nil)
	// IMPORTANT: Set User-Agent to prevent the "Inaccessible" 403 errors
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 400
}

func resolveURL(base, link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	baseURL, _ := url.Parse(base)
	return baseURL.ResolveReference(u).String()
}

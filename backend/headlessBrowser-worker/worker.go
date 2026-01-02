package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func GetRenderedHTML(targetURL string) (string, error) {
	// 1. Setup options (Headless mode is default)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,                           // Crucial for Docker
		chromedp.DisableGPU,                          // Usually necessary in containers
		chromedp.Flag("disable-dev-shm-usage", true), // Prevents crashes in small containers
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// 2. Set a generous timeout for Elakiri's heavy load
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	var htmlContent string

	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 5000),
		chromedp.Navigate(targetURL),

		// 1. Wait for the 'body' instead of 'footer' (every site has a body)
		chromedp.WaitVisible(`body`, chromedp.ByQuery),

		// 2. Optional: Scroll to trigger any lazy-loading
		chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),

		// 3. Instead of waiting for a footer, just sleep briefly
		// to let dynamic JS finish loading.
		chromedp.Sleep(5*time.Second),

		chromedp.OuterHTML(`html`, &htmlContent),
	)

	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

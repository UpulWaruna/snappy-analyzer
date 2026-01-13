package analysis

import (
	"bytes"
	"headlessBrowser-worker/application/core"
	"headlessBrowser-worker/domain/model"
	"headlessBrowser-worker/domain/service"
	"log/slog"
)

type AnalyzeURLUseCase struct {
	Browser   core.BrowserProvider
	Publisher core.ResultPublisher
}

func (uc *AnalyzeURLUseCase) Execute(targetURL string, l *slog.Logger) {
	html, err := uc.Browser.GetRenderedHTML(targetURL)
	if err != nil {
		l.Info("failed to get rendered HTML", "error", err.Error())

		// 2. Prepare a specific "Unreachable" error result
		errorResult := model.AnalysisResult{
			URL: targetURL,
			Error: &model.ErrorDetail{
				Message:    "The site could not be reached. Please check the URL and try again.",
				StatusCode: 502,
			},
		}

		// 3. Publish the error result so the UI updates
		if pubErr := uc.Publisher.Publish(errorResult); pubErr != nil {
			l.Error("Failed to publish error state", "error", pubErr)
		}
		return
	}

	result, _ := service.ParseHTML(bytes.NewReader([]byte(html)))
	result.URL = targetURL

	links := service.ProcessLinks(targetURL, result.DiscoveredLinks)
	for _, li := range links {
		if li.IsExternal {
			result.Links.ExternalCount++
		} else {
			result.Links.InternalCount++
		}
		if !li.Accessible {
			result.Links.Inaccessible++
		}
	}
	if err := uc.Publisher.Publish(result); err != nil {
		l.Error("Result delivery failed", "error", err, "url", targetURL)
	} else {
		l.Info("Result delivered to socket server", "url", targetURL)
	}
}

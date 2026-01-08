package analysis

import (
	"bytes"
	"headlessBrowser-worker/application/core"
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
		uc.Publisher.Publish(map[string]string{"url": targetURL, "error": err.Error()})
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

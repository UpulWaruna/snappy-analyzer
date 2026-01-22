package repo

import "headlessBrowser-worker/domain/model"

// AnalysisRepository is the domain interface
type AnalysisRepository interface {
	Save(result *model.AnalysisResult) error
	Get(url string) (*model.AnalysisResult, error)
}

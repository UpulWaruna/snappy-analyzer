package service

import (
	"headlessBrowser-worker/domain/model"
	"headlessBrowser-worker/domain/repo"
)

type AnalysisDataService struct {
	repo repo.AnalysisRepository // Depends on interface, not implementation
}

func (s *AnalysisDataService) SaveAnalysisResult(result *model.AnalysisResult) any {
	panic("unimplemented")
}

func NewAnalysisDataService(r repo.AnalysisRepository) *AnalysisDataService {
	return &AnalysisDataService{repo: r}
}

func (s *AnalysisDataService) StoreAnalysisResult(res *model.AnalysisResult) error {
	return s.repo.Save(res)
}
func (s *AnalysisDataService) RetrieveAnalysisResult(url string) (*model.AnalysisResult, error) {
	return s.repo.Get(url)
}

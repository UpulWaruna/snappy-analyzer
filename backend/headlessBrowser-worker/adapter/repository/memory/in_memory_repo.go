package memory

import (
	"headlessBrowser-worker/domain/model"
	"sync" // Added for thread-safety
)

type InMemoryRepository struct {
	// sync.RWMutex ensures multiple goroutines don't crash the map
	mu   sync.RWMutex
	data map[string]*model.AnalysisResult
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*model.AnalysisResult),
	}
}

func (repo *InMemoryRepository) Save(result *model.AnalysisResult) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.data[result.URL] = result
	return nil
}

func (repo *InMemoryRepository) Get(url string) (*model.AnalysisResult, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	if result, exists := repo.data[url]; exists {
		return result, nil
	}
	return nil, nil
}

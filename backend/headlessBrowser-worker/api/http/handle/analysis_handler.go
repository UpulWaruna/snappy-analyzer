package handle

import (
	"common/logger"
	"encoding/json"
	"headlessBrowser-worker/application/analysis"
	"net/http"
)

type AnalysisHandler struct {
	UseCase *analysis.AnalyzeURLUseCase
}

func (h *AnalysisHandler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	l := logger.Scoped("worker", "system", "req_id", r.URL.Path, r.Method)
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "processing"})

	go h.UseCase.Execute(req.URL, l)
}

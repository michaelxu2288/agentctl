package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/michaelxu2288/swarmboard/internal/orchestration"
)

type Handlers struct {
	Engine *orchestration.Engine
}

func NewHandlers(engine *orchestration.Engine) *Handlers {
	return &Handlers{Engine: engine}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.health)
	mux.HandleFunc("/v1/run", h.runWorkflow)
	mux.HandleFunc("/v1/events", h.events)
}

func (h *Handlers) health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{"ok": true, "time": time.Now().UTC()})
}

func (h *Handlers) runWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	var req struct {
		Goal    string `json:"goal"`
		TraceID string `json:"trace_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	result, err := h.Engine.Run(r.Context(), orchestration.RunInput{Goal: req.Goal, TraceID: req.TraceID})
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handlers) events(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.Engine.Events.History())
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

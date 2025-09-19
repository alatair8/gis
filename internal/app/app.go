package app

import (
	"encoding/json"
	"log"
	"net/http"

	"zemlya-prosto/internal/assistant"
	"zemlya-prosto/internal/catalog"
	"zemlya-prosto/internal/documents"
	"zemlya-prosto/internal/layer"
	"zemlya-prosto/internal/plot"
	"zemlya-prosto/internal/workflow"
)

// App агрегирует доменные сервисы и настраивает HTTP-роуты.
type App struct {
	config          Config
	plotService     plot.Service
	documentService documents.Service
	workflowService workflow.Service
	assistant       assistant.Service
	catalogService  catalog.Service
	layerService    layer.Service
}

// NewApp создаёт приложение с базовыми (пока in-memory) реализациями сервисов.
func NewApp(cfg Config) (*App, error) {
	return &App{
		config:          cfg,
		plotService:     plot.NewInMemoryService(),
		documentService: documents.NewStubService(),
		workflowService: workflow.NewStubService(),
		assistant:       assistant.NewScenarioAssistant(),
		catalogService:  catalog.NewInMemoryService(),
		layerService:    layer.NewStubService(),
	}, nil
}

// RegisterRoutes регистрирует HTTP-обработчики.
func (a *App) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/v1/plots", a.handlePlots)
	mux.HandleFunc("/api/v1/document-packages", a.handleDocumentPackages)
}

// Close освобождает ресурсы приложения.
func (a *App) Close() {
	if err := a.layerService.Close(); err != nil {
		log.Printf("layer service close error: %v", err)
	}
}

func (a *App) handlePlots(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req plot.ContourDraft
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		contour, err := a.plotService.CreateContour(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, contour)
	case http.MethodGet:
		ownerID := r.URL.Query().Get("ownerId")
		contours, err := a.plotService.ListContours(r.Context(), ownerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, contours)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *App) handleDocumentPackages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req documents.PackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	pkg, err := a.documentService.PreparePackage(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := a.workflowService.NotifyPackageReady(r.Context(), pkg.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, pkg)
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

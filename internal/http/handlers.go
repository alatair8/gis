// Package httpapi содержит HTTP-обработчики сервиса.
package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"zemlyaprosto/internal/assistant"
	"zemlyaprosto/internal/model"
	"zemlyaprosto/internal/service"
)

// Handler инкапсулирует работу с HTTP-запросами.
type Handler struct {
	service *service.Service
}

// New создаёт новый набор обработчиков.
func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Register регистрирует маршруты в HTTP-мультиплексоре.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/contours/drawn", h.handleCreateDrawnContour)
	mux.HandleFunc("/api/contours/coordinates", h.handleCreateContourFromCoordinates)
	mux.HandleFunc("/api/contours/import", h.handleImportContour)
	mux.HandleFunc("/api/contours", h.handleListContours)

	mux.HandleFunc("/api/cards", h.handleCreateCard)

	mux.HandleFunc("/api/parcels", h.handleListParcels)

	mux.HandleFunc("/api/document-packages", h.handleDocumentPackages)

	mux.HandleFunc("/api/assistant/suggest", h.handleAssistantSuggest)

	mux.HandleFunc("/api/business/processes", h.handleBusinessProcesses)

	mux.HandleFunc("/api/layer/publish", h.handlePublishLayerFeature)
	mux.HandleFunc("/api/layer", h.handleGetLayer)
}

// writeJSON — вспомогательная функция для формирования ответа.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// writeError упрощает формирование сообщения об ошибке.
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

// handleCreateDrawnContour обрабатывает создание контура методом рисования.
func (h *Handler) handleCreateDrawnContour(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Description string        `json:"description"`
		Points      []model.Point `json:"points"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	contour, err := h.service.CreateContourFromDrawing(req.Description, req.Points)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, contour)
}

// handleCreateContourFromCoordinates создаёт контур на основе координат.
func (h *Handler) handleCreateContourFromCoordinates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Description string        `json:"description"`
		Points      []model.Point `json:"points"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	contour, err := h.service.CreateContourFromCoordinates(req.Description, req.Points)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, contour)
}

// handleImportContour загружает контур из внешней системы.
func (h *Handler) handleImportContour(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Description string        `json:"description"`
		Points      []model.Point `json:"points"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	contour, err := h.service.ImportContour(req.Description, req.Points)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, contour)
}

// handleListContours возвращает все контуры.
func (h *Handler) handleListContours(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	contours := h.service.ListContours()
	writeJSON(w, http.StatusOK, contours)
}

// handleCreateCard формирует информационную карточку.
func (h *Handler) handleCreateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ContourID        string            `json:"contour_id"`
		AutoAttributes   []model.Attribute `json:"auto_attributes"`
		ManualAttributes []model.Attribute `json:"manual_attributes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	card, err := h.service.CreateInformationCard(req.ContourID, req.AutoAttributes, req.ManualAttributes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, card)
}

// handleListParcels возвращает готовые участки.
func (h *Handler) handleListParcels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	category := model.ParcelCategory(r.URL.Query().Get("category"))
	parcels := h.service.ListReadyParcels(category)
	writeJSON(w, http.StatusOK, parcels)
}

// handleDocumentPackages обрабатывает как создание, так и получение пакетов документов.
func (h *Handler) handleDocumentPackages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req struct {
			ContourID string `json:"contour_id"`
			ParcelID  string `json:"parcel_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		pkg, err := h.service.GenerateDocumentPackage(req.ContourID, req.ParcelID)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, pkg)
	case http.MethodGet:
		packages := h.service.GetDocumentPackages()
		writeJSON(w, http.StatusOK, packages)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleAssistantSuggest возвращает подсказки цифрового помощника.
func (h *Handler) handleAssistantSuggest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req assistant.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	suggestions := h.service.GetAssistantSuggestions(req)
	writeJSON(w, http.StatusOK, suggestions)
}

// handleBusinessProcesses обрабатывает создание и изменение процессов.
func (h *Handler) handleBusinessProcesses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// POST /api/business/processes — создание нового процесса
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		process := h.service.CreateBusinessProcess(req.Name)
		writeJSON(w, http.StatusCreated, process)
	case http.MethodPatch:
		// PATCH /api/business/processes?id=...&action=advance
		processID := r.URL.Query().Get("id")
		action := r.URL.Query().Get("action")
		switch action {
		case "advance":
			process, err := h.service.AdvanceBusinessProcess(processID)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusOK, process)
		case "complete":
			stageID := r.URL.Query().Get("stage_id")
			success := strings.ToLower(r.URL.Query().Get("success")) != "false"
			process, err := h.service.CompleteBusinessStage(processID, stageID, success)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusOK, process)
		default:
			writeError(w, http.StatusBadRequest, fmt.Errorf("неизвестное действие: %s", action))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handlePublishLayerFeature публикует контур в слое «Земля просто».
func (h *Handler) handlePublishLayerFeature(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ContourID  string            `json:"contour_id"`
		Attributes map[string]string `json:"attributes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	feature, err := h.service.PublishContourToLayer(req.ContourID, req.Attributes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, feature)
}

// handleGetLayer возвращает слой «Земля просто».
func (h *Handler) handleGetLayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	layer := h.service.GetLayer()
	writeJSON(w, http.StatusOK, layer)
}

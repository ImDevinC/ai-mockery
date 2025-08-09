package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/imdevinc/mockery/services"
)

type MockeryRequest struct {
	Class   string `json:"class"`
	Species string `json:"species"`
}

type MockeryResponse struct {
	Insult string `json:"insult"`
}

type MockeryHandler struct {
	mockeryService *services.MockeryService
}

func NewMockeryHandler(mockeryService *services.MockeryService) *MockeryHandler {
	return &MockeryHandler{
		mockeryService: mockeryService,
	}
}

func (h *MockeryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /mockery", h.handleMockery)
}

func (h *MockeryHandler) handleMockery(w http.ResponseWriter, r *http.Request) {
	var request MockeryRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		slog.Error("failed to decode request body", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	prompt := ""
	if request.Class != "" {
		prompt += "Class: " + request.Class + "\n"
	}
	if request.Species != "" {
		prompt += "Species: " + request.Species + "\n"
	}
	if prompt == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Class or Species must be provided")
		return
	}

	response, err := h.mockeryService.GenerateInsult(r.Context(), prompt)
	if err != nil {
		slog.Error("failed to generate insult", "error", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate insult")
		return
	}
	h.writeJSONResponse(w, http.StatusOK, MockeryResponse{Insult: response})
}

func (h *MockeryHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *MockeryHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

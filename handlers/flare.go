package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/imdevinc/mockery/services"
)

type FlareRequest struct {
	Description       string `json:"description"`
	PreviousResponses string `json:"previousResponses,omitempty"`
}

type FlareResponse struct {
	Responses string `json:"responses"`
}

type FlareHandler struct {
	flareService *services.FlareService
}

func NewFlareHandler(flareService *services.FlareService) *FlareHandler {
	return &FlareHandler{
		flareService: flareService,
	}
}

func (h *FlareHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /flare", h.handleFlare)
}

func (h *FlareHandler) handleFlare(w http.ResponseWriter, r *http.Request) {
	var request FlareRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		slog.Error("failed to decode request body", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	prompt := strings.TrimSpace(request.Description)

	if prompt == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Description must be provided")
		return
	}

	previous, err := base64.StdEncoding.DecodeString(request.PreviousResponses)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid previous responses format")
		return
	}
	var previousResponses []string
	err = json.Unmarshal(previous, &previousResponses)
	if err != nil {
		slog.Error("failed to decode previous responses", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid previous responses format")
		return
	}

	response, err := h.flareService.GenerateFlares(r.Context(), prompt, previousResponses)
	if err != nil {
		slog.Error("failed to generate insult", "error", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate insult")
		return
	}
	h.writeJSONResponse(w, http.StatusOK, FlareResponse{Responses: response})
}

func (h *FlareHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *FlareHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

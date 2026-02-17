package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sudankdk/offense/internal/helper"
	"github.com/sudankdk/offense/internal/tools"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RunTool(w http.ResponseWriter, r *http.Request) {
	toolName := chi.URLParam(r, "tool")
	tools, ok := tools.GetTool(toolName)
	if !ok {
		http.Error(w, "Tool not found", http.StatusNotFound)
		return
	}

	var input map[string]string

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input JSON", http.StatusBadRequest)
		return
	}
	// result, err := tools.Run(r.Context(), input)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	resp := helper.ExecuteResponse(r.Context(), tools, input)

	WriteJson(w, http.StatusOK, resp)

}

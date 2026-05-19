package greeter

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rizky/go-scaffold/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Greet(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		response.Error(w, http.StatusBadRequest, "invalid_request", "name is required")
		return
	}

	greeting, err := h.svc.Greet(r.Context(), name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal_error", "failed to greet")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": greeting})
}

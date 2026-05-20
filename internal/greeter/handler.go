package greeter

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rizky/go-scaffold/internal/middleware"
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
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, greeting)
}

func (h *Handler) CreateGreeting(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		response.Error(w, http.StatusBadRequest, "invalid_request", "name is required")
		return
	}

	greeting, err := h.svc.Greet(r.Context(), name)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, greeting)
}

func RegisterRoutes(r chi.Router, h *Handler, jwtSecret string) {
	r.Get("/greet/{name}", h.Greet)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(jwtSecret))
		r.Use(middleware.RequireRole("admin"))
		r.Post("/greet/{name}", h.CreateGreeting)
	})
}

func handleError(w http.ResponseWriter, err error) {
	switch err {
	case ErrNameRequired:
		response.Error(w, http.StatusBadRequest, "invalid_request", "name is required")
	default:
		response.Error(w, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
	}
}

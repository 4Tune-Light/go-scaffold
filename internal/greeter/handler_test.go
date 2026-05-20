package greeter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rizky/go-scaffold/internal/greeter/dto"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	greetResult *dto.GreetResponse
	greetErr    error
}

func (m *mockService) Greet(ctx context.Context, name string) (*dto.GreetResponse, error) {
	return m.greetResult, m.greetErr
}

func router(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/greet/{name}", h.Greet)
	r.Post("/greet/{name}", h.CreateGreeting)
	return r
}

func TestGreetHandler_Success(t *testing.T) {
	svc := &mockService{greetResult: &dto.GreetResponse{Message: "Hello, John!"}}
	h := NewHandler(svc)
	r := router(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/greet/John", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Hello, John!"`)
}

func TestGreetHandler_ServiceError(t *testing.T) {
	svc := &mockService{greetErr: assert.AnError}
	h := NewHandler(svc)
	r := router(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/greet/John", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateGreetingHandler_Success(t *testing.T) {
	svc := &mockService{greetResult: &dto.GreetResponse{Message: "Hello, Jane!"}}
	h := NewHandler(svc)
	r := router(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/greet/Jane", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"Hello, Jane!"`)
}

func TestCreateGreetingHandler_ServiceError(t *testing.T) {
	svc := &mockService{greetErr: ErrNameRequired}
	h := NewHandler(svc)
	r := router(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/greet/someone", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

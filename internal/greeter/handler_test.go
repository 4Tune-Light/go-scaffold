package greeter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	greetResult string
	greetErr    error
}

func (m *mockService) Greet(ctx context.Context, name string) (string, error) {
	return m.greetResult, m.greetErr
}

func router(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/greet/{name}", h.Greet)
	return r
}

func TestGreetHandler_Success(t *testing.T) {
	svc := &mockService{greetResult: "Hello, John!"}
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

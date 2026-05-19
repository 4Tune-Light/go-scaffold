package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	_ = logger

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test-path", nil)

	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestID_AddsHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, r)

	assert.NotEmpty(t, w.Header().Get("X-Request-Id"))
}

func TestRequestID_PreservesExisting(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Request-Id", "existing-id")

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, r)

	assert.Equal(t, "existing-id", w.Header().Get("X-Request-Id"))
}

func TestRecovery_HandlesPanic(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, strings.TrimSpace(w.Body.String()), "Internal Server Error")
}

func TestCORS_SetsHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("OPTIONS", "/", nil)
	r.Header.Set("Origin", "http://example.com")
	r.Header.Set("Access-Control-Request-Method", "POST")

	handler := CORS([]string{"*"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, r)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

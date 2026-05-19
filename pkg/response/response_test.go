package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusOK, map[string]string{"key": "value"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp API
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
}

func TestJSON_NilData(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusNoContent, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.Bytes())
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusBadRequest, "bad_request", "invalid input")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp API
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "bad_request", resp.Error.Code)
	assert.Equal(t, "invalid input", resp.Error.Message)
}

func TestPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	Paginated(w, []string{"a", "b"}, 1, 10, 2)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp API
	json.Unmarshal(w.Body.Bytes(), &resp)
	meta := resp.Meta.(map[string]interface{})
	assert.Equal(t, float64(1), meta["page"])
	assert.Equal(t, float64(10), meta["per_page"])
	assert.Equal(t, float64(2), meta["total"])
}

func TestCursorPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	CursorPaginated(w, []string{"a"}, "next-cursor-abc", true)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp API
	json.Unmarshal(w.Body.Bytes(), &resp)
	meta := resp.Meta.(map[string]interface{})
	assert.Equal(t, "next-cursor-abc", meta["next_cursor"])
	assert.Equal(t, true, meta["has_more"])
}

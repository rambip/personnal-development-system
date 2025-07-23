package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValuesHandler(t *testing.T) {
	t.Run("GET /values", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/values", nil)
		w := httptest.NewRecorder()

		ValuesHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	})

	t.Run("POST /values", func(t *testing.T) {
		body := strings.NewReader("name=TestValue&description=TestDescription")
		req := httptest.NewRequest(http.MethodPost, "/values", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		ValuesHandler(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status SeeOther, got %v", w.Code)
		}
	})
}

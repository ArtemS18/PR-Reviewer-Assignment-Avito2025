package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserActivity(t *testing.T) {
	r := setupRouter()
	body := `{"user_id":"123","is_active":true}`
	req, _ := http.NewRequest("POST", "/users/setIsActive", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUserActivity_BadRequest(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/users/setIsActive", bytes.NewBufferString("bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserReview(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/users/getReview?user_id=123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserReview_NotFound(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/users/getReview", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

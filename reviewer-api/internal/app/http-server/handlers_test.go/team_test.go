package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTeam(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/team/get?team_name=TestTeam", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTeam_NotFound(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/team/get", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddTeam(t *testing.T) {
	r := setupRouter()
	body := `{"name":"Team","members":[]}`
	req, _ := http.NewRequest("POST", "/team/add", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddTeam_BadRequest(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/team/add", bytes.NewBufferString("bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

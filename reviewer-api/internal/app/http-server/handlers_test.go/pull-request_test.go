package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewPullRequest(t *testing.T) {
	r := setupRouter()
	body := `{"pull_request_id":"1","pull_request_name":"Test PR","author_id":"42"}`
	req, _ := http.NewRequest("POST", "/pullRequest/create", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateNewPullRequest_BadRequest(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/pullRequest/create", bytes.NewBufferString("bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReassignPullRequest(t *testing.T) {
	r := setupRouter()
	body := `{"pull_request_id":"1","old_reviewer_id":"2"}`
	req, _ := http.NewRequest("POST", "/pullRequest/reassign", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReassignPullRequest_BadRequest(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/pullRequest/reassign", bytes.NewBufferString("bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMergedPR(t *testing.T) {
	r := setupRouter()
	body := `{"pull_request_id":"1"}`
	req, _ := http.NewRequest("POST", "/pullRequest/merge", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMergedPR_BadRequest(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/pullRequest/merge", bytes.NewBufferString("bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

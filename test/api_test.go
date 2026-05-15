package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/gin-go-monolithic-boilerplate/internal/app"
	"github.com/example/gin-go-monolithic-boilerplate/internal/config"
)

func newTestServer(t *testing.T) http.Handler {
	t.Helper()

	cfg := config.Default()
	cfg.AppEnv = "test"
	cfg.GinMode = "test"
	cfg.Port = "0"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return app.NewRouter(cfg, logger)
}

func TestRoot(t *testing.T) {
	server := newTestServer(t)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	assertOK(t, response.Body.Bytes())
}

func TestHealth(t *testing.T) {
	server := newTestServer(t)
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	assertOK(t, response.Body.Bytes())
}

func TestCreateAndGetUser(t *testing.T) {
	server := newTestServer(t)

	payload := []byte(`{"name":"Ada Lovelace","email":"ada@example.com"}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(payload))
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()

	server.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, createResponse.Code, createResponse.Body.String())
	}

	var created struct {
		OK   bool `json:"ok"`
		Data struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"data"`
	}
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !created.OK || created.Data.ID == "" || created.Data.Email != "ada@example.com" {
		t.Fatalf("unexpected create response: %+v", created)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+created.Data.ID, nil)
	getResponse := httptest.NewRecorder()

	server.ServeHTTP(getResponse, getRequest)

	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResponse.Code)
	}
	assertOK(t, getResponse.Body.Bytes())
}

func TestDuplicateEmailReturns409(t *testing.T) {
	server := newTestServer(t)
	payload := []byte(`{"name":"Grace Hopper","email":"grace@example.com"}`)

	firstRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(payload))
	firstRequest.Header.Set("Content-Type", "application/json")
	firstResponse := httptest.NewRecorder()
	server.ServeHTTP(firstResponse, firstRequest)

	secondRequest := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(payload))
	secondRequest.Header.Set("Content-Type", "application/json")
	secondResponse := httptest.NewRecorder()
	server.ServeHTTP(secondResponse, secondRequest)

	if secondResponse.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, secondResponse.Code)
	}
}

func TestInvalidJSONReturns400(t *testing.T) {
	server := newTestServer(t)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader([]byte(`{"name":`)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
	assertNotOK(t, response.Body.Bytes())
}

func TestUnknownRouteReturns404(t *testing.T) {
	server := newTestServer(t)
	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
	assertNotOK(t, response.Body.Bytes())
}

func TestUnsupportedMethodReturns405(t *testing.T) {
	server := newTestServer(t)
	request := httptest.NewRequest(http.MethodPost, "/health", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
	assertNotOK(t, response.Body.Bytes())
}

func assertOK(t *testing.T, body []byte) {
	t.Helper()
	var envelope struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if !envelope.OK {
		t.Fatalf("expected ok=true, body=%s", string(body))
	}
}

func assertNotOK(t *testing.T, body []byte) {
	t.Helper()
	var envelope struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if envelope.OK {
		t.Fatalf("expected ok=false, body=%s", string(body))
	}
}

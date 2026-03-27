package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthAlwaysReturns200(t *testing.T) {
	server := newServer(NewTaskService(NewMemoryRepo()))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	server.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content-type, got %q", got)
	}
}

func TestListTasksReturnsEmptyArrayOnFreshServer(t *testing.T) {
	server := newServer(NewTaskService(NewMemoryRepo()))
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()

	server.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content-type, got %q", got)
	}

	var payload struct {
		Data  []Task `json:"data"`
		Count int    `json:"count"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode list payload: %v", err)
	}
	if payload.Data == nil {
		t.Fatal("expected data to be an empty array, got null")
	}
	if len(payload.Data) != 0 {
		t.Fatalf("expected no tasks, got %d", len(payload.Data))
	}
	if payload.Count != 0 {
		t.Fatalf("expected count 0, got %d", payload.Count)
	}
}

func TestCreateTaskRejectsMissingTitle(t *testing.T) {
	server := newServer(NewTaskService(NewMemoryRepo()))
	body := []byte(`{"description":"missing title"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content-type, got %q", got)
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode error payload: %v", err)
	}
	expected := "title is required and must be under 255 characters"
	if payload["error"] != expected {
		t.Fatalf("expected error %q, got %q", expected, payload["error"])
	}
}

func TestUnknownRouteReturnsJSONContentType(t *testing.T) {
	server := newServer(NewTaskService(NewMemoryRepo()))
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	rec := httptest.NewRecorder()

	server.routes().ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content-type, got %q", got)
	}
}

func TestOptionsReturnsJSONContentType(t *testing.T) {
	server := newServer(NewTaskService(NewMemoryRepo()))
	req := httptest.NewRequest(http.MethodOptions, "/api/tasks", nil)
	rec := httptest.NewRecorder()

	server.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content-type, got %q", got)
	}
}

func TestPatchTaskRejectsInvalidUUID(t *testing.T) {
	server := newServer(NewTaskService(NewMemoryRepo()))

	body := []byte(`{"title":"new title"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/tasks/not-a-uuid", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content-type, got %q", got)
	}
	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode error payload: %v", err)
	}
	if payload["error"] != "invalid UUID" {
		t.Fatalf("expected error %q, got %q", "invalid UUID", payload["error"])
	}
}

func TestPatchTaskReturns404ForUnknownUUID(t *testing.T) {
	server := newServer(NewTaskService(NewMemoryRepo()))

	body := []byte(`{"title":"new title"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/tasks/00000000-0000-0000-0000-000000000000", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content-type, got %q", got)
	}
	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode error payload: %v", err)
	}
	if payload["error"] != "Task not found" {
		t.Fatalf("expected error %q, got %q", "Task not found", payload["error"])
	}
}

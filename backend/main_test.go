package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
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

func TestCreateTaskReturnsUUID(t *testing.T) {
	server := newServer(NewTaskService(NewMemoryRepo()))
	body := []byte(`{"title":"ship cli add output"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content-type, got %q", got)
	}

	var payload struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode create payload: %v", err)
	}
	if payload.ID == "" {
		t.Fatal("expected create response to include id")
	}
	if _, err := uuid.Parse(payload.ID); err != nil {
		t.Fatalf("expected create response id to be a UUID, got %q: %v", payload.ID, err)
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

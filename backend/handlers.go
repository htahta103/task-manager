package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type listResponse struct {
	Data  []Task `json:"data"`
	Count int    `json:"count"`
}

type createTaskRequest struct {
	Title       string        `json:"title"`
	Description *string       `json:"description"`
	Status      *TaskStatus   `json:"status"`
	Priority    *TaskPriority `json:"priority"`
	DueDate     *string       `json:"due_date"`
}

type patchTaskRequest struct {
	Title       *string       `json:"title"`
	Description **string      `json:"description"`
	Status      *TaskStatus   `json:"status"`
	Priority    *TaskPriority `json:"priority"`
	DueDate     **string      `json:"due_date"`
}

func (s *apiServer) handleListTasks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var f TaskFilters
	if v := strings.TrimSpace(q.Get("status")); v != "" {
		st := TaskStatus(v)
		f.Status = &st
	}
	if v := strings.TrimSpace(q.Get("priority")); v != "" {
		p := TaskPriority(v)
		f.Priority = &p
	}
	if v := q.Get("search"); v != "" {
		f.Search = &v
	}

	tasks, err := s.svc.List(r.Context(), f)
	if err != nil {
		s.mapError(w, err)
		return
	}
	if tasks == nil {
		tasks = make([]Task, 0)
	}
	writeJSON(w, http.StatusOK, listResponse{Data: tasks, Count: len(tasks)})
}

func (s *apiServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	task, err := s.svc.Create(r.Context(), CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	})
	if err != nil {
		s.mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (s *apiServer) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := s.svc.Get(r.Context(), id)
	if err != nil {
		s.mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (s *apiServer) handlePatchTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	// OpenAPI has minProperties=1; enforce at least one field provided.
	if req.Title == nil && req.Description == nil && req.Status == nil && req.Priority == nil && req.DueDate == nil {
		writeError(w, http.StatusBadRequest, "at least one field must be provided")
		return
	}

	task, err := s.svc.Patch(r.Context(), id, PatchTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	})
	if err != nil {
		s.mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (s *apiServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.Delete(r.Context(), id); err != nil {
		s.mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Task deleted",
		"id":      id,
	})
}

func (s *apiServer) handleClearDoneTasks(w http.ResponseWriter, r *http.Request) {
	deleted, err := s.svc.ClearDone(r.Context())
	if err != nil {
		s.mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Done tasks cleared",
		"deleted": deleted,
	})
}

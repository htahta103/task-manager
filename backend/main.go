package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func main() {
	server := newServer()
	log.Println("task manager API listening on :8080")
	if err := http.ListenAndServe(":8080", server.routes()); err != nil {
		log.Fatal(err)
	}
}

type apiServer struct {
	tasks []task
}

func newServer() *apiServer {
	return &apiServer{tasks: make([]task, 0)}
}

func (s *apiServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /api/tasks", s.handleListTasks)
	mux.HandleFunc("POST /api/tasks", s.handleCreateTask)
	return mux
}

func (s *apiServer) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *apiServer) handleListTasks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data":  s.tasks,
		"count": len(s.tasks),
	})
}

func (s *apiServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" || len(title) > 255 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title is required and must be under 255 characters",
		})
		return
	}

	now := time.Now().UTC()
	item := task{
		ID:        now.Format("20060102150405.000000000"),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.tasks = append(s.tasks, item)
	writeJSON(w, http.StatusCreated, item)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

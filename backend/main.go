package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		port = "8080"
	}

	ctx := context.Background()
	repo, cleanup, err := buildRepository(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	svc := NewTaskService(repo)
	server := newServer(svc)

	log.Printf("task manager API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, server.routes()); err != nil {
		log.Fatal(err)
	}
}

type apiServer struct {
	svc        *TaskService
	corsOrigin string
	authToken  string
}

func newServer(svc *TaskService) *apiServer {
	return &apiServer{
		svc:        svc,
		corsOrigin: strings.TrimSpace(os.Getenv("CORS_ORIGIN")),
		authToken:  strings.TrimSpace(os.Getenv("AUTH_TOKEN")),
	}
}

func (s *apiServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.handleHealth(w, r)
	})

	mux.HandleFunc("/api/tasks/clear/done", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.handleClearDoneTasks(w, r)
	})

	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleListTasks(w, r)
		case http.MethodPost:
			s.handleCreateTask(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
		if strings.TrimSpace(id) == "" || strings.Contains(id, "/") {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		r.SetPathValue("id", id)
		switch r.Method {
		case http.MethodGet:
			s.handleGetTask(w, r)
		case http.MethodPatch:
			s.handlePatchTask(w, r)
		case http.MethodDelete:
			s.handleDeleteTask(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})

	return s.withMiddleware(mux)
}

func (s *apiServer) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func buildRepository(ctx context.Context) (TaskRepository, func(), error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		return NewMemoryRepo(), func() {}, nil
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, func() {}, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, func() {}, err
	}

	migrationsDir := "migrations"
	if err := runMigrations(ctx, pool, migrationsDir); err != nil {
		pool.Close()
		return nil, func() {}, fmt.Errorf("migrations failed: %w", err)
	}

	return NewPostgresRepo(pool), pool.Close, nil
}

func (s *apiServer) mapError(w http.ResponseWriter, err error) {
	switch {
	case err == nil:
		return
	case errors.Is(err, ErrInvalidUUID):
		writeError(w, http.StatusBadRequest, "invalid UUID")
	case IsValidation(err):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, "Task not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

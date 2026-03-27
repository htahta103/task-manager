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
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /api/tasks", s.handleListTasks)
	mux.HandleFunc("POST /api/tasks", s.handleCreateTask)
	mux.HandleFunc("GET /api/tasks/{id}", s.handleGetTask)
	mux.HandleFunc("PATCH /api/tasks/{id}", s.handlePatchTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", s.handleDeleteTask)
	mux.HandleFunc("DELETE /api/tasks/clear/done", s.handleClearDoneTasks)

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

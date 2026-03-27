package main

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type memoryRepo struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

func NewMemoryRepo() TaskRepository {
	return &memoryRepo{
		tasks: make(map[string]Task),
	}
}

func (r *memoryRepo) List(_ context.Context, f TaskFilters) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []Task
	for _, t := range r.tasks {
		if f.Status != nil && t.Status != *f.Status {
			continue
		}
		if f.Priority != nil && t.Priority != *f.Priority {
			continue
		}
		if f.Search != nil {
			needle := strings.ToLower(*f.Search)
			if !strings.Contains(strings.ToLower(t.Title), needle) {
				continue
			}
		}
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (r *memoryRepo) Create(_ context.Context, in CreateTaskInput) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	status := TaskStatusPending
	priority := TaskPriorityMedium
	if in.Status != nil {
		status = *in.Status
	}
	if in.Priority != nil {
		priority = *in.Priority
	}
	id := uuid.NewString()

	var due *time.Time
	var err error
	if in.DueDate != nil {
		due, err = parseISODatePtr(in.DueDate)
		if err != nil {
			return Task{}, ErrValidation("invalid due_date")
		}
	}

	t := Task{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     formatISODatePtr(due),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.tasks[id] = t
	return t, nil
}

func (r *memoryRepo) Get(_ context.Context, id string) (Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return t, nil
}

func (r *memoryRepo) Patch(_ context.Context, id string, in PatchTaskInput) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	if in.Title != nil {
		t.Title = *in.Title
	}
	if in.Description != nil {
		t.Description = *in.Description
	}
	if in.Status != nil {
		t.Status = *in.Status
	}
	if in.Priority != nil {
		t.Priority = *in.Priority
	}
	if in.DueDate != nil {
		var due *time.Time
		var err error
		if *in.DueDate != nil {
			due, err = parseISODatePtr(*in.DueDate)
			if err != nil {
				return Task{}, ErrValidation("invalid due_date")
			}
		}
		t.DueDate = formatISODatePtr(due)
	}
	t.UpdatedAt = time.Now().UTC()

	r.tasks[id] = t
	return t, nil
}

func (r *memoryRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}

func (r *memoryRepo) ClearDone(_ context.Context) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var deleted int64
	for id, t := range r.tasks {
		if t.Status == TaskStatusDone {
			delete(r.tasks, id)
			deleted++
		}
	}
	return deleted, nil
}

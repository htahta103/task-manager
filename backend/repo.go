package main

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type TaskFilters struct {
	Status   *TaskStatus
	Priority *TaskPriority
	Search   *string
}

type TaskRepository interface {
	List(ctx context.Context, f TaskFilters) ([]Task, error)
	Create(ctx context.Context, in CreateTaskInput) (Task, error)
	Get(ctx context.Context, id string) (Task, error)
	Patch(ctx context.Context, id string, in PatchTaskInput) (Task, error)
	Delete(ctx context.Context, id string) error
	ClearDone(ctx context.Context) (int64, error)
}

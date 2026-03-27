package main

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type CreateTaskInput struct {
	Title       string
	Description *string
	Status      *TaskStatus
	Priority    *TaskPriority
	DueDate     *string
}

type PatchTaskInput struct {
	Title       *string
	Description **string
	Status      *TaskStatus
	Priority    *TaskPriority
	DueDate     **string
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) List(ctx context.Context, f TaskFilters) ([]Task, error) {
	if f.Status != nil && !isValidStatus(*f.Status) {
		return nil, ErrValidation("invalid status")
	}
	if f.Priority != nil && !isValidPriority(*f.Priority) {
		return nil, ErrValidation("invalid priority")
	}
	if f.Search != nil && len(*f.Search) > 255 {
		return nil, ErrValidation("search must be under 255 characters")
	}
	return s.repo.List(ctx, f)
}

func (s *TaskService) Create(ctx context.Context, in CreateTaskInput) (Task, error) {
	title := normalizeTitle(in.Title)
	if !validateTitle(title) {
		return Task{}, ErrValidation("title is required and must be under 255 characters")
	}

	if in.Status != nil && !isValidStatus(*in.Status) {
		return Task{}, ErrValidation("invalid status")
	}
	if in.Priority != nil && !isValidPriority(*in.Priority) {
		return Task{}, ErrValidation("invalid priority")
	}
	if in.Description != nil {
		d := strings.TrimSpace(*in.Description)
		in.Description = &d
	}

	in.Title = title
	return s.repo.Create(ctx, in)
}

func (s *TaskService) Get(ctx context.Context, id string) (Task, error) {
	if _, err := uuid.Parse(id); err != nil {
		return Task{}, ErrInvalidUUID
	}
	return s.repo.Get(ctx, id)
}

func (s *TaskService) Patch(ctx context.Context, id string, in PatchTaskInput) (Task, error) {
	if _, err := uuid.Parse(id); err != nil {
		return Task{}, ErrInvalidUUID
	}
	if in.Title != nil {
		t := normalizeTitle(*in.Title)
		if !validateTitle(t) {
			return Task{}, ErrValidation("title is required and must be under 255 characters")
		}
		in.Title = &t
	}
	if in.Status != nil && !isValidStatus(*in.Status) {
		return Task{}, ErrValidation("invalid status")
	}
	if in.Priority != nil && !isValidPriority(*in.Priority) {
		return Task{}, ErrValidation("invalid priority")
	}
	return s.repo.Patch(ctx, id, in)
}

func (s *TaskService) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidUUID
	}
	return s.repo.Delete(ctx, id)
}

func (s *TaskService) ClearDone(ctx context.Context) (int64, error) {
	return s.repo.ClearDone(ctx)
}

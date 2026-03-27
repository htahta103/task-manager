package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) TaskRepository {
	return &postgresRepo{pool: pool}
}

func (r *postgresRepo) List(ctx context.Context, f TaskFilters) ([]Task, error) {
	var where []string
	args := make([]any, 0, 3)
	argn := 1

	if f.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", argn))
		args = append(args, string(*f.Status))
		argn++
	}
	if f.Priority != nil {
		where = append(where, fmt.Sprintf("priority = $%d", argn))
		args = append(args, string(*f.Priority))
		argn++
	}
	if f.Search != nil && strings.TrimSpace(*f.Search) != "" {
		where = append(where, fmt.Sprintf("lower(title) LIKE $%d", argn))
		args = append(args, "%"+strings.ToLower(*f.Search)+"%")
		argn++
	}

	q := `
SELECT id, title, description, status, priority, due_date, created_at, updated_at
FROM tasks
`
	if len(where) > 0 {
		q += "WHERE " + strings.Join(where, " AND ") + "\n"
	}
	q += "ORDER BY created_at ASC"

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var (
			id        string
			title     string
			desc      *string
			status    string
			priority  string
			due       *time.Time
			createdAt time.Time
			updatedAt time.Time
		)
		if err := rows.Scan(&id, &title, &desc, &status, &priority, &due, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		out = append(out, Task{
			ID:          id,
			Title:       title,
			Description: desc,
			Status:      TaskStatus(status),
			Priority:    TaskPriority(priority),
			DueDate:     formatISODatePtr(due),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *postgresRepo) Create(ctx context.Context, in CreateTaskInput) (Task, error) {
	status := string(TaskStatusPending)
	priority := string(TaskPriorityMedium)
	if in.Status != nil {
		status = string(*in.Status)
	}
	if in.Priority != nil {
		priority = string(*in.Priority)
	}
	due, err := parseISODatePtr(in.DueDate)
	if err != nil {
		return Task{}, ErrValidation("invalid due_date")
	}

	var (
		id        string
		title     string
		desc      *string
		st        string
		pr        string
		dueOut    *time.Time
		createdAt time.Time
		updatedAt time.Time
	)
	err = r.pool.QueryRow(ctx, `
INSERT INTO tasks (title, description, status, priority, due_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, description, status, priority, due_date, created_at, updated_at
`, in.Title, in.Description, status, priority, due).Scan(&id, &title, &desc, &st, &pr, &dueOut, &createdAt, &updatedAt)
	if err != nil {
		return Task{}, err
	}
	return Task{
		ID:          id,
		Title:       title,
		Description: desc,
		Status:      TaskStatus(st),
		Priority:    TaskPriority(pr),
		DueDate:     formatISODatePtr(dueOut),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (r *postgresRepo) Get(ctx context.Context, id string) (Task, error) {
	var (
		title     string
		desc      *string
		st        string
		pr        string
		due       *time.Time
		createdAt time.Time
		updatedAt time.Time
	)
	err := r.pool.QueryRow(ctx, `
SELECT title, description, status, priority, due_date, created_at, updated_at
FROM tasks
WHERE id = $1
`, id).Scan(&title, &desc, &st, &pr, &due, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, err
	}
	return Task{
		ID:          id,
		Title:       title,
		Description: desc,
		Status:      TaskStatus(st),
		Priority:    TaskPriority(pr),
		DueDate:     formatISODatePtr(due),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (r *postgresRepo) Patch(ctx context.Context, id string, in PatchTaskInput) (Task, error) {
	// Load existing row first for simplicity and predictable behavior.
	current, err := r.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}

	title := current.Title
	if in.Title != nil {
		title = *in.Title
	}
	desc := current.Description
	if in.Description != nil {
		desc = *in.Description
	}
	status := string(current.Status)
	if in.Status != nil {
		status = string(*in.Status)
	}
	priority := string(current.Priority)
	if in.Priority != nil {
		priority = string(*in.Priority)
	}

	var due *time.Time
	if current.DueDate != nil {
		due, _ = parseISODatePtr(current.DueDate) // current already validated from DB
	}
	if in.DueDate != nil {
		if *in.DueDate == nil {
			due = nil
		} else {
			parsed, err := parseISODatePtr(*in.DueDate)
			if err != nil {
				return Task{}, ErrValidation("invalid due_date")
			}
			due = parsed
		}
	}

	var (
		outID     string
		outTitle  string
		outDesc   *string
		outStatus string
		outPrio   string
		outDue    *time.Time
		createdAt time.Time
		updatedAt time.Time
	)
	err = r.pool.QueryRow(ctx, `
UPDATE tasks
SET title = $2, description = $3, status = $4, priority = $5, due_date = $6
WHERE id = $1
RETURNING id, title, description, status, priority, due_date, created_at, updated_at
`, id, title, desc, status, priority, due).Scan(&outID, &outTitle, &outDesc, &outStatus, &outPrio, &outDue, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, err
	}

	return Task{
		ID:          outID,
		Title:       outTitle,
		Description: outDesc,
		Status:      TaskStatus(outStatus),
		Priority:    TaskPriority(outPrio),
		DueDate:     formatISODatePtr(outDue),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id string) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepo) ClearDone(ctx context.Context) (int64, error) {
	ct, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE status = 'done'`)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}

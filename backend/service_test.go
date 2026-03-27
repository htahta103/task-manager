package main

import (
	"context"
	"errors"
	"testing"
)

type stubRepo struct {
	listFn      func(context.Context, TaskFilters) ([]Task, error)
	createFn    func(context.Context, CreateTaskInput) (Task, error)
	getFn       func(context.Context, string) (Task, error)
	patchFn     func(context.Context, string, PatchTaskInput) (Task, error)
	deleteFn    func(context.Context, string) error
	clearDoneFn func(context.Context) (int64, error)
}

func (r stubRepo) List(ctx context.Context, f TaskFilters) ([]Task, error) {
	return r.listFn(ctx, f)
}
func (r stubRepo) Create(ctx context.Context, in CreateTaskInput) (Task, error) {
	return r.createFn(ctx, in)
}
func (r stubRepo) Get(ctx context.Context, id string) (Task, error) {
	return r.getFn(ctx, id)
}
func (r stubRepo) Patch(ctx context.Context, id string, in PatchTaskInput) (Task, error) {
	return r.patchFn(ctx, id, in)
}
func (r stubRepo) Delete(ctx context.Context, id string) error {
	return r.deleteFn(ctx, id)
}
func (r stubRepo) ClearDone(ctx context.Context) (int64, error) {
	return r.clearDoneFn(ctx)
}

func TestTaskService_Create_ValidatesTitle(t *testing.T) {
	svc := NewTaskService(stubRepo{
		createFn: func(context.Context, CreateTaskInput) (Task, error) {
			t.Fatal("repo.Create should not be called on validation failure")
			return Task{}, nil
		},
	})

	_, err := svc.Create(context.Background(), CreateTaskInput{Title: "   "})
	if err == nil || !IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestTaskService_Create_TrimsTitleAndCallsRepo(t *testing.T) {
	called := false
	svc := NewTaskService(stubRepo{
		createFn: func(_ context.Context, in CreateTaskInput) (Task, error) {
			called = true
			if in.Title != "hello" {
				t.Fatalf("expected trimmed title, got %q", in.Title)
			}
			return Task{ID: "00000000-0000-0000-0000-000000000000", Title: in.Title, Status: TaskStatusPending, Priority: TaskPriorityMedium}, nil
		},
	})

	_, err := svc.Create(context.Background(), CreateTaskInput{Title: "  hello  "})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected repo.Create to be called")
	}
}

func TestTaskService_List_RejectsInvalidFilters(t *testing.T) {
	svc := NewTaskService(stubRepo{
		listFn: func(context.Context, TaskFilters) ([]Task, error) {
			t.Fatal("repo.List should not be called on validation failure")
			return nil, nil
		},
	})

	badStatus := TaskStatus("nope")
	_, err := svc.List(context.Background(), TaskFilters{Status: &badStatus})
	if err == nil || !IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestTaskService_Get_RejectsInvalidUUID(t *testing.T) {
	svc := NewTaskService(stubRepo{
		getFn: func(context.Context, string) (Task, error) {
			t.Fatal("repo.Get should not be called on invalid uuid")
			return Task{}, nil
		},
	})

	_, err := svc.Get(context.Background(), "not-a-uuid")
	var invalidID InvalidIDError
	if !errors.As(err, &invalidID) || invalidID.Problem != "invalid" || invalidID.Param != "id" {
		t.Fatalf("expected InvalidIDError(invalid id), got %v", err)
	}
}

func TestTaskService_Get_RejectsMissingID(t *testing.T) {
	svc := NewTaskService(stubRepo{
		getFn: func(context.Context, string) (Task, error) {
			t.Fatal("repo.Get should not be called on missing id")
			return Task{}, nil
		},
	})

	_, err := svc.Get(context.Background(), "   ")
	var invalidID InvalidIDError
	if !errors.As(err, &invalidID) || invalidID.Problem != "missing" || invalidID.Param != "id" {
		t.Fatalf("expected InvalidIDError(missing id), got %v", err)
	}
}

func TestTaskService_Patch_RejectsInvalidUUID(t *testing.T) {
	svc := NewTaskService(stubRepo{
		patchFn: func(context.Context, string, PatchTaskInput) (Task, error) {
			t.Fatal("repo.Patch should not be called on invalid uuid")
			return Task{}, nil
		},
	})

	_, err := svc.Patch(context.Background(), "not-a-uuid", PatchTaskInput{Title: ptr("x")})
	var invalidID InvalidIDError
	if !errors.As(err, &invalidID) || invalidID.Problem != "invalid" || invalidID.Param != "id" {
		t.Fatalf("expected InvalidIDError(invalid id), got %v", err)
	}
}

func TestTaskService_Patch_ValidatesTitleAndEnums(t *testing.T) {
	goodID := "00000000-0000-0000-0000-000000000000"

	svc := NewTaskService(stubRepo{
		patchFn: func(context.Context, string, PatchTaskInput) (Task, error) {
			t.Fatal("repo.Patch should not be called on validation failure")
			return Task{}, nil
		},
	})

	_, err := svc.Patch(context.Background(), goodID, PatchTaskInput{Title: ptr("   ")})
	if err == nil || !IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}

	badStatus := TaskStatus("nope")
	_, err = svc.Patch(context.Background(), goodID, PatchTaskInput{Status: &badStatus})
	if err == nil || !IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}

	badPriority := TaskPriority("nope")
	_, err = svc.Patch(context.Background(), goodID, PatchTaskInput{Priority: &badPriority})
	if err == nil || !IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestTaskService_Patch_TrimsTitleAndCallsRepo(t *testing.T) {
	goodID := "00000000-0000-0000-0000-000000000000"
	called := false

	svc := NewTaskService(stubRepo{
		patchFn: func(_ context.Context, id string, in PatchTaskInput) (Task, error) {
			called = true
			if id != goodID {
				t.Fatalf("expected id %q, got %q", goodID, id)
			}
			if in.Title == nil || *in.Title != "hello" {
				t.Fatalf("expected trimmed title \"hello\", got %#v", in.Title)
			}
			return Task{ID: id, Title: *in.Title, Status: TaskStatusPending, Priority: TaskPriorityMedium}, nil
		},
	})

	_, err := svc.Patch(context.Background(), goodID, PatchTaskInput{Title: ptr("  hello  ")})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected repo.Patch to be called")
	}
}

func TestTaskService_Create_ValidatesEnumsAndDueDate(t *testing.T) {
	svc := NewTaskService(stubRepo{
		createFn: func(context.Context, CreateTaskInput) (Task, error) {
			t.Fatal("repo.Create should not be called on validation failure")
			return Task{}, nil
		},
	})

	badStatus := TaskStatus("nope")
	_, err := svc.Create(context.Background(), CreateTaskInput{Title: "ok", Status: &badStatus})
	if err == nil || !IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}

	badPriority := TaskPriority("nope")
	_, err = svc.Create(context.Background(), CreateTaskInput{Title: "ok", Priority: &badPriority})
	if err == nil || !IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestTaskService_List_RejectsLongSearch(t *testing.T) {
	svc := NewTaskService(stubRepo{
		listFn: func(context.Context, TaskFilters) ([]Task, error) {
			t.Fatal("repo.List should not be called on validation failure")
			return nil, nil
		},
	})

	long := make([]byte, 256)
	for i := range long {
		long[i] = 'a'
	}
	s := string(long)
	_, err := svc.List(context.Background(), TaskFilters{Search: &s})
	if err == nil || !IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestTaskService_Delete_RejectsInvalidUUID(t *testing.T) {
	svc := NewTaskService(stubRepo{
		deleteFn: func(context.Context, string) error {
			t.Fatal("repo.Delete should not be called on invalid uuid")
			return nil
		},
	})

	err := svc.Delete(context.Background(), "not-a-uuid")
	var invalidID InvalidIDError
	if !errors.As(err, &invalidID) || invalidID.Problem != "invalid" || invalidID.Param != "id" {
		t.Fatalf("expected InvalidIDError(invalid id), got %v", err)
	}
}

func TestTaskService_Delete_RejectsMissingID(t *testing.T) {
	svc := NewTaskService(stubRepo{
		deleteFn: func(context.Context, string) error {
			t.Fatal("repo.Delete should not be called on missing id")
			return nil
		},
	})

	err := svc.Delete(context.Background(), "")
	var invalidID InvalidIDError
	if !errors.As(err, &invalidID) || invalidID.Problem != "missing" || invalidID.Param != "id" {
		t.Fatalf("expected InvalidIDError(missing id), got %v", err)
	}
}

func TestTaskService_ClearDone_PassesThrough(t *testing.T) {
	svc := NewTaskService(stubRepo{
		clearDoneFn: func(context.Context) (int64, error) { return 3, nil },
	})

	n, err := svc.ClearDone(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3, got %d", n)
	}
}

func ptr(s string) *string { return &s }

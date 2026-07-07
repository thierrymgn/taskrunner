package task

import (
	"errors"
	"strings"
	"testing"
)

func TestTaskError_Error(t *testing.T) {
	inner := errors.New("boom")
	err := NewTaskError(CodeExecution, "t1", inner)

	got := err.Error()
	if !strings.Contains(got, "t1") {
		t.Errorf("Error() = %q, doit contenir l'ID de la tâche", got)
	}
	if !strings.Contains(got, "boom") {
		t.Errorf("Error() = %q, doit contenir l'erreur sous-jacente", got)
	}
}

func TestTaskError_UnwrapAndIs(t *testing.T) {
	inner := errors.New("boom")
	err := NewTaskError(CodeTimeout, "t2", inner)

	if !errors.Is(err, inner) {
		t.Errorf("errors.Is doit retrouver l'erreur sous-jacente via Unwrap")
	}
}

func TestTaskError_As(t *testing.T) {
	inner := errors.New("boom")
	wrapped := NewTaskError(CodeConfig, "t3", inner)

	var te *TaskError
	if !errors.As(wrapped, &te) {
		t.Fatalf("errors.As doit extraire un *TaskError")
	}
	if te.TaskID != "t3" {
		t.Errorf("TaskID = %q, want %q", te.TaskID, "t3")
	}
	if te.Code != CodeConfig {
		t.Errorf("Code = %d, want %d", te.Code, CodeConfig)
	}
}

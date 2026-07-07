package task

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFakeTask_Success(t *testing.T) {
	ft := NewFakeTask("f1", BehaviorSuccess, time.Millisecond)
	if err := ft.Execute(context.Background()); err != nil {
		t.Fatalf("attendu nil, got %v", err)
	}
	if ft.ID() != "f1" {
		t.Errorf("ID = %q, want %q", ft.ID(), "f1")
	}
}

func TestFakeTask_Fail(t *testing.T) {
	ft := NewFakeTask("f2", BehaviorFail, time.Millisecond)
	if err := ft.Execute(context.Background()); err == nil {
		t.Fatal("attendu une erreur, got nil")
	}
}

func TestFakeTask_Timeout(t *testing.T) {
	ft := NewFakeTask("f3", BehaviorTimeout, time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := ft.Execute(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("attendu DeadlineExceeded, got %v", err)
	}
}

func TestFakeTask_SuccessCancelledByTimeout(t *testing.T) {
	ft := NewFakeTask("f4", BehaviorSuccess, 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	err := ft.Execute(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("un delay dépassant le timeout doit annuler la tâche, got %v", err)
	}
}

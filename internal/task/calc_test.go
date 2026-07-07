package task

import (
	"context"
	"errors"
	"testing"
)

func TestCalcTask_Sum(t *testing.T) {
	ct := NewCalcTask("c1", 5)
	if err := ct.Execute(context.Background()); err != nil {
		t.Fatalf("attendu nil, got %v", err)
	}
	if ct.Result() != 15 {
		t.Errorf("Result = %d, want 15", ct.Result())
	}
}

func TestCalcTask_ZeroValue(t *testing.T) {
	ct := NewCalcTask("c2", 0)
	if err := ct.Execute(context.Background()); err != nil {
		t.Fatalf("attendu nil, got %v", err)
	}
	if ct.Result() != 0 {
		t.Errorf("Result = %d, want 0", ct.Result())
	}
}

func TestCalcTask_NegativeValueIsConfigError(t *testing.T) {
	ct := NewCalcTask("c3", -1)
	err := ct.Execute(context.Background())

	var te *TaskError
	if !errors.As(err, &te) {
		t.Fatalf("attendu un *TaskError, got %v", err)
	}
	if te.Code != CodeConfig {
		t.Errorf("Code = %d, want CodeConfig (%d)", te.Code, CodeConfig)
	}
}

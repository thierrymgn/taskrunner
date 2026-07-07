package task

import (
	"context"
	"testing"
	"time"
)

func TestScheduled_PromotesTaskAndExposesPolicy(t *testing.T) {
	inner := NewFakeTask("s1", BehaviorSuccess, time.Millisecond)
	s := NewScheduled(inner, 2*time.Second, 3)

	if s.ID() != "s1" {
		t.Errorf("ID promu = %q, want s1", s.ID())
	}
	if err := s.Execute(context.Background()); err != nil {
		t.Errorf("Execute promu doit déléguer à la tâche, got %v", err)
	}

	var p Policy = s
	if p.Timeout() != 2*time.Second {
		t.Errorf("Timeout = %v, want 2s", p.Timeout())
	}
	if p.Retries() != 3 {
		t.Errorf("Retries = %d, want 3", p.Retries())
	}
}

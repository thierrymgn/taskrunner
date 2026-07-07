package orchestrator

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"taskrunner/internal/report"
	"taskrunner/internal/task"
)

func fake(id string, b task.FakeTaskBehavior, delay, timeout time.Duration, retries int) task.Task {
	return task.NewScheduled(task.NewFakeTask(id, b, delay), timeout, retries)
}

func TestRun_AllSuccess(t *testing.T) {
	tasks := []task.Task{
		fake("t1", task.BehaviorSuccess, time.Millisecond, time.Second, 0),
		fake("t2", task.BehaviorSuccess, time.Millisecond, time.Second, 0),
		fake("t3", task.BehaviorSuccess, time.Millisecond, time.Second, 0),
	}

	rep := New(WithWorkers(2)).Run(context.Background(), tasks)
	if len(rep.Results) != 3 {
		t.Fatalf("len = %d, want 3", len(rep.Results))
	}
	for _, r := range rep.Results {
		if r.Status != report.StatusSuccess {
			t.Errorf("%s: status = %q, want success", r.ID, r.Status)
		}
		if r.Attempts != 1 {
			t.Errorf("%s: attempts = %d, want 1", r.ID, r.Attempts)
		}
	}
}

func TestRun_PreservesInputOrder(t *testing.T) {
	tasks := []task.Task{
		fake("a", task.BehaviorSuccess, 30*time.Millisecond, time.Second, 0),
		fake("b", task.BehaviorSuccess, 5*time.Millisecond, time.Second, 0),
		fake("c", task.BehaviorSuccess, 15*time.Millisecond, time.Second, 0),
	}

	rep := New(WithWorkers(3)).Run(context.Background(), tasks)
	want := []string{"a", "b", "c"}
	for i, id := range want {
		if rep.Results[i].ID != id {
			t.Errorf("Results[%d].ID = %q, want %q", i, rep.Results[i].ID, id)
		}
	}
}

func TestRun_RetriesOnFailure(t *testing.T) {
	tasks := []task.Task{
		fake("f", task.BehaviorFail, time.Millisecond, time.Second, 2),
	}

	rep := New().Run(context.Background(), tasks)
	r := rep.Results[0]
	if r.Status != report.StatusFailed {
		t.Errorf("status = %q, want failed", r.Status)
	}
	if r.Attempts != 3 {
		t.Errorf("attempts = %d, want 3 (1 initiale + 2 retries)", r.Attempts)
	}
}

func TestRun_TimeoutStatusAndRetries(t *testing.T) {
	tasks := []task.Task{
		fake("to", task.BehaviorTimeout, 0, 20*time.Millisecond, 1),
	}

	rep := New().Run(context.Background(), tasks)
	r := rep.Results[0]
	if r.Status != report.StatusTimeout {
		t.Errorf("status = %q, want timeout", r.Status)
	}
	if r.Attempts != 2 {
		t.Errorf("attempts = %d, want 2", r.Attempts)
	}
}

func TestRun_VerboseLogsToWriter(t *testing.T) {
	var buf bytes.Buffer
	tasks := []task.Task{
		fake("v1", task.BehaviorSuccess, time.Millisecond, time.Second, 0),
	}

	New(WithVerbose(true), WithLogOutput(&buf)).Run(context.Background(), tasks)
	if !strings.Contains(buf.String(), "v1") {
		t.Errorf("le log verbose doit mentionner la tâche, got %q", buf.String())
	}
}

type concurrencyProbe struct {
	id      string
	current *int32
	max     *int32
	dur     time.Duration
}

func (p *concurrencyProbe) ID() string { return p.id }

func (p *concurrencyProbe) Execute(ctx context.Context) error {
	n := atomic.AddInt32(p.current, 1)
	for {
		old := atomic.LoadInt32(p.max)
		if n <= old || atomic.CompareAndSwapInt32(p.max, old, n) {
			break
		}
	}
	time.Sleep(p.dur)
	atomic.AddInt32(p.current, -1)
	return nil
}

func TestRun_RespectsWorkerLimit(t *testing.T) {
	var current, max int32
	tasks := make([]task.Task, 6)
	for i := range tasks {
		probe := &concurrencyProbe{
			id:      fmt.Sprintf("p%d", i),
			current: &current,
			max:     &max,
			dur:     20 * time.Millisecond,
		}
		tasks[i] = task.NewScheduled(probe, time.Second, 0)
	}

	New(WithWorkers(2)).Run(context.Background(), tasks)
	if got := atomic.LoadInt32(&max); got > 2 {
		t.Errorf("exécutions simultanées max = %d, doit être <= 2", got)
	}
}

func TestRun_ContextCancelProducesPartialReport(t *testing.T) {
	tasks := make([]task.Task, 20)
	for i := range tasks {
		tasks[i] = fake(fmt.Sprintf("s%d", i), task.BehaviorSuccess, 50*time.Millisecond, time.Second, 0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	rep := New(WithWorkers(2)).Run(ctx, tasks)
	if len(rep.Results) >= len(tasks) {
		t.Errorf("rapport partiel attendu, got %d/%d résultats", len(rep.Results), len(tasks))
	}
}

func TestOrchestrate_EntryPoint(t *testing.T) {
	tasks := []task.Task{
		fake("t1", task.BehaviorSuccess, time.Millisecond, time.Second, 0),
	}

	rep, err := Orchestrate(context.Background(), tasks, 2)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(rep.Results) != 1 || rep.Results[0].Status != report.StatusSuccess {
		t.Errorf("rapport inattendu: %+v", rep.Results)
	}
}

func TestOrchestrate_InvalidWorkersFallsBack(t *testing.T) {
	var buf bytes.Buffer
	tasks := []task.Task{
		fake("t1", task.BehaviorSuccess, time.Millisecond, time.Second, 0),
	}

	rep, err := Orchestrate(context.Background(), tasks, 999, WithLogOutput(&buf))
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(rep.Results) != 1 {
		t.Fatalf("len = %d, want 1", len(rep.Results))
	}
	if !strings.Contains(buf.String(), "avertissement") {
		t.Errorf("un avertissement sur le nombre de workers était attendu, got %q", buf.String())
	}
}

package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"taskrunner/internal/report"
	"taskrunner/internal/task"
)

const defaultTimeout = 30 * time.Second

type Orchestrator struct {
	cfg OrchestratorConfig
}

func New(opts ...Option) *Orchestrator {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	workers, err := ValidateWorkers(cfg.workers)
	if err != nil {
		fmt.Fprintf(cfg.logOut, "avertissement: %v (repli sur %d workers)\n", err, workers)
	}
	cfg.workers = workers

	return &Orchestrator{cfg: cfg}
}

func Orchestrate(ctx context.Context, tasks []task.Task, workers int, opts ...Option) (report.Report, error) {
	all := append([]Option{WithWorkers(workers)}, opts...)
	return New(all...).Run(ctx, tasks), nil
}

func (o *Orchestrator) Run(ctx context.Context, tasks []task.Task) report.Report {
	jobs := make(chan task.Task)
	results := make(chan report.TaskResult)

	var wg sync.WaitGroup
	for i := 0; i < o.cfg.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range jobs {
				results <- o.execTask(ctx, t)
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, t := range tasks {
			select {
			case <-ctx.Done():
				return
			case jobs <- t:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	collected := make(map[string]report.TaskResult, len(tasks))
	for r := range results {
		collected[r.ID] = r
	}

	ordered := make([]report.TaskResult, 0, len(collected))
	for _, t := range tasks {
		if r, ok := collected[t.ID()]; ok {
			ordered = append(ordered, r)
		}
	}
	return report.Report{Results: ordered}
}

func (o *Orchestrator) execTask(ctx context.Context, t task.Task) report.TaskResult {
	timeout, retries := policyFor(t)
	start := time.Now()

	var (
		status   report.Status
		attempts int
	)

	for attempt := 1; attempt <= retries+1; attempt++ {
		attempts = attempt

		if ctx.Err() != nil {
			status = report.StatusFailed
			break
		}

		o.logf("[%s] début (tentative %d)", t.ID(), attempt)

		attemptCtx, cancel := context.WithTimeout(ctx, timeout)
		err := t.Execute(attemptCtx)

		timedOut := errors.Is(attemptCtx.Err(), context.DeadlineExceeded)
		cancel()

		switch {
		case err == nil:
			status = report.StatusSuccess
			o.logf("[%s] succès", t.ID())
		case timedOut:
			status = report.StatusTimeout
			o.logf("[%s] timeout", t.ID())
		default:
			status = report.StatusFailed
			o.logf("[%s] échec: %v", t.ID(), err)
		}

		if status == report.StatusSuccess || ctx.Err() != nil {
			break
		}
	}

	return report.TaskResult{
		ID:       t.ID(),
		Status:   status,
		Duration: time.Since(start),
		Attempts: attempts,
	}
}

func policyFor(t task.Task) (time.Duration, int) {
	if p, ok := t.(task.Policy); ok {
		return p.Timeout(), p.Retries()
	}
	return defaultTimeout, 0
}

func (o *Orchestrator) logf(format string, args ...any) {
	if !o.cfg.verbose {
		return
	}
	fmt.Fprintf(o.cfg.logOut, format+"\n", args...)
}

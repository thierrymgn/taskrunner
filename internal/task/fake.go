package task

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type FakeTaskBehavior int

const (
	BehaviorSuccess FakeTaskBehavior = iota
	BehaviorFail
	BehaviorTimeout
)

var errFakeFailure = errors.New("fake task: échec simulé")

func ParseFakeBehavior(s string) (FakeTaskBehavior, error) {
	switch s {
	case "success":
		return BehaviorSuccess, nil
	case "fail":
		return BehaviorFail, nil
	case "timeout":
		return BehaviorTimeout, nil
	default:
		return 0, fmt.Errorf("comportement fake inconnu: %q", s)
	}
}

type FakeTask struct {
	id       string
	behavior FakeTaskBehavior
	delay    time.Duration
}

func NewFakeTask(id string, behavior FakeTaskBehavior, delay time.Duration) *FakeTask {
	return &FakeTask{id: id, behavior: behavior, delay: delay}
}

func (t *FakeTask) ID() string { return t.id }

func (t *FakeTask) Execute(ctx context.Context) error {
	if t.behavior == BehaviorTimeout {
		<-ctx.Done()
		return ctx.Err()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(t.delay):
	}

	if t.behavior == BehaviorFail {
		return errFakeFailure
	}
	return nil
}

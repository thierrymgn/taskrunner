package task

import "time"

type Policy interface {
	Timeout() time.Duration
	Retries() int
}

type Scheduled struct {
	Task
	timeout time.Duration
	retries int
}

func NewScheduled(t Task, timeout time.Duration, retries int) *Scheduled {
	return &Scheduled{Task: t, timeout: timeout, retries: retries}
}

func (s *Scheduled) Timeout() time.Duration { return s.timeout }

func (s *Scheduled) Retries() int { return s.retries }

var (
	_ Task   = (*Scheduled)(nil)
	_ Policy = (*Scheduled)(nil)
)

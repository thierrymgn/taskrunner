package task

import (
	"context"
	"fmt"
)

type CalcTask struct {
	id     string
	value  int
	result int
}

func NewCalcTask(id string, value int) *CalcTask {
	return &CalcTask{id: id, value: value}
}

func (t *CalcTask) ID() string { return t.id }

func (t *CalcTask) Result() int { return t.result }

func (t *CalcTask) Execute(ctx context.Context) error {
	if t.value < 0 {
		return NewTaskError(CodeConfig, t.id, fmt.Errorf("value doit être >= 0, reçu %d", t.value))
	}

	sum := 0
	for i := 1; i <= t.value; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		sum += i
	}
	t.result = sum
	return nil
}

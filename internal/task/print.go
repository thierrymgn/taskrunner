package task

import (
	"context"
	"fmt"
	"io"
	"os"
)

type PrintTask struct {
	id      string
	message string
	out     io.Writer
}

func NewPrintTask(id, message string, out io.Writer) *PrintTask {
	if out == nil {
		out = os.Stdout
	}
	return &PrintTask{id: id, message: message, out: out}
}

func (t *PrintTask) ID() string { return t.id }

func (t *PrintTask) Execute(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err := fmt.Fprintln(t.out, t.message)
	return err
}

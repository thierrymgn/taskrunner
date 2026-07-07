package task

import "fmt"

const (
	CodeUnknown = iota
	CodeConfig
	CodeExecution
	CodeTimeout
)

type TaskError struct {
	Code   int
	TaskID string
	Err    error
}

func NewTaskError(code int, taskID string, err error) *TaskError {
	return &TaskError{Code: code, TaskID: taskID, Err: err}
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("task %q (code %d): %v", e.TaskID, e.Code, e.Err)
}

func (e *TaskError) Unwrap() error {
	return e.Err
}

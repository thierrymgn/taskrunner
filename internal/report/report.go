package report

import (
	"encoding/json"
	"io"
	"time"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusTimeout Status = "timeout"
)

type TaskResult struct {
	ID       string
	Status   Status
	Duration time.Duration
	Attempts int
}

func (r TaskResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID       string `json:"id"`
		Status   Status `json:"status"`
		Duration string `json:"duration"`
		Attempts int    `json:"attempts"`
	}{
		ID:       r.ID,
		Status:   r.Status,
		Duration: r.Duration.String(),
		Attempts: r.Attempts,
	})
}

type Report struct {
	Results []TaskResult `json:"results"`
}

func (r Report) WriteTo(w io.Writer) (int64, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return 0, err
	}
	data = append(data, '\n')

	n, err := w.Write(data)
	return int64(n), err
}

var _ io.WriterTo = Report{}

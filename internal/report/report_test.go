package report

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

func TestTaskResult_DurationFormattedAsString(t *testing.T) {
	res := TaskResult{ID: "t1", Status: StatusSuccess, Duration: 12 * time.Millisecond, Attempts: 1}

	data, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(data), `"duration":"12ms"`) {
		t.Errorf("JSON = %s, doit contenir \"duration\":\"12ms\"", data)
	}
}

func TestReport_WriteTo_RoundTrip(t *testing.T) {
	rep := Report{Results: []TaskResult{
		{ID: "t1", Status: StatusSuccess, Duration: 12 * time.Millisecond, Attempts: 1},
		{ID: "t2", Status: StatusTimeout, Duration: 3001 * time.Millisecond, Attempts: 3},
	}}

	var buf bytes.Buffer
	n, err := rep.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo: %v", err)
	}
	if n != int64(buf.Len()) {
		t.Errorf("n = %d, want %d (nombre d'octets écrits)", n, buf.Len())
	}

	var back struct {
		Results []struct {
			ID       string `json:"id"`
			Status   string `json:"status"`
			Duration string `json:"duration"`
			Attempts int    `json:"attempts"`
		} `json:"results"`
	}
	if err := json.Unmarshal(buf.Bytes(), &back); err != nil {
		t.Fatalf("JSON de sortie invalide: %v", err)
	}
	if len(back.Results) != 2 {
		t.Fatalf("len = %d, want 2", len(back.Results))
	}
	if back.Results[0].ID != "t1" || back.Results[0].Status != "success" || back.Results[0].Attempts != 1 {
		t.Errorf("résultat[0] inattendu: %+v", back.Results[0])
	}
	if back.Results[1].Status != "timeout" || back.Results[1].Duration != "3.001s" {
		t.Errorf("résultat[1] inattendu: %+v", back.Results[1])
	}
}

func TestReport_ImplementsWriterTo(t *testing.T) {
	var _ io.WriterTo = Report{}
}

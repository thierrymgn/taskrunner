package metrics

import (
	"strings"
	"testing"

	"taskrunner/internal/report"
)

func TestWriteMetrics_CountsByStatus(t *testing.T) {
	results := []report.TaskResult{
		{ID: "t1", Status: report.StatusSuccess},
		{ID: "t2", Status: report.StatusTimeout},
		{ID: "t3", Status: report.StatusFailed},
	}

	md := WriteMetrics(results)

	wants := []string{
		"# Métriques d'exécution",
		"Tâches exécutées : 3",
		"Tâches réussies : 1",
		"Tâches en échec : 1",
		"Tâches en timeout : 1",
		"Goroutines actives à la fin :",
	}
	for _, w := range wants {
		if !strings.Contains(md, w) {
			t.Errorf("METRICS.md doit contenir %q, got:\n%s", w, md)
		}
	}
}

func TestWriteMetrics_Empty(t *testing.T) {
	md := WriteMetrics(nil)
	if !strings.Contains(md, "Tâches exécutées : 0") {
		t.Errorf("un rapport vide doit indiquer 0 tâche, got:\n%s", md)
	}
}

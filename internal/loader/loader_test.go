package loader

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"taskrunner/internal/task"
)

func TestParse_ValidFileWithSeveralTypes(t *testing.T) {
	data := []byte(`{
		"tasks": [
			{"id":"t1","type":"print","params":{"message":"hi"},"timeout":"2s","retries":0},
			{"id":"t2","type":"calc","params":{"value":10},"timeout":"1s","retries":1},
			{"id":"t3","type":"fake","params":{"behavior":"fail","delay":"5ms"},"timeout":"500ms","retries":2}
		]
	}`)

	tasks, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("len = %d, want 3", len(tasks))
	}
	if tasks[0].ID() != "t1" {
		t.Errorf("ID[0] = %q, want t1", tasks[0].ID())
	}

	p, ok := tasks[2].(task.Policy)
	if !ok {
		t.Fatalf("la tâche doit implémenter task.Policy")
	}
	if p.Timeout() != 500*time.Millisecond {
		t.Errorf("Timeout = %v, want 500ms", p.Timeout())
	}
	if p.Retries() != 2 {
		t.Errorf("Retries = %d, want 2", p.Retries())
	}
}

func TestParse_UnknownTypeIsConfigError(t *testing.T) {
	data := []byte(`{"tasks":[{"id":"x","type":"???","timeout":"1s"}]}`)

	_, err := Parse(data)
	var te *task.TaskError
	if !errors.As(err, &te) {
		t.Fatalf("attendu un *task.TaskError, got %v", err)
	}
	if te.Code != task.CodeConfig {
		t.Errorf("Code = %d, want CodeConfig (%d)", te.Code, task.CodeConfig)
	}
}

func TestParse_InvalidTimeout(t *testing.T) {
	data := []byte(`{"tasks":[{"id":"x","type":"calc","params":{"value":1},"timeout":"nope"}]}`)

	if _, err := Parse(data); err == nil {
		t.Fatal("attendu une erreur pour un timeout invalide")
	}
}

func TestParse_MissingTimeout(t *testing.T) {
	data := []byte(`{"tasks":[{"id":"x","type":"calc","params":{"value":1}}]}`)

	if _, err := Parse(data); err == nil {
		t.Fatal("attendu une erreur pour un timeout manquant")
	}
}

func TestParse_DuplicateID(t *testing.T) {
	data := []byte(`{"tasks":[
		{"id":"dup","type":"calc","params":{"value":1},"timeout":"1s"},
		{"id":"dup","type":"calc","params":{"value":2},"timeout":"1s"}
	]}`)

	if _, err := Parse(data); err == nil {
		t.Fatal("attendu une erreur pour un id dupliqué")
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	if _, err := Parse([]byte("{not json")); err == nil {
		t.Fatal("attendu une erreur pour un JSON invalide")
	}
}

func TestParse_EmptyTasks(t *testing.T) {
	if _, err := Parse([]byte(`{"tasks":[]}`)); err == nil {
		t.Fatal("attendu une erreur pour une liste vide")
	}
}

func TestLoad_FromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	content := `{"tasks":[{"id":"t1","type":"print","params":{"message":"hi"},"timeout":"1s","retries":0}]}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	tasks, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("len = %d, want 1", len(tasks))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	if _, err := Load("/chemin/inexistant/tasks.json"); err == nil {
		t.Fatal("attendu une erreur pour un fichier absent")
	}
}

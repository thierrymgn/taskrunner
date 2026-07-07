package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"taskrunner/internal/task"
)

type rawFile struct {
	Tasks []rawTask `json:"tasks"`
}

type rawTask struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Params  json.RawMessage `json:"params"`
	Timeout string          `json:"timeout"`
	Retries int             `json:"retries"`
}

func Load(path string) ([]task.Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loader: lecture %s: %w", path, err)
	}
	return Parse(data)
}

func Parse(data []byte) ([]task.Task, error) {
	var file rawFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("loader: JSON invalide: %w", err)
	}
	if len(file.Tasks) == 0 {
		return nil, fmt.Errorf("loader: aucune tâche dans le fichier")
	}

	seen := make(map[string]bool, len(file.Tasks))
	tasks := make([]task.Task, 0, len(file.Tasks))

	for _, rt := range file.Tasks {
		if rt.ID == "" {
			return nil, fmt.Errorf("loader: tâche sans id")
		}
		if seen[rt.ID] {
			return nil, fmt.Errorf("loader: id dupliqué %q", rt.ID)
		}
		seen[rt.ID] = true

		if rt.Retries < 0 {
			return nil, task.NewTaskError(task.CodeConfig, rt.ID, fmt.Errorf("retries doit être >= 0"))
		}

		timeout, err := parseTimeout(rt.Timeout)
		if err != nil {
			return nil, task.NewTaskError(task.CodeConfig, rt.ID, err)
		}

		t, err := buildTask(rt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task.NewScheduled(t, timeout, rt.Retries))
	}
	return tasks, nil
}

func buildTask(rt rawTask) (task.Task, error) {
	switch rt.Type {
	case "print":
		var p struct {
			Message string `json:"message"`
		}
		if err := decodeParams(rt, &p); err != nil {
			return nil, err
		}
		return task.NewPrintTask(rt.ID, p.Message, os.Stderr), nil

	case "calc":
		var p struct {
			Value int `json:"value"`
		}
		if err := decodeParams(rt, &p); err != nil {
			return nil, err
		}
		return task.NewCalcTask(rt.ID, p.Value), nil

	case "download":
		var p struct {
			URL  string `json:"url"`
			Dest string `json:"dest"`
		}
		if err := decodeParams(rt, &p); err != nil {
			return nil, err
		}
		return task.NewDownloadTask(rt.ID, p.URL, p.Dest, nil), nil

	case "fake":
		var p struct {
			Behavior string `json:"behavior"`
			Delay    string `json:"delay"`
		}
		if err := decodeParams(rt, &p); err != nil {
			return nil, err
		}
		behavior, err := task.ParseFakeBehavior(p.Behavior)
		if err != nil {
			return nil, task.NewTaskError(task.CodeConfig, rt.ID, err)
		}
		delay, err := parseDelay(rt.ID, p.Delay)
		if err != nil {
			return nil, err
		}
		return task.NewFakeTask(rt.ID, behavior, delay), nil

	default:
		return nil, task.NewTaskError(task.CodeConfig, rt.ID, fmt.Errorf("type de tâche inconnu: %q", rt.Type))
	}
}

func decodeParams(rt rawTask, dst any) error {
	if len(rt.Params) == 0 {
		return nil
	}
	if err := json.Unmarshal(rt.Params, dst); err != nil {
		return task.NewTaskError(task.CodeConfig, rt.ID, fmt.Errorf("params invalides: %w", err))
	}
	return nil
}

func parseTimeout(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("timeout manquant")
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("timeout invalide %q: %w", s, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("timeout doit être > 0, reçu %q", s)
	}
	return d, nil
}

func parseDelay(id, s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, task.NewTaskError(task.CodeConfig, id, fmt.Errorf("delay invalide %q: %w", s, err))
	}
	return d, nil
}

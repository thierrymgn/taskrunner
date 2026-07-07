package task

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DownloadTask struct {
	id     string
	url    string
	dest   string
	client *http.Client
}

func NewDownloadTask(id, url, dest string, client *http.Client) *DownloadTask {
	if client == nil {
		client = http.DefaultClient
	}
	return &DownloadTask{id: id, url: url, dest: dest, client: client}
}

func (t *DownloadTask) ID() string { return t.id }

func (t *DownloadTask) Execute(ctx context.Context) error {
	if t.url == "" || t.dest == "" {
		return NewTaskError(CodeConfig, t.id, fmt.Errorf("url et dest sont obligatoires"))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.url, nil)
	if err != nil {
		return NewTaskError(CodeConfig, t.id, err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("download %s: %w", t.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewTaskError(CodeExecution, t.id, fmt.Errorf("statut HTTP inattendu: %s", resp.Status))
	}

	f, err := os.Create(t.dest)
	if err != nil {
		return fmt.Errorf("création %s: %w", t.dest, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("écriture %s: %w", t.dest, err)
	}
	return nil
}

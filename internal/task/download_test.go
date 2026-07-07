package task

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDownloadTask_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "payload-data")
	}))
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "out.txt")
	dt := NewDownloadTask("d1", srv.URL, dest, srv.Client())

	if err := dt.Execute(context.Background()); err != nil {
		t.Fatalf("attendu nil, got %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("lecture du fichier: %v", err)
	}
	if string(got) != "payload-data" {
		t.Errorf("contenu = %q, want %q", got, "payload-data")
	}
}

func TestDownloadTask_HTTPErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "out.txt")
	dt := NewDownloadTask("d2", srv.URL, dest, srv.Client())

	err := dt.Execute(context.Background())
	var te *TaskError
	if !errors.As(err, &te) {
		t.Fatalf("attendu un *TaskError, got %v", err)
	}
	if te.Code != CodeExecution {
		t.Errorf("Code = %d, want CodeExecution (%d)", te.Code, CodeExecution)
	}
}

func TestDownloadTask_MissingParamsIsConfigError(t *testing.T) {
	dt := NewDownloadTask("d3", "", "", nil)

	err := dt.Execute(context.Background())
	var te *TaskError
	if !errors.As(err, &te) {
		t.Fatalf("attendu un *TaskError, got %v", err)
	}
	if te.Code != CodeConfig {
		t.Errorf("Code = %d, want CodeConfig (%d)", te.Code, CodeConfig)
	}
}

func TestDownloadTask_ContextTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprint(w, "trop tard")
	}))
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "out.txt")
	dt := NewDownloadTask("d4", srv.URL, dest, srv.Client())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := dt.Execute(ctx); err == nil {
		t.Fatal("un timeout pendant le download doit produire une erreur")
	}
}

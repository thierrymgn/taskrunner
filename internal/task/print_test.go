package task

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestPrintTask_WritesMessage(t *testing.T) {
	var buf bytes.Buffer
	pt := NewPrintTask("p1", "hello", &buf)

	if err := pt.Execute(context.Background()); err != nil {
		t.Fatalf("attendu nil, got %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "hello" {
		t.Errorf("sortie = %q, want %q", got, "hello")
	}
	if pt.ID() != "p1" {
		t.Errorf("ID = %q, want %q", pt.ID(), "p1")
	}
}

func TestPrintTask_RespectsCancelledContext(t *testing.T) {
	var buf bytes.Buffer
	pt := NewPrintTask("p2", "hello", &buf)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := pt.Execute(ctx); err == nil {
		t.Fatal("un contexte déjà annulé doit produire une erreur")
	}
	if buf.Len() != 0 {
		t.Errorf("rien ne doit être écrit si le contexte est annulé, got %q", buf.String())
	}
}

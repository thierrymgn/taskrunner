package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"taskrunner/internal/loader"
	"taskrunner/internal/metrics"
	"taskrunner/internal/orchestrator"
)

const metricsFile = "METRICS.md"

func main() {
	file := flag.String("file", "", "chemin du fichier JSON de tâches (obligatoire)")
	workers := flag.Int("workers", 3, "nombre de workers simultanés")
	verbose := flag.Bool("verbose", false, "affiche le statut des tâches en temps réel sur stderr")
	flag.Parse()

	if *file == "" {
		fmt.Fprintln(os.Stderr, "erreur: le flag -file est obligatoire")
		flag.Usage()
		os.Exit(2)
	}

	tasks, err := loader.Load(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erreur de chargement: %v\n", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rep, err := orchestrator.Orchestrate(ctx, tasks, *workers, orchestrator.WithVerbose(*verbose))
	if err != nil {
		fmt.Fprintf(os.Stderr, "erreur d'orchestration: %v\n", err)
	}

	if _, err := rep.WriteTo(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "erreur d'écriture du rapport: %v\n", err)
	}

	if err := os.WriteFile(metricsFile, []byte(metrics.WriteMetrics(rep.Results)), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "erreur d'écriture de %s: %v\n", metricsFile, err)
	}
}

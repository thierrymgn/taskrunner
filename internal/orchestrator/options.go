package orchestrator

import (
	"io"
	"os"
)

type OrchestratorConfig struct {
	workers int
	verbose bool
	logOut  io.Writer
}

type Option func(*OrchestratorConfig)

func WithWorkers(n int) Option {
	return func(c *OrchestratorConfig) { c.workers = n }
}

func WithVerbose(v bool) Option {
	return func(c *OrchestratorConfig) { c.verbose = v }
}

func WithLogOutput(w io.Writer) Option {
	return func(c *OrchestratorConfig) { c.logOut = w }
}

func defaultConfig() OrchestratorConfig {
	return OrchestratorConfig{
		workers: defaultWorkers,
		verbose: false,
		logOut:  os.Stderr,
	}
}

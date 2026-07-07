package orchestrator

import "fmt"

const (
	minWorkers     = 1
	maxWorkers     = 100
	defaultWorkers = 3
)

func ValidateWorkers(n int) (int, error) {
	if n < minWorkers || n > maxWorkers {
		return defaultWorkers, fmt.Errorf("workers doit être entre %d et %d, reçu %d", minWorkers, maxWorkers, n)
	}
	return n, nil
}

# taskrunner

Orchestrateur de tâches concurrentes en Go (bibliothèque standard uniquement).
Lit un fichier JSON de tâches, les exécute en parallèle via un pool de workers
avec timeout et retries par tâche, puis produit un rapport JSON.

## Utilisation

```bash
make build
./bin/taskrunner -file tasks.json -workers 3
```

| Flag | Défaut | Rôle |
|---|---|---|
| `-file` | *(obligatoire)* | chemin du fichier JSON de tâches |
| `-workers` | `3` | nombre de workers simultanés (borné à `[1, 100]`) |
| `-verbose` | `false` | affiche le statut des tâches en temps réel sur stderr |

- Le **rapport JSON** est écrit sur **stdout**.
- Le fichier **`METRICS.md`** est généré à la fin (goroutines actives, totaux par statut).
- `Ctrl+C` (SIGINT) arrête proprement : les tâches en cours sont annulées et le
  rapport partiel est quand même produit.

## Types de tâches

| type | params | action |
|---|---|---|
| `print` | `message` | affiche un message |
| `calc` | `value` | somme de `1..value` |
| `download` | `url`, `dest` | télécharge une URL vers un fichier |
| `fake` | `behavior` (`success`/`fail`/`timeout`), `delay` | simule une tâche (tests) |

Exemple minimal :

```json
{
  "tasks": [
    { "id": "t1", "type": "print", "params": { "message": "hello" }, "timeout": "2s", "retries": 0 }
  ]
}
```

## Architecture

```
cmd/taskrunner/       flags + SIGINT + appel Orchestrate
internal/task/        interface Task, TaskError, implémentations, décorateur Scheduled
internal/loader/      parsing tasks.json (switch sur "type")
internal/report/      Report (io.WriterTo) + TaskResult
internal/orchestrator/ worker pool, timeout/retries, functional options
internal/metrics/     génération de METRICS.md
```

Le timeout et les retries voyagent avec chaque tâche via le décorateur
`Scheduled` (interface optionnelle `Policy`), ce qui garde l'interface `Task`
minimale tout en respectant la signature imposée `Orchestrate(ctx, []task.Task, workers)`.

## Développement

```bash
make test    # go test ./...
make lint    # go vet + gofmt
make run     # build + exécution sur tasks.json
```

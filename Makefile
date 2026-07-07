BINARY := taskrunner
CMD    := ./cmd/taskrunner

.PHONY: build test run lint

build:
	go build -o bin/$(BINARY) $(CMD)

test:
	go test ./...

run: build
	./bin/$(BINARY) -file tasks.json

lint:
	go vet ./...
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "gofmt: fichiers non formatés:"; echo "$$out"; exit 1; fi

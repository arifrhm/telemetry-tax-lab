.PHONY: all build run test ci bench bench-changelog clean

BINARY_NAME=bin/server

all: build test

build:
	@mkdir -p bin
	go build -o $(BINARY_NAME) cmd/server/main.go
	@echo "✅ Server built: $(BINARY_NAME)"

run: build
	@PORT=8081 ./$(BINARY_NAME)

test:
	go test -v -race ./internal/...

ci:
	@bash scripts/ci_test.sh

bench:
	@bash scripts/bench.sh

bench-changelog:
	@bash scripts/bench.sh --append-changelog $(VERSION) "$(TITLE)"

clean:
	rm -rf bin/ docs/flamegraphs/*.out
	@echo "🧹 Cleaned build artifacts"

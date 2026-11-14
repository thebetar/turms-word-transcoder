BINARY         = bin/turms-word-transcoder
MAIN_PACKAGE   = cmd/turms-word-transcoder
COVERAGE_FILE  = coverage.out

.PHONY: run test build coverage clean

run:
	go run ./$(MAIN_PACKAGE) $(ARGS)

build:
	mkdir -p bin
	go build -o $(BINARY) ./$(MAIN_PACKAGE)

test:
	go test ./...

coverage:
	go test ./internal/... -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic
	go tool cover -func=$(COVERAGE_FILE)

clean:
	rm -rf bin $(COVERAGE_FILE)

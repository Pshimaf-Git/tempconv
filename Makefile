BINARY_NAME=./out/bin/tempconv.exe
MAIN_PKG=./main.go
CMD_DIR=./cmd
COVERAGE_DIR=./out/coverage
COVERAGE_OUT=$(COVERAGE_DIR)/coverage.out
COVERAGE_HTML=$(COVERAGE_DIR)/coverage.html

.PHONY: all build test run clean install fmt lint vet cover help

all: fmt vet test build

build:
	go build -o $(BINARY_NAME) $(MAIN_PKG)

test:
	go test -v $(CMD_DIR)/...

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME) $(COVERAGE_OUT) $(COVERAGE_HTML)

install:
	go install $(MAIN_PKG)

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

cover:
	if not exist $(COVERAGE_DIR) mkdir $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_OUT) $(CMD_DIR)/...
	go tool cover -func=$(COVERAGE_OUT)
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  test     - Run tests"
	@echo "  run      - Build and run the binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install binary to GOPATH/bin"
	@echo "  fmt      - Format code"
	@echo "  vet      - Run go vet"
	@echo "  lint     - Run golangci-lint (requires installation)"
	@echo "  cover    - Generate test coverage report"

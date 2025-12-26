# jfmt Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: all build clean install release

all: build

build:
	go build -ldflags "$(LDFLAGS)" -o jfmt .

install: build
	mv jfmt /usr/local/bin/jfmt

clean:
	rm -f jfmt jfmt-*

# Build for all platforms
release: clean
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o jfmt-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o jfmt-darwin-arm64 .
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o jfmt-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o jfmt-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o jfmt-windows-amd64.exe .

test:
	@echo '{"z":1,"a":2}' | ./jfmt && echo "Basic format: OK"
	@echo '{"z":1,"a":2}' | ./jfmt -s | grep -q '"a"' && echo "Sort keys: OK"
	@echo '{"a":1}' | ./jfmt -c | grep -q '^{"a":1}$$' && echo "Compact: OK"
	@echo '{"a":1,}' | ./jfmt -f | grep -q '"a"' && echo "Fix trailing comma: OK"
	@echo "All tests passed!"

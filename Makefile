DISABLED_LINTERS="depguard,exhaustruct"

all: build

build:
	find cmd -name "*.go" -print0 | xargs -0 -n1 -t go build -o bin/

test:
	go test -v -count=1 ./...

fmt:
	go fmt ./...

lint:
	go vet ./...
	golangci-lint run --enable-all --disable=$(DISABLED_LINTERS)

clean:
	rm -rf bin

generate:

.PHONY: all example test fmt lint clean generate build

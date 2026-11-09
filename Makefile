.PHONY: build test

build:
	go build -o ~/bin/stuntdb ./cmd/stuntdb

test:
	go test ./...

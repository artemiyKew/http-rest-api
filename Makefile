.PHONY: build

build: 
	go build -v ./cmd/apiserver

.PHONY: test

run: build
	./apiserver

test: 
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL = build
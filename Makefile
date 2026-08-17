.DEFAULT_GOAL := run
.PHONY: fmt vet build clean run

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build

run: build
	./book-storage-system
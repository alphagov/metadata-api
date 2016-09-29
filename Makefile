.PHONY: build run test clean

BINARY ?= $(PWD)/metadata-api

all: clean build test

clean:
	rm -rf $(BINARY)

build:
	go build -o $(BINARY)

test: build
	go test -race -v $$(go list ./... | grep -v '/vendor/')

run: build
	$(BINARY)

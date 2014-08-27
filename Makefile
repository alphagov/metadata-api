.PHONY: godep deps test build

BINARY := metadata_api

all: godep deps test build

godep:
	go get -v -u github.com/tools/godep

deps: godep
	godep get

test: deps
	godep go test -v ./...

build: deps
	godep go build -v -o $(BINARY)

clean:
	@rm -f $(BINARY)

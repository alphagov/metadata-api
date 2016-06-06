.PHONY: deps test build

all: test build

test:
	go test -v $$(go list ./... | grep -v '/vendor/')

build:
	go build -o metadata-api

run: build
	./metadata-api

clean:
	rm -rf bin metadata-api

.PHONY: deps test build

all: test build

test:
	go test -v ./ \
		./content_api \
		./need_api \
		./performance_platform \
		./request

build:
	go build -o metadata-api

run: build
	./metadata-api

clean:
	rm -rf bin metadata-api

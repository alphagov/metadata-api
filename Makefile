.PHONY: deps test build

BINARY := metadata-api
ORG_PATH := github.com/alphagov
REPO_PATH := $(ORG_PATH)/$(BINARY)

all: test build

deps:
	gom install

vendor: deps
	mkdir -p vendor

test: vendor
	gom test -v $(REPO_PATH) \
		$(REPO_PATH)/content_api \
		$(REPO_PATH)/need_api \
		$(REPO_PATH)/performance_platform \
		$(REPO_PATH)/request

build: vendor
	gom build -o $(BINARY)

run: build
	./metadata-api

clean:
	rm -rf bin metadata-api vendor

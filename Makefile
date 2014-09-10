.PHONY: deps test build rm_compiled_self

BINARY := metadata-api
ORG_PATH := github.com/alphagov
REPO_PATH := $(ORG_PATH)/$(BINARY)

all: deps test build

deps: third_party/src/$(REPO_PATH) rm_compiled_self
	go run third_party.go get -t -v .

rm_compiled_self:
	rm -rf third_party/pkg/*/$(REPO_PATH)

third_party/src/$(REPO_PATH):
	mkdir -p third_party/src/$(ORG_PATH)
	ln -s ../../../.. third_party/src/$(REPO_PATH)

test: deps
	go run third_party.go test -v $(REPO_PATH) \
		$(REPO_PATH)/content_api \
		$(REPO_PATH)/need_api \
		$(REPO_PATH)/performance_platform \
		$(REPO_PATH)/request

build: deps
	go run third_party.go build -o $(BINARY)

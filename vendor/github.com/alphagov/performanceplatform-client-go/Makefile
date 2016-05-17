.PHONY: deps fmt test

IMPORT_BASE := github.com/alphagov
IMPORT_PATH := $(IMPORT_BASE)/performanceplatform-client.go

all: deps _vendor fmt test

deps:
	go get github.com/mattn/gom
	go get github.com/onsi/ginkgo/ginkgo
	go get golang.org/x/tools/cmd/cover

fmt:
	gofmt -w=1 *.go

test:
	gom exec ginkgo -cover . 
	# rewrite the generated .coverprofile files so that you can run the command
	# gom tool cover -html=./pkg/handlers/handlers.coverprofile and other lovely stuff
	find . -name '*.coverprofile' -type f -exec sed -i '' 's|_'$(CURDIR)'|\.|' {} \;

_vendor: Gomfile _vendor/src/$(IMPORT_PATH)
	gom -test install
	touch _vendor

_vendor/src/$(IMPORT_PATH):
	rm -f _vendor/src/$(IMPORT_PATH)
	mkdir -p _vendor/src/$(IMPORT_BASE)
	ln -s $(CURDIR) _vendor/src/$(IMPORT_PATH)

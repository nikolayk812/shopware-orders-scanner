.PHONY: clean test build.local build.linux build.osx build.docker

BINARY        ?= shopware-orders-scanner
VERSION       ?= $(shell git describe --tags --always --dirty)
IMAGE         ?= $(BINARY)
SOURCES       = $(shell find . -name '*.go')
BUILD_FLAGS   ?=-mod=readonly

help:
	@grep -E '^[a-zA-Z0-9_.$$()-/%]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

default: build.local

clean: ## Clean the existing build files
	rm -rf build

test.unit: ## Run unit test
	go test -race  -cover -p 8  ./...

build.local: build/$(BINARY)
build.linux: build/linux/$(BINARY)
build.osx: build/osx/$(BINARY)

build/$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS)  app.go

build/linux/$(BINARY): $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/linux/$(BINARY) app.go

build/osx/$(BINARY): $(SOURCES)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/osx/$(BINARY) app.go

build.docker: build.linux ## Build local docker image
	docker build --rm -t "$(IMAGE)" -f Dockerfile .

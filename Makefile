DOCKER = $(shell which docker)

ifeq ($(DOCKER), "docker not found")
	$(error Docker is needed for build, exiting.)
endif

ifeq ($(OS), Windows_NT)
	FILE_PREFIX = ".exe"
endif
OUTPUT_PATH = "build/micro-ddns$(FILE_PREFIX)"

GO ?= $(shell which go)
ifeq ($(GO), "go not found")
	$(error Golang SDK not detected, exiting.)
endif

VERSION ?= "0.0.1"
IMG ?= "docker.io/masteryyh/micro-ddns"
TAG ?= $(VERSION)

BUILD_TIME = $(shell date --iso=seconds)

GO_VERSION = "go1.22.6"

COMMIT_HASH = $(shell git rev-parse HEAD)

LDFLAGS = "-X 'github.com/masteryyh/micro-ddns/internal/version.Version=$(VERSION)' \
			-X 'github.com/masteryyh/micro-ddns/internal/version.BuildTime=$(BUILD_TIME)' \
			-X 'github.com/masteryyh/micro-ddns/internal/version.GoVersion=$(GO_VERSION)' \
			-X 'github.com/masteryyh/micro-ddns/internal/version.CommitHash=$(COMMIT_HASH)'"

all: clean build build-image

clean:
	@echo "Cleaning build artifacts and images..."
	rm -rf build
	$(DOCKER) image rm -f $(IMG):$(TAG)
	$(DOCKER) image rm -f $(IMG):latest

build:
	@echo "Building project binary..."
	@if [ ! -f $(OUTPUT_PATH) ]; then \
		$(GO) build -ldflags $(LDFLAGS) -o $(OUTPUT_PATH) cmd/main.go ; \
	else \
	  	echo "Already built, skipping binary build."; \
	fi

build-image: build
	@echo "Building Docker images..."
	$(DOCKER) build -t $(IMG):$(TAG) .
	$(DOCKER) tag $(IMG):$(TAG) $(IMG):latest
	@echo "Pushing images to Docker Hub..."
	$(DOCKER) push $(IMG):$(TAG)
	$(DOCKER) push $(IMG):latest


.PHONY: all clean build build-image

DOCKER = $(shell which docker)

ifeq ($(DOCKER), "docker not found")
	$(error Docker is needed for build, exiting.)
endif

OUTPUT_PATH = "build/micro-ddns"

GO ?= $(shell which go)
ifeq ($(GO), "go not found")
	$(error Golang SDK not detected, exiting.)
endif

USE_CN_MIRROR ?= ""

VERSION ?= "0.0.1"
IMG ?= "docker.io/masteryyh/micro-ddns"
TAG ?= $(VERSION)

BUILD_TIME = $(shell date --iso=seconds)

GO_VERSION = "go1.23.0"

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
	$(DOCKER) build -t $(IMG):$(TAG) --build-arg USE_CN_MIRROR=$(USE_CN_MIRROR) .
	$(DOCKER) tag $(IMG):$(TAG) $(IMG):latest
	@echo "Pushing images to Docker Hub..."
	$(DOCKER) push $(IMG):$(TAG)
	$(DOCKER) push $(IMG):latest


.PHONY: all clean build build-image

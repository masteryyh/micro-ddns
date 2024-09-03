DOCKER = $(shell which docker)

OUTPUT_PATH = bin/micro-ddns

GO ?= $(shell which go)
ifeq ($(GO), "go not found")
	$(error Golang SDK not detected, exiting.)
endif

USE_CN_MIRROR ?= ""

VERSION ?= 0.0.1
IMG ?= docker.io/masteryyh/micro-ddns
TAG ?= $(VERSION)

BUILD_TIME = $(shell date --iso=seconds)

GO_VERSION = "go1.23.0"

COMMIT_HASH = $(shell git rev-parse HEAD)

LDFLAGS = "-X 'github.com/masteryyh/micro-ddns/internal/version.Version=$(VERSION)' \
			-X 'github.com/masteryyh/micro-ddns/internal/version.BuildTime=$(BUILD_TIME)' \
			-X 'github.com/masteryyh/micro-ddns/internal/version.GoVersion=$(GO_VERSION)' \
			-X 'github.com/masteryyh/micro-ddns/internal/version.CommitHash=$(COMMIT_HASH)'"

DOCKER_BUILD_ARGS = --build-arg USE_CN_MIRROR=$(USE_CN_MIRROR) --build-arg BUILD_TIME=$(BUILD_TIME)

DOCKER_BUILD_PLATFORMS_DEBIAN = --platform linux/amd64,linux/arm64,linux/mips64le,linux/ppc64le,linux/s390x
DOCKER_BUILD_PLATFORMS_ALPINE = --platform linux/amd64,linux/arm64,linux/ppc64le,linux/s390x,linux/riscv64

clean:
ifeq ($(DOCKER), "docker not found")
    $(error Docker is needed for build, exiting.)
endif
	@echo "Cleaning build artifacts..."
	rm -rf bin
	$(DOCKER) image rm -f $(IMG):$(TAG)-bookworm-slim
	$(DOCKER) image rm -f $(IMG):bookworm-slim
	$(DOCKER) image rm -f $(IMG):$(TAG)-alpine3.20
	$(DOCKER) image rm -f $(IMG):alpine3.20
	$(DOCKER) image rm -f $(IMG):alpine
	$(DOCKER) image rm -f $(IMG):$(TAG)
	$(DOCKER) image rm -f $(IMG):latest

build:
	@echo "Building project binary..."
	$(GO) build -ldflags $(LDFLAGS) -o $(OUTPUT_PATH) cmd/main.go

build-image:
ifeq ($(DOCKER), "docker not found")
    $(error Docker is needed for build, exiting.)
endif
	@echo "Building Docker images..."

	@echo "Building Debian images..."
	$(DOCKER) build -t $(IMG):$(TAG) $(DOCKER_BUILD_ARGS) $(DOCKER_BUILD_PLATFORMS_DEBIAN) -f build/Dockerfile-debian .
	$(DOCKER) tag $(IMG):$(TAG) $(IMG):$(TAG)-bookworm-slim
	$(DOCKER) tag $(IMG):$(TAG) $(IMG):bookworm-slim
	$(DOCKER) tag $(IMG):$(TAG) $(IMG):latest

	@echo "Building Alpine images..."
	$(DOCKER) build -t $(IMG):$(TAG)-alpine3.20 $(DOCKER_BUILD_ARGS) $(DOCKER_BUILD_PLATFORMS_ALPINE) -f build/Dockerfile-alpine .
	$(DOCKER) tag $(IMG):$(TAG)-alpine3.20 $(IMG):alpine3.20
	$(DOCKER) tag $(IMG):$(TAG)-alpine3.20 $(IMG):alpine

	@echo "Pushing images to Docker Hub..."
	$(DOCKER) push $(IMG):$(TAG)-bookworm-slim
	$(DOCKER) push $(IMG):bookworm-slim
	$(DOCKER) push $(IMG):$(TAG)-alpine3.20
	$(DOCKER) push $(IMG):alpine3.20
	$(DOCKER) push $(IMG):$(TAG)
	$(DOCKER) push $(IMG):latest

.PHONY: clean build build-image

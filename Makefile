APPNAME := $(notdir $(CURDIR))

GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
GIT_SHA := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags --exact-match 2>/dev/null || echo "$(GIT_TAG)-$(GIT_SHA)$(if $(shell git status --porcelain),-dirty)")

CMDPATH := ./cmd/$(APPNAME)
BUILDPATH := ./build

# Go parameters
GOVERSION=1.22.2
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

LDFLAGS := -ldflags '-s -w -X=github.com/mikejoh/$(APPNAME)/internal/buildinfo.Version=$(VERSION) -X=github.com/mikejoh/$(APPNAME)/internal/buildinfo.Name=$(APPNAME) -X=github.com/mikejoh/$(APPNAME)/internal/buildinfo.GitSHA=$(GIT_SHA)'

# Container image parameters
IMAGE_NAME=$(APPNAME)
IMAGE_REGISTRY=mikejoh
CHART_REPOSITORY=""

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -v -o $(BUILDPATH)/$(APPNAME) $(CMDPATH)

docker-build:
	docker build \
		-t $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(VERSION) \
		--build-arg=GOVERSION=$(GOVERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg APPNAME=$(APPNAME) \
		--build-arg GIT_SHA=$(GIT_SHA) \
		.

docker-push:
	docker push \
		$(IMAGE_REGISTRY)/$(IMAGE_NAME):$(VERSION)

docker-release: docker-build docker-push

test: 
	$(GOTEST) -v ./...

testcov:
	$(GOTEST) ./... -coverprofile=coverage.out

dep:
	$(GOCMD) mod download

vet:
	$(GOCMD) vet ./...

lint:
	golangci-lint run -v --timeout=15m ./...

clean: 
	$(GOCLEAN)
	rm -f $(BUILDPATH)/$(APPNAME)

run:
	$(GOBUILD) -v -o $(BUILDPATH)/$(APPNAME) $(CMDPATH)
	$(BUILDPATH)/$(APPNAME)

install:
	cp $(BUILDPATH)/$(APPNAME) ~/.local/bin

.PHONY: all build test testcov clean run install dep vet lint docker-build docker-push docker-release


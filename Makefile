# Image URL to use all building/pushing image targets
REGISTRY ?= quay.io
REPOSITORY ?= $(REGISTRY)/opstree
ARTIFACT_NAME=k8s-secret-injector
VERSION = 4.0

all: build-code build-image

# Build k8s-secret-injector binary
build-code:
	go build -o k8s-secret-injector

# Build the docker image
build-image:
	docker build -t $(REPOSITORY)/$(ARTIFACT_NAME):$(VERSION) -f Dockerfile .

image-push:
	docker push $(REPOSITORY)/$(ARTIFACT_NAME):$(VERSION)

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

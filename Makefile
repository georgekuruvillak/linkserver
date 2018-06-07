VERSION             := $(shell cat VERSION)
REGISTRY            := georgekk
IMAGE_REPOSITORY    := $(REGISTRY)/linkserver
IMAGE_TAG           := $(VERSION)
DOCKER_DIR           := docker
BIN_DIR             := ./bin

.PHONY: revendor
revendor:
	@dep ensure -update -v

.PHONY: build
build: 
	@./build

.PHONY: build-local
build-local:
	@env LOCAL_BUILD=1 ./build

.PHONY: docker-image
docker-image: 
	@if [[ ! -f $(BIN_DIR)/linux-amd64/linkserver ]]; then echo "No binary found. Please run 'make build'"; false; fi
	@docker build -t $(IMAGE_REPOSITORY):$(IMAGE_TAG) -f $(DOCKER_DIR)/Dockerfile --rm .

.PHONY: docker-push
docker-push:
	@if ! docker images $(IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(IMAGE_TAG); then echo "$(IMAGE_REPOSITORY) version $(IMAGE_TAG) is not yet built. Please run 'make docker-image'"; false; fi
	@docker push $(IMAGE_REPOSITORY):$(IMAGE_TAG)

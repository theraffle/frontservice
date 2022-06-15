# Current  Version
VERSION ?= v0.1.0
REGISTRY ?= changjjjjjjjj

# Image URL to use all building/pushing image targets
IMG_FRONT_SERVICE ?= $(REGISTRY)/raffle-front-service:$(VERSION)

# Build the docker image
.PHONY: docker-build
docker-build:
	docker build . -f Dockerfile -t ${IMG_FRONT_SERVICE}

# Push the docker image
.PHONY: docker-push
docker-push:
	docker push ${IMG_FRONT_SERVICE}

# Test code lint
test-lint:
	golangci-lint run ./... -v

# Generate manifests.
manifests:
	bash ./hack/release-manifest.sh $(VERSION) $(REGISTRY)

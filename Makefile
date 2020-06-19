GOOS := darwin
GOARCH := amd64
VERSION := 0.2.0
BUILD := $(shell git rev-parse --short HEAD)-dev
DATE := $(shell date "+%Y-%m-%d")

NAMESPACE := docker
REPO := oscalkit
BINARY=oscalkit_$(GOOS)_$(GOARCH)

.DEFAULT_GOAL := $(BINARY)
.PHONY: test build-docker push $(BINARY) clean generate build-docs build docs

generate:
	docker build -t $(NAMESPACE)/$(REPO):generate -f Dockerfile.generate .
	docker container run \
		-v $$PWD:/go/src/github.com/docker/oscalkit \
		$(NAMESPACE)/$(REPO):generate \
		sh -c "go generate"

test:
	@echo "Running Oscalkit test Utility"
	@echo "Running remaining tests"
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v $(shell go list ./... | grep -v "/vendor/\|/test_util/src")

build-docker:
	docker image build --build-arg VERSION=$(VERSION) --build-arg BUILD=$(BUILD) --build-arg DATE=$(DATE) -t $(NAMESPACE)/$(REPO):$(VERSION)-$(BUILD) .

push: build-docker
	docker image push $(NAMESPACE)/$(REPO):$(BUILD)

build:
	docker image build -f Dockerfile.build \
		--build-arg GOOS=$(GOOS) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--build-arg BINARY=$(BINARY) \
		-t $(NAMESPACE)/$(REPO):$(VERSION)-$(BUILD)-builder .

build-docs:
	docker image build --no-cache -t $(NAMESPACE)/$(REPO):$(VERSION)-$(BUILD)-docs docs/

docs: build-docs
	rm -rf docs/public
	git worktree prune
	rm -rf .git/worktrees/docs/public/
	git worktree add -B gh-pages docs/public upstream/gh-pages
	$(eval ID := $(shell docker create $(NAMESPACE)/$(REPO):$(VERSION)-$(BUILD)-docs))
	@docker cp $(ID):/docs/public docs/
	@docker rm $(ID) >/dev/null

# Builds binary for the OS/arch. Assumes that types have already been generated
# via the "generate" target
$(BINARY): build
	$(eval ID := $(shell docker create $(NAMESPACE)/$(REPO):$(VERSION)-$(BUILD)-builder))
	@docker cp $(ID):/$(BINARY) .
	@docker rm $(ID) >/dev/null

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

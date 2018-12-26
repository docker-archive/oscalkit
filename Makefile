# oscalkit - OSCAL conversion utility
# Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

# To the extent possible under law, the author(s) have dedicated all copyright
# and related and neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.

# You should have received a copy of the CC0 Public Domain Dedication along with this software.
# If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

GOOS := darwin
GOARCH := amd64
VERSION := 0.2.0
BUILD := $(shell git rev-parse --short HEAD)-dev
DATE := $(shell date "+%Y-%m-%d")

NAMESPACE := opencontrolorg
REPO := oscalkit
BINARY=oscalkit_$(GOOS)_$(GOARCH)

.DEFAULT_GOAL := $(BINARY)
.PHONY: test build-docker push $(BINARY) clean

generate:
	docker build -t $(NAMESPACE)/$(REPO):generate -f Dockerfile.generate .
	docker container run \
		-v $$PWD:/go/src/github.com/opencontrol/oscalkit \
		$(NAMESPACE)/$(REPO):generate \
		sh -c "go generate"

test:
	docker container run \
		-v $$PWD:/go/src/github.com/opencontrol/oscalkit \
		-w /go/src/github.com/opencontrol/oscalkit \
		circleci/golang:1.11 \
		sh -c "go test \$$(go list ./... | grep -v /vendor/)"

build-docker:
	docker image build --build-arg VERSION=$(VERSION) --build-arg BUILD=$(BUILD) --build-arg DATE=$(DATE) -t $(NAMESPACE)/$(REPO):$(VERSION)-$(BUILD) .

push: build-docker
	docker image push $(NAMESPACE)/$(REPO):$(BUILD)

# Builds binary for the OS/arch. Assumes that types have already been generated
# via the "generate" target
$(BINARY):
	docker image build -f Dockerfile.build \
		--build-arg GOOS=$(GOOS) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--build-arg BINARY=$(BINARY) \
		-t $(NAMESPACE)/$(REPO):$(VERSION)-$(BUILD)-builder .;
	$(eval ID := $(shell docker create $(NAMESPACE)/$(REPO):$(VERSION)-$(BUILD)-builder))
	@docker cp $(ID):/$(BINARY) .
	@docker rm $(ID) >/dev/null

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

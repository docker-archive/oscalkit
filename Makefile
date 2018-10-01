# oscalkit - OSCAL conversion utility
# Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

# To the extent possible under law, the author(s) have dedicated all copyright
# and related and neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.

# You should have received a copy of the CC0 Public Domain Dedication along with this software.
# If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

GOOS ?= darwin
GOARCH ?= amd64
LDFLAGS=-ldflags "-s -w"
NAMESPACE ?= opencontrolorg
REPO ?= oscalkit
BUILD ?= dev
BINARY=oscalkit_$(GOOS)_$(GOARCH)

.DEFAULT_GOAL := $(BINARY)
.PHONY: test build-docker push $(BINARY) clean

generate:
	docker build -t $(NAMESPACE)/$(REPO):generate -f Dockerfile.generate .
	docker container run \
		-v $$PWD:/go/src/github.com/opencontrol/oscalkit \
		$(NAMESPACE)/$(REPO):generate \
		sh -c "go generate"

test: generate
	docker container run \
		-v $$PWD:/go/src/github.com/opencontrol/oscalkit \
		-w /go/src/github.com/opencontrol/oscalkit \
		golang:1.11 \
		sh -c "go test \$$(go list ./... | grep -v /vendor/)"

build-docker:
	docker image build -t $(NAMESPACE)/$(REPO):$(BUILD) .

push: build-docker
	docker image push $(NAMESPACE)/$(REPO):$(BUILD)

$(BINARY): generate
	docker container run --rm \
		-v $$PWD:/go/src/github.com/opencontrol/oscalkit \
		-w /go/src/github.com/opencontrol/oscalkit/cli \
		golang:1.11-alpine \
		sh -c 'GOOS=${GOOS} GOARCH=${GOARCH} go build -v ${LDFLAGS} -o ../${BINARY}'

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

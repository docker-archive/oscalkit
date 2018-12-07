# oscalkit - OSCAL conversion utility
# Written in 2017 by Andrew Weiss <andrew.weiss@docker.com>

# To the extent possible under law, the author(s) have dedicated all copyright
# and related and neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.

# You should have received a copy of the CC0 Public Domain Dedication along with this software.
# If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

FROM golang:1.11-alpine AS builder
WORKDIR /go/src/github.com/opencontrol/oscalkit
ARG VERSION
ARG BUILD
ARG DATE
COPY . .
WORKDIR /go/src/github.com/opencontrol/oscalkit/cli
RUN CGO_ENABLED=0 go build -o oscalkit -v -ldflags "-s -w -X github.com/opencontrol/oscalkit/cli/version.Version=${VERSION} -X github.com/opencontrol/oscalkit/cli/version.Build=${BUILD} -X github.com/opencontrol/oscalkit/cli/version.Date=${DATE}"

FROM alpine:3.7
RUN apk --no-cache add ca-certificates libxml2-utils
WORKDIR /oscalkit
COPY --from=builder /go/src/github.com/opencontrol/oscalkit/cli/oscalkit /oscalkit-linux-x86_64
RUN ln -s /oscalkit-linux-x86_64 /usr/local/bin/oscalkit
ENTRYPOINT ["oscalkit"]

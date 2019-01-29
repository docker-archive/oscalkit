FROM golang:1.11-alpine AS builder
WORKDIR /go/src/github.com/docker/oscalkit
ARG VERSION
ARG BUILD
ARG DATE
COPY . .
WORKDIR /go/src/github.com/docker/oscalkit/cli
RUN CGO_ENABLED=0 go build -o oscalkit -v -ldflags "-s -w -X github.com/docker/oscalkit/cli/version.Version=${VERSION} -X github.com/docker/oscalkit/cli/version.Build=${BUILD} -X github.com/docker/oscalkit/cli/version.Date=${DATE}"

FROM alpine:3.7
RUN apk --no-cache add ca-certificates libxml2-utils
WORKDIR /oscalkit
COPY --from=builder /go/src/github.com/docker/oscalkit/cli/oscalkit /oscalkit-linux-x86_64
RUN ln -s /oscalkit-linux-x86_64 /usr/local/bin/oscalkit
ENTRYPOINT ["oscalkit"]

FROM golang:1.13.1-alpine3.10 as builder

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up path
ENV REPO_PATH    $GOPATH/src/github.com/irisnet/rainbow-sync
ENV GO111MODULE on

# RUN mkdir -p $REPO_PATH

COPY . $REPO_PATH
WORKDIR $REPO_PATH

# Install minimum necessary dependencies, build binary
RUN apk add --no-cache $PACKAGES && \
    cd $REPO_PATH && make all

FROM alpine:3.10

ENV BINARY_NAME	rainbow-sync-iris
COPY --from=builder /go/src/$BINARY_NAME /usr/local/bin/$BINARY_NAME

CMD $BINARY_NAME
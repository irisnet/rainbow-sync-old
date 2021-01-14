FROM golang:1.15.5-alpine3.12 as builder

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up path
ENV REPO_PATH	$GOPATH/src
ENV GO111MODULE on

# RUN mkdir -p $REPO_PATH

COPY . $REPO_PATH
WORKDIR $REPO_PATH

# Install minimum necessary dependencies, build binary
RUN apk add --no-cache $PACKAGES && \
    cd $REPO_PATH && make all

FROM alpine:3.12

ENV BINARY_NAME	rainbow-sync-iris
COPY --from=builder /go/src/$BINARY_NAME /usr/local/bin/$BINARY_NAME

CMD $BINARY_NAME
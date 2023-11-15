FROM golang:1.18-alpine3.15 as builder

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up path
ENV REPO_PATH	$GOPATH/src
ENV GO111MODULE on

# RUN mkdir -p $REPO_PATH

ARG GOPROXY=http://192.168.0.60:8081/repository/go-bianjie/,http://nexus.bianjie.ai/repository/golang-group,https://goproxy.cn,direct
COPY . $REPO_PATH
WORKDIR $REPO_PATH

# Install minimum necessary dependencies, build binary
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
apk add --no-cache $PACKAGES && \
    cd $REPO_PATH && make all

FROM alpine:3.12

ENV BINARY_NAME	rainbow-sync-iris
COPY --from=builder /go/src/$BINARY_NAME /usr/local/bin/$BINARY_NAME

CMD $BINARY_NAME
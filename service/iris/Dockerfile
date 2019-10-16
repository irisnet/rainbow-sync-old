FROM alpine:20190925

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up path
ENV BINARY_NAME rainbow-sync
ENV GOPATH       /root/go
ENV REPO_PATH    $GOPATH/src/github.com/irisnet/rainbow-sync/service/iris
ENV PATH         $GOPATH/bin:$PATH
ENV GO111MODULE  on

RUN mkdir -p $GOPATH/bin $REPO_PATH

COPY . $REPO_PATH
WORKDIR $REPO_PATH

VOLUME $GOPATH/src/github.com/irisnet/rainbow-sync/logs

# Install minimum necessary dependencies, build binary
RUN apk add --no-cache $PACKAGES && \
    cd $REPO_PATH && make all && \
    mv $REPO_PATH/$BINARY_NAME $GOPATH/bin && \
    rm -rf $REPO_PATH/vendor && \
    rm -rf $GOPATH/pkg/* && \
    apk del $PACKAGES

CMD $BINARY_NAME
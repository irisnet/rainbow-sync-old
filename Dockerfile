FROM alpine:3.11.5  as builder

# Set up dependencies
ENV PACKAGES go make git libc-dev bash

# Set up path
ENV BINARY_NAME rainbow-sync
ENV GOPATH       /root/go
ENV REPO_PATH    $GOPATH/src/github.com/irisnet/rainbow-sync
ENV PATH         $GOPATH/bin:$PATH
ENV GO111MODULE  on

RUN mkdir -p $REPO_PATH

COPY . $REPO_PATH
WORKDIR $REPO_PATH

# Install minimum necessary dependencies, build binary
RUN apk add --no-cache $PACKAGES && make all

FROM alpine:3.11.5
WORKDIR /root/
COPY --from=builder /root/go/src/github.com/irisnet/rainbow-sync/rainbow-sync /root/
CMD ./rainbow-sync
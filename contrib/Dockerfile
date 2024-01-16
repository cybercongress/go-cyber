# Simple usage with a mounted data directory:
# > docker build -t deepchain .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.deepchain:/root/.deepchain -v ~/.cybercli:/root/.cybercli deepchain deepchain init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.deepchain:/root/.deepchain -v ~/.cybercli:/root/.cybercli deepchain deepchain start
FROM golang:1.16-alpine3.13 AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/cybercongress/cyber

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add $PACKAGES && \
  make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/cyber /usr/bin/cyber

COPY ./contrib/single-node.sh .

EXPOSE 26657

ENTRYPOINT [ "./single-node.sh" ]
# NOTE: to run this image, docker run -d -p 26657:26657 ./single-node.sh {{chain_id}} {{genesis_account}}
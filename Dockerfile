# Simple usage with a mounted data directory:
# > docker build -t gaia .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.gaiad:/gaia/.gaiad -v ~/.gaiacli:/gaia/.gaiacli gaia gaiad init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.gaiad:/gaia/.gaiad -v ~/.gaiacli:/gaia/.gaiacli gaia gaiad start
FROM golang:1.15-alpine3.12 AS build-env

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

# Set up dependencies
#ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
#WORKDIR /go/src/github.com/cybercongress/cyber

# Add source files
#COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
#RUN apk add --no-cache $PACKAGES && \
#    make install

WORKDIR /code
COPY . .

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 39dc389cc6b556280cbeaebeda2b62cf884993137b83f90d1398ac47d09d3900

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc make install

# --------------------------------------------------------

# Final image
FROM alpine:edge

ENV CYBER /cyber

# Install ca-certificates
RUN apk add --update ca-certificates

RUN addgroup cyber && \
    adduser -S -G cyber cyber -h "$CYBER"

USER cyber

WORKDIR $CYBER

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/cyber /usr/bin/cyber

# Run gaiad by default, omit entrypoint to ease using container with gaiacli
CMD ["cyber", "--help"]

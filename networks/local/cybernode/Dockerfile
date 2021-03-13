FROM alpine:3.12

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

# Changed from Alpine to Ubuntu because the keyring PR is linking to libc
# Alpine uses muslc instead of libc

#FROM ubuntu:18.04
#
#RUN apt-get update && \
#    apt-get -y upgrade && \
#    apt-get -y install curl jq file

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 39dc389cc6b556280cbeaebeda2b62cf884993137b83f90d1398ac47d09d3900

VOLUME ["/cyber"]
WORKDIR /cyber
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh


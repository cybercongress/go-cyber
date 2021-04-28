FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.14.0/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 220b85158d1ae72008f099a7ddafe27f6374518816dd5873fd8be272c5418026

VOLUME ["/cyber"]
WORKDIR /cyber
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh


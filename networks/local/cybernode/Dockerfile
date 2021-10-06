FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.16.1/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 0e62296b9f24cf3a05f8513f99cee536c7087079855ea6ffb4f89b35eccdaa66

VOLUME ["/cyber"]
WORKDIR /cyber
EXPOSE 26656 26657 1317
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh


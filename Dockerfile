FROM golang:1.11-alpine as builder

WORKDIR $GOPATH/src/github.com/cybercongress/cyberd/cosmos/poc
COPY cosmos/poc .

RUN apk add --no-cache git

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/cyberd ./cyberd
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/cyberdcli ./cyberdcli
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/cyberdproxy ./proxy

FROM alpine:edge

RUN apk add --update ca-certificates
WORKDIR /root

COPY --from=builder /go/bin/cyberd /usr/bin/cyberd
COPY --from=builder /go/bin/cyberdcli /usr/bin/cyberdcli
COPY --from=builder /go/bin/cyberdproxy /usr/bin/cyberdproxy

COPY start_script.sh start_script.sh
RUN chmod +x start_script.sh

COPY testnet /genesis.json
COPY testnet /config.toml

EXPOSE 26656 26657 26660

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]

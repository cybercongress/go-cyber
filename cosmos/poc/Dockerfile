FROM golang:1.11-alpine as builder

WORKDIR $GOPATH/src/github.com/cybercongress/cyberd/cosmos/poc
COPY . .

RUN apk add --no-cache git

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/cyberd ./cyberd

FROM alpine:edge

RUN apk add --update ca-certificates
WORKDIR /root

COPY --from=builder /go/bin/cyberd /usr/bin/cyberd

EXPOSE 26656 26657

ENTRYPOINT ["cyberd", "start"]

# Bostrom: Superintelligence grow here

Chain-ID: bostromdev-1

Genesis: [QmSB76Ggfswc9AxwHmSAP7QCigW7fqaX9RfXs51uUreVwH](http://cloudflare-ipfs.com/ipfs/QmSB76Ggfswc9AxwHmSAP7QCigW7fqaX9RfXs51uUreVwH)

Build: ```make install```

Run: ```cyber start ```

* default mode for now CPU but GPU also supported

Bostrom (cyber) - v0.2.0-alpha1 ()
- RPC: 167.172.103.118:26657
- REST: 167.172.103.118:1317
- gRPC: 167.172.103.118:9090
- Faucet: 
- Seeds: 5f6a49a68abc391a07b76eedf253b64a2d87f2fa@167.172.99.185:26656,444b1df99e031bafcdf6d421017c187c5491d704@167.172.103.118:26656

```
cyber tx resources convert 1000000nick volt 10000 --from <name> --chain-id bostromdev-1 --gas-prices 0.001nick
cyber tx resources convert 1000000nick amper 10000 --from <name> --chain-id bostromdev-1 --gas-prices 0.001nick
cyber tx graph create QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe3 QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe8 --from <name> --chain-id bostromdev-1
```
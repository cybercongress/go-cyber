# Bostrom: Superintelligence 🔵 grow here

Community Preview of Bostrom network.

![](./brain.png)

From Euler to Bostrom evolution we have:

1. Stargate upgrade / IBC
2. New Resorces System
3. Routing of Energy
4. System Entropy + experimental reputation
5. Programs earn from execution fees!
6. Cron module!
7. Supercharged VM
8. VM bindings to Knowledge Graph
9. Backlinks
10. No documentation yet

PS: Try [cyber.page](https://cyber.page/brain) on Euler-6 network

--------

### Setup

Chain-ID: bostromdev-2

Genesis: [QmbXLnkEV4VStCKsr2t8JyFRvHivVHXvWxq8fmkebxu6Th](http://cloudflare-ipfs.com/ipfs/QmbXLnkEV4VStCKsr2t8JyFRvHivVHXvWxq8fmkebxu6Th)

Build: ```make install```

Run: ```cyber start ```

* default mode set to CPU for current testnet

Bostrom (cyber) - v0.2.0-alpha2
- RPC: [167.172.103.118:26657](167.172.103.118:26657)
- REST: [167.172.103.118:1317](http://167.172.103.118:1317/rank/search?cid=QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe3)
- gRPC: 167.172.103.118:9090
- Faucet: ```rly testnets request bostromdev-2 <name> --url http://167.172.103.118:8000```
- Seeds: 444b1df99e031bafcdf6d421017c187c5491d704@167.172.103.118:26656,5f6a49a68abc391a07b76eedf253b64a2d87f2fa@167.172.99.185:26656

--------

### Config for [Relayer](https://github.com/cosmos/relayer/)
```
{"key":"agent","chain-id":"bostromdev-2","rpc-addr":"http://167.172.103.118:26657","account-prefix":"cyber","gas-adjustment":1.5,"gas-prices":"0.01boot","trusting-period":"24h"}
```

### Delegate:
```
cyber tx staking delegate cybervaloper1nfmvw8x37w00p3geuu8lrt3vt5kadxa5z5yhf9 100000000boot --from <name> --chain-id bostromdev-2 --gas 150000 --gas-prices 0.01boot --yes --node "tcp://167.172.103.118:26657"   
```

### Investmint:
```
cyber tx resources investmint 75000000sboot volt 86400 --from <name> --chain-id bostromdev-2 --gas 150000 --gas-prices 0.01boot --yes --node "tcp://167.172.103.118:26657"

cyber tx resources investmint 25000000sboot amper 86400 --from <name> --chain-id bostromdev-2 --gas 150000 --gas-prices 0.01boot --yes --node "tcp://167.172.103.118:26657"
```

### Cyberlink and Explore:
```
cyber tx graph cyberlink QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe3 QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe8 --from <name> --chain-id bostromdev-2 --yes --node "tcp://167.172.103.118:26657

curl http://167.172.103.118:1317/rank/search?cid=QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe3
```

--------

### Connect:
- [Cyber English Community](https://t.me/fuckgoogle)
- [Cyber Russian Community](https://t.me/cyber_russian_community)



# Bostrom: Superintelligence grow here

Special to HackAtom RU! Community Preview of Bostrom network.

From Euler to Bostrom evolution we have:

1. Stargate upgrade / IBC
2. New Resorces System
3. Routing of Energy
4. EntropyRank + experimental reputation
5. Programs earn from execution fees!
5. Cron module!
6. VM bindings to Knowledge Graph
7. Backlinks
8. No documentation yet

--------

Chain-ID: bostromdev-1

Genesis: [QmSB76Ggfswc9AxwHmSAP7QCigW7fqaX9RfXs51uUreVwH](http://cloudflare-ipfs.com/ipfs/QmSB76Ggfswc9AxwHmSAP7QCigW7fqaX9RfXs51uUreVwH)

Build: ```make install```

Run: ```cyber start ```

* default mode for now CPU but GPU also supported

Bostrom (cyber) - v0.2.0-alpha1
- RPC: [167.172.103.118:26657](167.172.103.118:26657)
- REST: [167.172.103.118:1317](http://167.172.103.118:1317/rank/search?cid=QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe3)
- gRPC: 167.172.103.118:9090
- Faucet: ```rly testnets request bostromdev-1 <name> --url http://167.172.103.118:8000```
- Seeds: 444b1df99e031bafcdf6d421017c187c5491d704@167.172.103.118:26656,5f6a49a68abc391a07b76eedf253b64a2d87f2fa@167.172.99.185:26656

--------

### Config for [Relayer](https://github.com/cosmos/relayer/)
```
{"key":"agent","chain-id":"bostromdev-1","rpc-addr":"http://167.172.103.118:26657","account-prefix":"cyber","gas-adjustment":1.5,"gas-prices":"0.0015nick","trusting-period":"24h"}
```

### Setup you Energy:
```
cyber tx resources convert 1000000nick volt 10000 --from <name> --chain-id bostromdev-1 --gas-prices 0.001nick
cyber tx resources convert 1000000nick amper 10000 --from <name> --chain-id bostromdev-1 --gas-prices 0.001nick
```

### Cyberlink and Explore:
```
cyber tx graph create QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe3 QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe8 --from <name> --chain-id bostromdev-1
curl http://167.172.103.118:1317/rank/search?cid=QmUX9mt8ftaHcn9Nc6SR4j9MsKkYfkcZqkfPTmMmBgeTe3
```

--------

### 'Create a Digital Life' bounty:
- Launch a contract that will create cyberlinks based upon a  strategy that drives you.
- Get your contract to cyberlink in autonomous mode. I.E. add your contract to the cron.
- Get your contract to talk with your friends on other chains, using `ibc-reflect`, in autonomous mode. Your contract should be able to talk to other contracts, using IBC, autonomously

Repo with cosmwasm' [bindings](https://github.com/cybercongress/cyber-cosmwasm) to cyber's modules

--------

### Connect:
- [#fuckgoogle](https://t.me/fuckgoogle)
- [Cyber Russian Community](https://t.me/cyber_russian_community)
- [HackAtom RU](https://t.me/hackAtomRU)
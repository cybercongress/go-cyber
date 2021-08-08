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

PS: Try [cyber.page](https://rebyc.cyber.page)

--------

### Setup

Chain-ID: bostrom-testnet-4

Genesis: [QmTywXJjKQew1wcfLc6X4i94Gp76Gjg8qbvtaEQcvQVcxV](http://cloudflare-ipfs.com/ipfs/QmTywXJjKQew1wcfLc6X4i94Gp76Gjg8qbvtaEQcvQVcxV)

Build: ```make install```

Run: ```cyber start ```

- Version - v0.2.0-beta6
- RPC: [167.172.103.118:26657](167.172.103.118:26657)
- REST: [167.172.103.118:1317](http://167.172.103.118:1317/rank/search?cid=QmdVWtX17m7UvF8FcvNLTJxcpxv2fSJd7Z3VBoYxxW9Qpu)
- gRPC: 167.172.103.118:9090
- Faucet: ```
  curl --header "Content-Type: application/json" --request POST --data '{"denom":"boot","address":"bostrom1sm9sq4wnn62tk5yz0x3fvvx2ea9efguqwvdu64"}' http://titan.cybernode.ai:8000/credit```
- Seed: d0518ce9881a4b0c5872e5e9b7c4ea8d760dad3f@85.10.207.173:26656
- Peers: 5d542c0eb40ae48dc2cac0c140aedb605ded77dc@195.201.105.229:26656,444b1df99e031bafcdf6d421017c187c5491d704@167.172.103.118:26656,5f6a49a68abc391a07b76eedf253b64a2d87f2fa@167.172.99.185:26656

For better network stability please update your `.cyber/config/config.toml` lines as following: 

```
addr_book_strict = false

persistent_peers_max_dial_period = "300s"

allow_duplicate_ip = true
```

--------

### How to migrate from bostrom-testnet-2 to bostrom-testnet-3

If you have your bostrom-testnet-2 node running on our docker container do:

```bash
docker stop bostrom-testnet-2
docker rm bostrom-testnet-2
docker rmi cyberd/cyber:bostrom-testnet-3
docker run -d --gpus all --name=bostrom-testnet-3 --restart always -p 26656:26656 -p 26657:26657 -p 1317:1317 -e ALLOW_SEARCH=true -v $HOME/.cyber:/root/.cyber  cyberd/cyber:bostrom-testnet-3
```
This will pull new image and replace genesis and cyber binary to correct versions.

If you have your node on somehow custom setup, you need to:

1. Replace your `genesis.json` to [new one](http://cloudflare-ipfs.com/ipfs/QmfCp3M7cAu6PfANSpffM8mo45begGvRxrkgkkEBvtgyq9)

2. Buid new cyber binary from release v0.2.0-beta2 and replace old one with it

3. Start the node

--------

### Config for [Relayer](https://github.com/cosmos/relayer/)
```
{"key":"agent","chain-id":"bostrom-testnet-4","rpc-addr":"http://167.172.103.118:26657","account-prefix":"bostrom","gas-adjustment":1.5,"gas-prices":"0.01boot","trusting-period":"24h"}
```

### Delegate:
```
cyber tx staking delegate bostromvaloper1hmkqhy8ygl6tnl5g8tc503rwrmmrkjcqf92r73 100000000boot --from <name> --chain-id bostrom-testnet-4 --gas 150000 --gas-prices 0.01boot --yes --node "tcp://167.172.103.118:26657"   
```

### Investmint:
```
cyber tx resources investmint 75000000hydrogen mvolt 86400 --from <name> --chain-id bostrom-testnet-4 --gas 160000 --gas-prices 0.01boot --yes --node "tcp://167.172.103.118:26657"

cyber tx resources investmint 25000000hydrogen mamper 86400 --from <name> --chain-id bostrom-testnet-4 --gas 160000 --gas-prices 0.01boot --yes --node "tcp://167.172.103.118:26657"
```

### Cyberlink and Explore:
```
cyber tx graph cyberlink QmdVWtX17m7UvF8FcvNLTJxcpxv2fSJd7Z3VBoYxxW9Qpu Qmb9xPYYwHt1F3bQysKCZzXRzAT8QLvAyMe5DyPy4rene8 --from <name> --chain-id bostrom-testnet-4 --yes --node "tcp://167.172.103.118:26657

curl http://167.172.103.118:1317/rank/search?cid=QmdVWtX17m7UvF8FcvNLTJxcpxv2fSJd7Z3VBoYxxW9Qpu
```

--------

### Connect:
- [Cyber English Community](https://t.me/fuckgoogle)
- [Cyber Russian Community](https://t.me/cyber_russian_community)



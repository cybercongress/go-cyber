# Launch network

## Export state from previous network

To export network at given state you should have fullnode stopped at required height.
You can use any existing fullnode\validator node, or sync new node from the begging.
 
Note: use `fail-before-height` flag to stop node before given height. Example:
```bash
docker run -d --runtime=nvidia -v /cyberd/daemon:/root/.cyberd  cyberd/cyberd:euler-x cyberd start --fail-before-height=322
```

Assuming you node files are located under `/cyberd/daemon` path and current chain is `euler-x`, run export command
```bash
docker run --rm --runtime=nvidia -v /cyberd/daemon:/root/.cyberd cyberd/cyberd:euler-x cyberd export
```

Now, you will have two genesis files under path `/cyberd/daemon/export`.
 
## Generate new validators gentx

Copy two genesis files into daemon config folder, for example `/cyberd/daemon-y/config`.
Copy validator key to the same directory. To add initial valdator to the genesis.json run:

```bash
docker run --rm --runtime=nvidia -v /cyberd/daemon-y:/root/.cyberd \
  -v /cyberd/cli:/root/.cyberdcli cyberd gentx --amount=10000000cyb  --name=wallet_key  --moniker=hlb
```

## Update Dockerfile

Upload new genesis.json and links files to IPFS. Update Dockerfile ipfs hashes.
Build new image locally, upload it to Dockerhub registry.

## Launch seed node.

You should backup seed node node_key.json from previous network. Using new image, launch seed node with given key. 

## Launch first validator node.

Using validators key from gentx step and new docker image launch first validator node. 

## Knowing issues

During cosmos sdk update a set of new params can be added to genesis.json, thus make current invalid. 
In such case, you need create new empty json via `testnet` command and 
use it as template for manually assembled genesis file.

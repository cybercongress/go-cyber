# Earth cheat sheets

```
# Copy from earth
scp -P 33224 earth@earth.cybernode.ai:/path/file /host/path/file
# Copy to earth
scp -P 33324 testnet/genesis.json earth@earth.cybernode.ai:/cyberdata/cyberd/config/genesis.json
scp -P 33324 testnet/config.toml earth@earth.cybernode.ai:/cyberdata/cyberd/config/config.toml
```

```
docker run -d --restart always --name=cyberd --runtime=nvidia \
 -p 34656:26656 -p 34657:26657 -p 34660:26660 \
 -v /cyberdata/cyberd:/root/.cyberd \
 cyberd/cyberd:euler-dev0
```
 
```
docker run --rm \
 -v /cyberdata/cyberd:/root/.cyberd \
 cyberd/cyberd:euler-dev0 cyberd unsafe-reset-all
```

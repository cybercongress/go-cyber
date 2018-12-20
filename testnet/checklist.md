# Before run new testnet check:

1. loose coins total amount  = initial supply
2. adjust genesis supply in x/mint
3. reset old state
4. change genesis file chain-id
5. upadte genesis gen tx

## Useful commands

```bash
# Generate genesis validator tx
docker run --rm  -v /cyberdata/cyberd:/root/.cyberd  cyberd/cyberd:euler-dev1 cyberd tendermint show-validator
./cyberd gentx --amount=100000CBD \
 --pubkey=cbdvalconspub1zcjduepqsvzxlunur5cl644ypm4tv8lt22aaedeh6uma2ev7ux7y7tdlnhnqd5f0q3 \
 --name=euler-dev_earth --chain-id=euler-dev2
```

```bash
# Copy to earth
scp -P 33324 /home/hlb/projects/cyberd/testnet/genesis.json   earth@93.125.26.210:/cyberdata/cyberd/config/
scp -P 33324 /home/hlb/projects/cyberd/testnet/config.toml   earth@93.125.26.210:/cyberdata/cyberd/config/
# Copy from earth
scp -P 33324 earth@93.125.26.210:/path/file /host/path/file
```

```bash
# Reset node
docker run --rm  -v /cyberdata/cyberd:/root/.cyberd  cyberd/cyberd:euler-dev2 cyberd unsafe-reset-all
```

```bash
# Run node
docker run -d --restart always --name=cyberd --runtime=nvidia \
 -p 34656:26656 -p 34657:26657 -p 34660:26660 -v /cyberdata/cyberd:/root/.cyberd  cyberd/cyberd:euler-dev2
```
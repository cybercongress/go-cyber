# Running two validators testnet

## First machine - First validator

Build cyberd daemon and cli with default instructions. Let's suppose binaries named `daemon` and `cli`.

Create new account and copy seed phrase.
```
./cli keys add testnet_v2
```

Create network files for one validator:
```
./daemon testnet --v 1
```

Add generated previously account to genesis.
```bash
./daemon add-genesis-account cbd1pnuzd4tvqvrawjtk6q8wwvg9cmw4mnnmy9skdy 500CBD --home=./mytestnet/node0/cyberd
cp testnet/config.toml  ./mytestnet/node0/cyberd/config/config.toml
```

Run daemon
```bash
./daemon start --home=./mytestnet/node0/cyberd
```

## Second machine - Second validator

Build cyberd daemon and cli with default instructions. Let's suppose binaries named `daemon` and `cli`.


Init daemon.
```asciidoc
./daemon init --home=./mytestnet/
```

Copy genesis and config
```bash
scp -P 33324 earth@earth.cybernode.ai:/home/earth/cyberd/mytestnet/node0/cyberd/config/genesis.json ./mytestnet/cyberd/config
scp -P 33324 earth@earth.cybernode.ai:/home/earth/cyberd/mytestnet/node0/cyberd/config/config.toml ./mytestnet/cyberd/config
```


```bash
./daemon tendermint show-validator --home=./mytestnet/

./cli tx stake create-validator \
  --amount=100CBD \
  --pubkey=cbdvalconspub1zcjduepq2pr9cqd9apves6av0xafzgzvvl0ky8904wlkecwu4vghxteutrsq8pcfa7 \
  --moniker="hlb" \
  --chain-id=chain-hFpwd3 \
  --from=testnet_hlb \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01"
```




# Cyberd Usage Guide

## Installing Cyberd Node

### Use docker

In order to start cyberd node locally using docker run following command (replace ${YOUR_DATA_LOCAL_FOLDER} and ${YOUR_CONFIG_LOCAL_FOLDER} with your local folders where you want to store cyberd data and configuration):
```bash
docker run -d --restart always --name=cyberd -p 26656:26656 -p 26657:26657 -v ${YOUR_DATA_LOCAL_FOLDER}:/root/.cyberd/data -v ${YOUR_CONFIG_LOCAL_FOLDER}:/root/.cyberd/config cybernode/cyberd:master
```
It will run NON-VALIDATOR local node and connect it to our seed node called "earth".

You could check that node is running by executing:
```bash
docker logs cyberd
```

### Use compiled binary

You could find latest binaries in our [releases](https://github.com/cybercongress/cyberd/releases).
Choose appropriate binary for your system, download them and add to $PATH.

To connect to our network you should copy `genesis.json` and `config.toml` files from [here](https://github.com/cybercongress/cyberd/tree/master/cosmos/poc)
and put them into `$HOME/.cyberd/config/` folder

After everything is set run:
```bash
cyberd init
```
to initialize your node.

To start node use following command:
```bash
cyberd start
```

Node is started! You should see logs with generated blocks.

### Build binaries manually

You need to have GO 1.11+ installed on your computer.

Install GO by following the [official docs](https://golang.org/doc/install). Remember to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables, for example:

```$bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=$PATH:$GOBIN" >> ~/.bash_profile
```

Once you have GO installed, run the command:
```$bash
go get github.com/cybercongress/cyberd
```

There will be an error stating `can't load package: package github.com/cosmos/cosmos-sdk: no Go files`,
however you can ignore this error, it doesn't affect us.

Now change directories to:
```$bash
cd $GOPATH/src/github.com/cybercongress/cyberd/cosmos/poc
```

And run:
```$bash
export GO111MODULE=on
go install ./cyberd
go install ./cyberdcli
```

This will create binaries for `cyberd` and `cyberdcli`

After you built binaries follow [this](#use-compiled-binary) instructions to run node.

## Using cyberdcli

You could find latest cli in our [releases section](https://github.com/cybercongress/cyberd/releases)
or build it by yourself using [this](#build-binaries-manually) guide.

### Generating keys

After you have cyberdcli binary installed you probably want to generate personal keys to sign and broadcast transactions.

To generate keys simply run:
```
cyberdcli keys add ${your-key-name}
```

Enter and confirm a passphrase:

```
Enter a passphrase for your key:
Repeat the passphrase:
```

You just created your first locally stored key, under the given name. 
Once you do this, the  `~/.cyberdcli` folder is created, which will hold the this key and any other keys you make.
Now that you have the key for alice, you can start broadcasting transactions.

### Link transaction

The following command will link cid1 with cid2.
```
cyberdcli link --from=${your_key_name} --cid1=42 --cid2=cyberd --sequence=0 --chain-id=test-chain-fbqPMq
```
You could find chain id in `$HOME/.cyberd/config/genesis.json`. For our zeronet chain id is `test-chain-fbqPMq`.
Every new transaction should have incremented `--sequence` parameter.

If everything went fine you should see similar message:
```
Committed at block 107 (tx hash: BE458B956646F8B3F25F071A958A5FD7E908791F)
```

### Search transaction by hash

To find transaction run:
```bash
cyberdcli tx BE458B956646F8B3F25F071A958A5FD7E908791F
```

### Help

```bash
cyberdcli --help
```
Please note that currently not all functions are available.

## RPC client

RPC client is available on `localhost:26657`. You could find URLs list [here](https://tendermint.github.io/slate/)

### Our public endpoint

```http://earth.cybernode.ai:34657```

Example
```
http://earth.cybernode.ai:34657/block?height=42
```

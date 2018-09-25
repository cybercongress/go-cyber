# Cyberd Usage Guide

You will need to have GO 1.11+ installed on your computer.

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

## Using cyberd and cyberdcli

Let's start by initializing the cyberd daemon. Run the command
```bash
cyberd init
```
And you should see something like this:
```json
{
  "chain_id": "test-chain-rBoHlB",
  "node_id": "1afb0a84c9b33eb87a6ee9d7da8b0cf14245c1f6",
  "app_message": {
    "secret": "enough gate sock devote move lumber weekend gesture illness lucky that story memory ocean mad horse allow easily tunnel room sand unlock honey onion"
  }
}
```

This creates the `~/.cyberd` folder, which has config.toml, genesis.json, node_key.json, priv_validator.json. Take some time to review what is contained in these files if you want to understand what is going on at a deeper level.

### Generating keys

The next thing we'll need to do is add the key from priv_validator.json to the gaiacli key manager. For this we need the 16 word seed that represents the private key, and a password. You can also get the 16 word seed from the output seen above, under `"secret"`. Then run the command:

```
cyberdcli keys add alice --recover
```

Which will give you three prompts:

```
Enter a passphrase for your key:
Repeat the passphrase:
Enter your recovery seed phrase:
```

You just created your first locally stored key, under the name alice, and this account is linked to the private key that is running the cyberd validator node. Once you do this, the  ~/.cyberdcli folder is created, which will hold the alice key and any other keys you make. Now that you have the key for alice, you can start up the blockchain by running

```
cyberd start
```

You should see blocks being created at a fast rate, with a lot of output in the terminal.

### Link transaction

The following command will link cid1 with cid2

```
cyberdcli link --from=alice --cid1=42 --cid2=cyberd --sequence=0 --chain-id=test-chain-rBoHlB
```

### See links

The following command will show links of cid

```
cyberdcli links <cid-here>
```
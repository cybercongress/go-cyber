# Ultimate cyberd CLI guide. Chain: euler-6

## Install client separately

It is possible to interract with cyber even if you don't have your own node. For that case you just need to install `cyberdcli` to your machine using script below (just paste it in the console):

```bash
bash < <(curl -s https://mars.cybernode.ai/go-cyber/install_cyberdcli_v0.1.6.sh)
```

To start using client without own node you have to configure some parameters as: [key storage](#keyring-manipulation-settings), [connetion to remote node](#usefull-client-configuration), etc.

After installation you would be able to use `cyberdcli` to [import account](#account-management), create [links](#linking-content) or, for example, running validator.

First of all I would like to encourage you to use  `--help` feature if you want to get better experience of using cyberdcli. This is really easy way to find all necessary commands with options and flags.

For example you can enter:

```bash
cyberdcli --help
```

You should see this message:

```bash
Command line interface for interacting with cyberd

Usage:
cyberdcli [command]

Available Commands:
status      Query remote node for status
config      Create or query an application CLI configuration file
query       Querying subcommands
tx          Transactions subcommands

rest-server Start LCD (light-client daemon), a local REST server

keys        Add or view local private keys

version     Print the app version
link        Create and sign a link tx
help        Help about any command

Flags:
    --chain-id string   Chain ID of tendermint node
-e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
-h, --help              help for cyberdcli
    --home string       directory for config and data (default "/root/.cyberdcli")
-o, --output string     Output format (text|json) (default "text")
    --trace             print out full stack trace on errors
```

Help feature working as a stairs - you can use it with any command to find available options, subcommands and flags. For example lets explore `query` subcommands:

```bash
cyberdcli query --help
```

now, you can see subcommand structure:

```bash
Usage:
cyberdcli query [command]
```

and available subcommands and flags:

```bash
Available Commands:
account                  Query account balance

tendermint-validator-set Get the full tendermint validator set at given height
block                    Get verified data for a the block at given height
txs                      Query for paginated transactions that match a set of tags
tx                       Query for a transaction by hash in a committed block

staking                  Querying commands for the staking module
slashing                 Querying commands for the slashing module
supply                   Querying commands for the supply module
bandwidth                Querying commands for the bandwidth module
auth                     Querying commands for the auth module
mint                     Querying commands for the minting module
distribution             Querying commands for the distribution module
gov                      Querying commands for the governance module
rank                     Querying commands for the rank module

Flags:
-h, --help   help for query

Global Flags:
    --chain-id string   Chain ID of tendermint node
-e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
    --home string       directory for config and data (default "/root/.cyberdcli")
-o, --output string     Output format (text|json) (default "text")
    --trace             print out full stack trace on errors
```

Alright, lets explore `account` subcommand:

```bash
cyberdcli query account --help
```

Now we see all options available at this subcommands, namely, account address and flags:

```bash
Usage:
cyberdcli query account [address] [flags]
```

In most cases you need just two extra flags:

```bash
--from <your_key_name> \
--chain-id euler-6
```

That it. This is very useful ability for using cyberdcli and troubleshooting.

## Glossary

**Bandwidth** - The recovered unit of your account. Used to complete transactions in the cyberd blockchain. The amount of your bandwidth calculates like:

`your_eul_tokens / all_eul_tokens_in_cyberd * 2000*1000*100`.

Messages cost is `500` (exclude link). Transaction consists of one or more messages `m_1, m_2, ..., m_n`. Transaction cost is `300 + c_1 + c_2 ... + c_n`, where `c_i` - cost of `m_i` message. Full bandwidth regeneration time is 86400 blocks (24 hours)

**commission** -  tokens that you've earned with validation. You can take them at any time.

**illiquid tokens** - non-transferable tokens that you've delegated to the validator. Delegation process duration - 1 block. **Unbonding** process, or taking back share -  5 days (for `euler-6` only).

**link** - reference between CID key and CID value. Link message cost is `100*n`, where `n` is quantity of links in message. Link finalization time is 1 block. New rank for CIDs of link will be recalculated at period from 100 to 200 blocks (from 100 to 200 seconds).

**liquid tokens** - transferable tokens in cyberd blockchain

**local keystore** - store with keys in you local machine

**rewards** - tokens that you've earned with the delegation. To reduce network loads all rewards storing in a pool. You can take your part of bounty at any time by commands at **delegator** section.

**<comission_rate_percentage>** - the commission that validator get for the work. Must be fraction >0 and <=1

**<delegator_address>** - delegator address. Starts with `cyber` most often coinciding with **<key_address>**

**<key_address>** - account address. Starts with `cyber`

**<key_name>** - name of account in cybercli

**<operator_address>** - validator address. Starts with `cybervaloper`

**<shares_percentage>** - the part of illiquid tokens that you want to unbonding or redelegate. Must be fraction >0 and <=1

**<testnet_chain_id>** - version of testnet(current is euler-6).

## General commands

### Show all validators

Return set of all active and jailed validators.

```bash
cyberdcli query staking validators --trust-node
```

### Show chain status

Return general chain information

```bash
 cyberdcli status --indent
```

### Distribution params

```bash
 cyberdcli query distribution params --trust-node
```

### The amount of outstanding rewards for validator

Return the sum of outstanding rewards for validator

```bash
 cyberdcli query distribution validator-outstanding-rewards <operator_address> --trust-node
```

### Staking params

Chain staking info

```bash
 cyberdcli query staking params --trust-node
```

### Staking pool

```bash
 cyberdcli query staking pool --trust-node
```

## Account management

Don't have an account? Check out if you have some [gift](https://cyber.page/gift) allocated to you!

### Import an account by seed phrase and store it in local keystore

```bash
 cyberdcli keys add <your_key_name> --recover
```

### Import an account by private key and store it in local keystore (private key could be your ETH private key)

```bash
 cyberdcli keys add private <your_key_name>
```

### Create a new account

```bash
 cyberdcli keys add <your_key_name>
```

### Show account information

Name, address and public key of current account

```bash
 cyberdcli keys show <your_key_name>
```

### Show account balance

Return account number, balance, public key in 16 and sequence.

```bash
 cyberdcli query account <your_key_address>
```

### List existing keys

Return all keys in cyberdcli

```bash
 cyberdcli keys list
```

### Delete account from cybercli

```bash
 cyberdcli keys delete <deleting_key_name>
```

### Keyring manipulation settings

**Important note**: Since v.38 cosmos-sdk uses os-native keyring to store all your keys. We've noticed that in several cases it does not work well by default (for example if you dont have GUI installed on you machine), so if during execituon `cyberdcli keys add` command you've got this kind of error:

```bash
panic: No such  interface 'org.freedesktop.DBus.Properties' on object at path /

goroutine 1 [running]:
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.writeInfo(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x147a6c0, 0xc000f1c780, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:479 +0x38c
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.writeLocalKey(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x147a6c0, 0xc000f1c780, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:465 +0x189
github.com/cosmos/cosmos-sdk/crypto/keys.baseKeybase.CreateAccount(0x1307a18, 0x1307a10, 0xc000b37160, 0x1, 0x1, 0xc000b37170, 0x1, 0x1, 0x146aa00, 0xc000b15630, ...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keybase_base.go:171 +0x192
github.com/cosmos/cosmos-sdk/crypto/keys.keyringKeybase.CreateAccount(...)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/crypto/keys/keyring.go:107
github.com/cosmos/cosmos-sdk/client/keys.RunAddCmd(0xc000f0b400, 0xc000f125f0, 0x1, 0x1, 0x148dcc0, 0xc000aca550, 0xc000ea75c0, 0xc000ae1c08, 0x5e93b7)
    /root/go/pkg/mod/github.com/cosmos/cosmos-sdk@v0.38.1/client/keys/add.go:273 +0xa8b
... etc
```

You will have to use another keyring backend to keep your keys. Here are 2 options: store the files within the cli folder or a `pass` manager.

Setting keyring backend to **local file**:

Execute:

```bash
cyberdcli config keyring-backend file
```

As a result you migth see following: `configuration saved to /root/.cybercli/config/config.toml`

Execute:

```bash
cyberdcli config --get keyring-backend
```

The result must be the following:

```bash
user@node:~# cyberdcli config --get keyring-backend
file
```

That means that you've set your keyring-backend to a local file. *Note*, in this case, all the keys in your keyring will be encrypted using the same password. If you would like to set up a unique password for each key, you should set a unique `--home` folder for each key. To do that, just use `--home=/<unique_path_to_key_folder>/` with setup keyring backend and at all interactions with keys when using cyberdcli:

```bash
cyberdcli config keyring-backend file --home=/<unique_path_to_key_folder>/
cyberdcli keys add <your_second_key_name> --home=/<unique_path_to_key_folder>/
cyberdcli keys list --home=/<unique_path_to_key_folder>/
```

Set keyring backend to [**pass manager**](https://github.com/cosmos/cosmos-sdk/blob/9cce836c08d14dc6836d07164dd964b2b7226f36/crypto/keyring/doc.go#L30):

Pass utility uses a GPG key to encrypt your keys (but again, it uses the same GPG for all the keys). To install and generate your GPG key you should follow [this guide](https://www.passwordstore.org/) or this very [detailed guide](http://tuxlabs.com/?p=450). When you'll get your `pass` set, configure `cyberdcli` to use it as a keyring backend:

```bash
cyberdcli config keyring-backend pass
```

And verify that all set as planned:

```bash
cyberdcli config --get keyring-backend
pass
```

### Usefull client configuration

There's some thing you would like to configure, that would simplify you interaction with cli. To see all possible parameters of those hints just run:

```bash
cyberdcli config --help
```

And here's what you can configure:

```bash
Create or query an application CLI configuration file

Usage:
  cybercli config <key> [value] [flags]

Flags:
      --get    print configuration value or its default if unset
  -h, --help   help for config

Global Flags:
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/Users/mr_laptanovi4/.cybercli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

First of all, if you using standalone client, that is good idea to set up addres of the node to process all of your transactions. To to that use:

```bash
cyberdcli config node <http://node_address:port>
```

We provide public api address: `http://titan.cybernode.ai:26657`.

> TODO add new addres for public api

Also it will be usefull to configure `chain-id` to avoid entering it aevery time:

```bash
cyberd config chain-id euler-6
```

And in case if you have some troubles with [key storage](#keyring-manipulation-settings), you may want to save home directory for your cli:

```bash
cyberdcli config --home /path_to_cli_home/.cybercli/
```

To check what is the current setup for any parameter run:

```bash
cyberdcli config --get <parameter_name>
```

### Send tokens

```bash
 cyberdcli tx send <from_key_or_address> <to_address> <amount_eul> --chain-id euler-6
```

### Linking content

>Just ipfs hashes available as a CID

```bash
 cyberdcli link \
  --from=<your_key_name> \
  --cid-from=<key_phrase_to_link> \
  --cid-to=<content_that_you_want_to_link> \
  --chain-id=euler-6
```

Example of link command:

```bash 
cyberdcli link --cid-from QmWDYzTXarWYy9UKC7Ro4xMCdSVseQPbmdnmTYsJ9zGTpK --cid-to QmVgxX3TVntSNRiQ1Kd8sE8zvEKkbEgb8PaMnA4N7w7pK3 --from fuckgoogle --chain-id euler-6 --yes
```

## Validator commands

### Get all validators

```bash
 cyberdcli query staking validators --trust-node
```

### The amount of commission

Available to withdraw validator commission.

```bash
 cyberdcli query distribution commission <operator_address>
```

### State of current validator

```bash
docker exec eiler-5 cyberdcli query staking validator <operator_address>
```

### Return all delegations to validator

```bash
 cyberdcli query staking delegations-to <operator_address>
```

### Edit commission in existing validator account

```bash
 cyberdcli tx staking edit-validator \
  --from=<your_key_name> \
  --commission-rate=<new_comission_rate_percentage> \
  --chain-id=euler-6
```

### Withdraw commission for either a delegation

```bash
 cyberdcli tx distribution withdraw-rewards <operator_address> \
  --from=<your_key_name> \
  --chain-id=euler-6 \
  --commission
```

### Edit site and description in existing validator account

```bash
 cyberdcli tx staking edit-validator \
  --from=<your_key_name> \
  --details="<description>" \
  --website=<your_website> \
  --chain-id=euler-6
```

### Unjail validator previously jailed for downtime

```bash
 cyberdcli tx slashing unjail --from=<your_key_name> --chain-id=euler-6
```

### Get info about redelegation process from validator

```bash
 cyberdcli query staking redelegations-from <operator_address>
```

## Delegator commands

### Return distribution delegator rewards according current validator

```bash
 cyberdcli query distribution rewards <delegator_address> <operator_address>
```

### Return delegator shares with current validator

```bash
 cyberdcli query staking delegation <delegator_address> <operator_address>
```

### Return all delegations made from one delegator

```bash
 cyberdcli query staking delegations <delegator_address>
```

### Return all unbonding delegatations from a validator

```bash
 cyberdcli query staking unbonding-delegations-from <operator_address>
```

### Withdraw rewards for either a delegation

```bash
 cyberdcli tx distribution withdraw-rewards <operator_address> \
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Withdraw all delegation rewards

```bash
 cyberdcli tx distribution withdraw-all-rewards \
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Change the default withdraw address for rewards associated with an address

```bash
 cyberdcli tx distribution set-withdraw-addr <your_new_address> \
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Delegate liquid tokens to a validator

```bash
 cyberdcli tx staking delegate <operator_address> <amount_cyb> \
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Redelegate illiquid tokens from one validator to another in absolute cyb value

> 5 days for unbonding first.

```bash
 cyberdcli tx staking redelegate <old_operator_address> <new_operator_address> <amount_cyb> --from=<your_key_name> --chain-id=euler-6
```

### Redelegate illiquid tokens from one validator to another in percentages

```bash
 cyberdcli tx staking redelegate <old_operator_address> <new_operator_address> <shares_percentage>
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Unbond shares from a validator in absolute cyb value

> 5 days for unbonding.

```bash
 cyberdcli tx staking unbond <operator_address> <amount_cyb>
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Unbond shares from a validator in percentages

>5 days unbonding.

```bash
 cyberdcli tx staking unbond <operator_address> <shares_percentage>
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Get info about unbonding delegation process to current validator

```bash
 cyberdcli query staking unbonding-delegation <delegator_address> <operator_address>
```

### Get info about unbonding delegation process to all unbonded validators

```bash
 cyberdcli query staking unbonding-delegation <delegator_address>
```

### Get info about redelegation process from to current validator

```bash
 cyberdcli query staking redelegation <delegator_address> <old_operator_address> <new_operator_address>
```

### Get info about all redelegation processes by one delegator

```bash
 cyberdcli query staking redelegations <delegator_address>
```

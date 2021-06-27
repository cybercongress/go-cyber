# Ultimate cyberd CLI guide. Chain: bostrom-testnet-1

** TODO upgrade for bostrom network**

## Install cyberd client

It is possible to interact with cyber even if you don't have your own node. All you need to do is install `cyberdcli` on your machine using the script below (just paste it in the console):

```bash
bash < <(curl -s https://raw.githubusercontent.com/cybercongress/go-cyber/master/scripts/install_cyberdcli_v0.1.6.sh)
```

To start using the client without an own node, you have to configure some parameter Such as: [key storage](#keyring-manipulation-settings), [connetion to a remote node](#usefull-client-configuration), etc.

After installation you will be able to use `cyberdcli` to [import accounts](#account-management), create [links](#linking-content) or, for example, run a validator.

First of all, I would like to encourage you to use the  `--help` feature if you want to get a better experience of using cyberdcli. This is a really easy way to find all the necessary commands with the appropriate options and flags.

For example, you can enter:

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

The help feature works like a pyramid, you can use it with any command to find available options, subcommands and flags. For example, lets explore the `query` subcommands:

```bash
cyberdcli query --help
```

You can see the structure of the subcommand:

```bash
Usage:
cyberdcli query [command]
```

And the available subcommands and flags:

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

Let's explore the `account` subcommand:

```bash
cyberdcli query account --help
```

We can see all of the options available for this subcommands, namely, account address and flags:

```bash
Usage:
cyberdcli query account [address] [flags]
```

In most cases you will need just two extra flags:

```bash
--from <your_key_name> \
--chain-id euler-6
```

That's it. This is a very useful tool for using cyberdcli and troubleshooting.

## Glossary

**Bandwidth** - The recovered unit of your account. Used to complete transactions in the cyber blockchain. The amount of your bandwidth calculates like:

`your_eul_tokens / all_eul_tokens_in_cyber * 2000*1000*100`.

Messages cost is `100` (exclude link). Transaction consists of one or more messages `m_1, m_2, ..., m_n`. Transaction cost is `300 + c_1 + c_2 ... + c_n`, where `c_i` - cost of `m_i` message. Full bandwidth regeneration time is 16000 blocks ( ~24 hours ). All bandwith prices, as well as other parameters could be found in our [launch kit](https://github.com/cybercongress/launch-kit/tree/0.1.0/params).

**commission** -  The tokens that you've earned via validating from delegators. You may take them at any time

**illiquid tokens** - Non-transferable tokens that you've delegated to the validator. Delegation process duration: 1 block

**Hero** - A validator node

**Unbonding** - The process of taking back your share (delegated tokens + any rewards). 5 days (for `euler-6` only)

**Ledger** - A hardware wallet

**link** - A reference between a CID key and a CID value. Link message cost is `100*n`, where `n` is the number of links in a message. Link finalization time is 1 block. New rank for CIDs of links will be recalculated at a period of 100 to 200 blocks (from 100 to 200 seconds)

**liquid tokens** - Transferable tokens within the cyber blockchain

**local keystore** - A store with keys on your local machine

**rewards** - Tokens that you've earned via delegation. To reduce network load all the rewards are stored in a pool. You can take your part of the bounty at any time with commands from the **delegator** section.

**<comission_rate_percentage>** - The commission that a validator gets for their work. Must be a fraction >0 and <=1

**<delegator_address>** - Delegator address. Starts with `cyber` most often coinciding with **<key_address>**

**<key_address>** - An account address. Starts with `cyber`

**<key_name>** - The name of the account in cybercli

**<operator_address>** - Validator address. Starts with `cybervaloper`

**<shares_percentage>** - The part of illiquid tokens that you want to unbond or redelegate. Must be a fraction >0 and <=1

**<testnet_chain_id>** - The current version of the testnet (current, euler-6).

## General commands

### Show all validators

Return the set of all active and jailed validators:

```bash
cyberdcli query staking validators --trust-node
```

### Show chain status

Return general chain information:

```bash
 cyberdcli status --indent
```

### Distribution params

```bash
 cyberdcli query distribution params --trust-node
```

### The number of outstanding rewards for a validator

Return the sum of outstanding rewards for a validator:

```bash
 cyberdcli query distribution validator-outstanding-rewards <operator_address> --trust-node
```

### Staking params

Chain staking info:

```bash
 cyberdcli query staking params --trust-node
```

### Staking pool

```bash
 cyberdcli query staking pool --trust-node
```

## Account management

Don't have an account? Check out if you have a [gift](https://cyber.page/gift) allocated to you!

### Import an account with a seed phrase and store it in the local keystore

```bash
 cyberdcli keys add <your_key_name> --recover
```

### Import an account with a private key and store it in the local keystore (private key could be your ETH private key)

```bash
 cyberdcli keys add private <your_key_name>
```

### Create a new account

```bash
 cyberdcli keys add <your_key_name>
```

### Show account information

Name, address and the public key of the current account

```bash
 cyberdcli keys show <your_key_name>
```

### Show account balance

Return account number and amount of tokens.

```bash
 cyberdcli query account <your_key_address>
```

### List existing keys

Return all the existing keys in cyberdcli:

```bash
 cyberdcli keys list
```

### Delete account from cybercli

```bash
 cyberdcli keys delete <deleting_key_name>
```

### Keyring manipulation settings

**Important note**: Starting with v.38, Cosmos-SDK uses os-native keyring to store all of your keys. We've noticed that in certain cases it does not work well by default (for example if you don't have any GUI installed on your machine). If during the execution `cyberdcli keys add` command, you are getting this type of error:

```bash
panic: No such interface 'org.freedesktop.DBus.Properties' on object at path /

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

Setting keyring backend to a **local file**:

Execute:

```bash
cyberdcli config keyring-backend file
```

As a result you migth see the following: `configuration saved to /root/.cybercli/config/config.toml`

Execute:

```bash
cyberdcli config --get keyring-backend
```

The result should be the following:

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

Verify that all of the settings went as planned:

```bash
cyberdcli config --get keyring-backend
pass
```

### Usefull client configuration

There's some hints you may configure to simplify you interaction with client. To see all possible parameters of those hints just run:

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

First of all, if you using a standalone client, it is a good idea to set up an address of a node to process all of your transactions. To do this use:

```bash
cyberdcli config node <http://node_address:port>
```

We provide a public API address: `http://titan.cybernode.ai:26657`.

> TO DO add a new address for public API

It will be useful to configure the `chain-id` to avoid entering it every time:

```bash
cyberd config chain-id euler-6
```

If you are having trouble with the [key storage](#keyring-manipulation-settings), you may want to save your home directory for your cli:

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

> Only IPFS hashes are available to use as CIDs

```bash
 cyberdcli link \
  --from=<your_key_name> \
  --cid-from=<key_phrase_to_link> \
  --cid-to=<content_that_you_want_to_link> \
  --chain-id=euler-6
```

The `cid-from` is the IPFS hash of the keyword that you want to make searchable. Example of a link command:

```bash 
cyberdcli link --cid-from QmWDYzTXarWYy9UKC7Ro4xMCdSVseQPbmdnmTYsJ9zGTpK --cid-to QmVgxX3TVntSNRiQ1Kd8sE8zvEKkbEgb8PaMnA4N7w7pK3 --from fuckgoogle --chain-id euler-6 --yes
```

## Validator commands

### Get all validators

```bash
 cyberdcli query staking validators --trust-node
```

### The amount of commission

The commission available to withdraw for a validator:

```bash
 cyberdcli query distribution commission <operator_address>
```

### State of a current validator

```bash
docker exec eiler-5 cyberdcli query staking validator <operator_address>
```

### Return all delegations to a validator

```bash
 cyberdcli query staking delegations-to <operator_address>
```

### Edit the commission in an existing validator account

```bash
 cyberdcli tx staking edit-validator \
  --from=<your_key_name> \
  --commission-rate=<new_comission_rate_percentage> \
  --chain-id=euler-6
```

### Withdraw the commission for any delegation

```bash
 cyberdcli tx distribution withdraw-rewards <operator_address> \
  --from=<your_key_name> \
  --chain-id=euler-6 \
  --commission
```

### Edit the site and description for an existing validator account

```bash
 cyberdcli tx staking edit-validator \
  --from=<your_key_name> \
  --details="<description>" \
  --website=<your_website> \
  --chain-id=euler-6
```

### Unjail a validator previously jailed for downtime

```bash
 cyberdcli tx slashing unjail --from=<your_key_name> --chain-id=euler-6
```

### Get info about a redelegation process from a validator

```bash
 cyberdcli query staking redelegations-from <operator_address>
```

## Delegator commands

### Return distribution delegator rewards for a specified validator

```bash
 cyberdcli query distribution rewards <delegator_address> <operator_address>
```

### Return delegator shares for the specified validator

```bash
 cyberdcli query staking delegation <delegator_address> <operator_address>
```

### Return all of the delegations made from a delegator

```bash
 cyberdcli query staking delegations <delegator_address>
```

### Return all unbonding delegations from a validator

```bash
 cyberdcli query staking unbonding-delegations-from <operator_address>
```

### Withdraw rewards for any delegation

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

### Change the default withdrawal address for rewards associated with an address

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

> There is a 5-day unbonding period

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

> 5 days for unbonding

```bash
 cyberdcli tx staking unbond <operator_address> <amount_cyb>
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Unbond shares from a validator in percentages

> 5 days for unbonding

```bash
 cyberdcli tx staking unbond <operator_address> <shares_percentage>
  --from=<your_key_name> \
  --chain-id=euler-6
```

### Get info about the unbonding delegation process to any validator

```bash
 cyberdcli query staking unbonding-delegation <delegator_address> <operator_address>
```

### Get info about the unbonding delegation process to all unbonded validators

```bash
 cyberdcli query staking unbonding-delegation <delegator_address>
```

### Get info about redelegation process from to current validator

```bash
 cyberdcli query staking redelegation <delegator_address> <old_operator_address> <new_operator_address>
```

### Get the info about all the redelegation processes by a delegator

```bash
 cyberdcli query staking redelegations <delegator_address>
```

## Governance and voting

### Query specific proposal

```bash
cyberdcli q gov proposal <proposal_id> --trust-node
```

### Query all proposals

```bash
cyberdcli q gov proposals --trust-node
```

### Query votes on proposal

```bash
cyberdcli q gov votes --trust-node
```

### Query parameters from governance module

```bash
cyberdcli q gov params
```

### Vote for specific proposal

```bash
cyberdcli tx gov vote <proposal_id> <vote_option:_yes_no_abstain> --from <your_key_name> --chain-id euler-6
```

### Submit text proposal

```bash
cyberdcli tx gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10eul" --from <your_key_name> --chain-id euler-6
```

### Submit community spend proposal

```bash
cyberdcli tx gov submit-proposal community-pool-spend <path/to/proposal.json> --from <key_or_address> --chain-id euler-6
```

Where `proposal.json` is a file, structured in following way:

```json
{
  "title": "Community Pool Spend",
  "description": "Pay me some Euls!",
  "recipient": "cyber1s5afhd6gxevu37mkqcvvsj8qeylhn0rz46zdlq",
  "amount": [
    {
      "denom": "eul",
      "amount": "10000"
    }
  ],
  "deposit": [
    {
      "denom": "eul",
      "amount": "10000"
    }
  ]
}
```

### Submit chain parameters change proposal

```bash
cyberdcli tx gov submit-proposal param-change <path/to/proposal.json> --from=<key_or_address> --chain-id euler-6
```

Where `proposal.json` is a file, structured in following way:

```json
{
  "title": "Staking Param Change",
  "description": "Update max validators",
  "changes": [
    {
      "subspace": "staking",
      "key": "MaxValidators",
      "value": 105
    }
  ],
  "deposit": [
    {
      "denom": "stake",
      "amount": "10000"
    }
  ]
}
```

# Ultimate cyberd CLI guide. Chain: Euler-5

## Install client separately

It is possible to interract with cyber even if you don't have your own node. For that case you just need to install `cyberdcli` to your machine using script below:



  First of all I would like to encourage you to use  `--help` feature if you want to get better experience of using cyberdcli. This is really easy way to find all necessary commands with options and flags.

  For example you can enter:

  ```bash
  docker exec euler-5 cyberdcli --help
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
  docker exec euler-5 cyberdcli query --help
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
  docker exec euler-5 cyberdcli query account --help
  ```

  Now we see all options available at this subcommands, namely, account address and flags:

  ```bash
  Usage:
  cyberdcli query account [address] [flags]
  ```

  In most cases you need just two extra flags:

  ```bash
  --from=<your_key_name> \
  --chain-id=euler-5
  ```

  That it. This is very useful ability for using cyberdcli and troubleshooting.

## Glossary

  **Bandwidth** - The recovered unit of your account. Used to complete transactions in the cyberd blockchain. The amount of your bandwidth calculates like:

  `your_eul_tokens / all_eul_tokens_in_cyberd * 2000*1000*100`.

  Messages cost is `500` (exclude link). Transaction consists of one or more messages `m_1, m_2, ..., m_n`. Transaction cost is `300 + c_1 + c_2 ... + c_n`, where `c_i` - cost of `m_i` message. Full bandwidth regeneration time is 86400 blocks (24 hours)

  **commission** -  tokens that you've earned with validation. You can take them at any time.

  **illiquid tokens** - non-transferable tokens that you've delegated to the validator. Delegation process duration - 1 block. **Unbonding** process, or taking back share -  5 days (for `euler-5` only).

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

  **<testnet_chain_id>** - version of testnet(current is euler-5).

## General commands

### Show all validators

Return set of all active and jailed validators.

```bash
docker exec euler-5 cyberdcli query staking validators --trust-node
```

### Show chain status

Return general chain information

```bash
docker exec euler-5 cyberdcli status --indent
```

### Distribution params

```bash
docker exec euler-5 cyberdcli query distribution params --trust-node
```

### The amount of outstanding rewards for validator

Return the sum of outstanding rewards for validator

```bash
docker exec euler-5 cyberdcli query distribution validator-outstanding-rewards <operator_address> --trust-node
```

### Staking params

Chain staking info

```bash
docker exec euler-5 cyberdcli query staking params --trust-node
```

### Staking pool

```bash
docker exec euler-5 cyberdcli query staking pool --trust-node
```

## Account management

### Import an account by seed phrase and store it in local keystore

```bash
docker exec -ti euler-5 cyberdcli keys add <your_key_name> --recover
```

### Import an account by private key and store it in local keystore (private key could be your ETH private key)

```bash
docker exec -ti euler-5 cyberdcli keys add import_private <your_key_name>
```

### Create a new account

```bash
docker exec -ti euler-5 cyberdcli keys add <your_key_name>
```

### Show account information

Name, address and public key of current account

```bash
docker exec euler-5 cyberdcli keys show <your_key_name>
```

### Show account balance

Return account number, balance, public key in 16 and sequence.
>Don't work if from current account no outgoing transactions. [Issue in progress](https://github.com/cybercongress/cyberd/issues/238)

```bash
docker exec euler-5 cyberdcli query account <your_key_address>
```

### List existing keys

Return all keys in cyberdcli

```bash
docker exec euler-5 cyberdcli keys list
```

### Delete account from cybercli

```bash
docker exec -ti euler-5 cyberdcli keys delete <deleting_key_name>
```

### Update account password

```bash
docker exec -ti euler-5 cyberdcli keys update <your_key_name>
```

### Send tokens

```bash
docker exec -ti euler-5 cyberdcli tx send <from_key_or_address> <to_address> <amount_eul> \
  --chain-id=euler-5
```

### Linking content

>Just ipfs hashes available as a CID

```bash
docker exec -ti euler-5 cyberdcli link \
  --from=<your_key_name> \
  --cid-from=<key_phrase_to_link> \
  --cid-to=<content_that_you_want_to_link> \
  --chain-id=euler-5
```

## Validator commands

### Get all validators

```bash
docker exec euler-5 cyberdcli query staking validators \
    --trust-node
```

### The amount of commission

Available to withdraw validator commission.

```bash
docker exec euler-5 cyberdcli query distribution commission <operator_address>
```

### State of current validator

```bash
docker exec eiler-5 cyberdcli query staking validator <operator_address>
```

### Return all delegations to validator

```bash
docker exec euler-5 cyberdcli query staking delegations-to <operator_address>
```

### Edit commission in existing validator account

```bash
docker exec -ti euler-5 cyberdcli tx staking edit-validator \
  --from=<your_key_name> \
  --commission-rate=<new_comission_rate_percentage> \
  --chain-id=euler-5
```

### Withdraw commission for either a delegation

```bash
docker exec -ti euler-5 cyberdcli tx distribution withdraw-rewards <operator_address> \
  --from=<your_key_name> \
  --chain-id=euler-5 \
  --commission
```

### Edit site and description in existing validator account

>Will be available at description section

```bash
docker exec -ti euler-5 cyberdcli tx staking edit-validator \
  --from=<your_key_name> \
  --details="<description>" \
  --website=<your_website> \
  --chain-id=euler-5
```

### Unjail validator previously jailed for downtime

```bash
docker exec -ti euler-5 cyberdcli tx slashing unjail \
  --from=<your_key_name> \
  --chain-id=euler-5
```

### Get info about redelegation process from validator

```bash
docker exec -ti euler-5 cyberdcli query staking redelegations-from <operator_address>
```

## Delegator commands

### Return distribution delegator rewards according current validator

```bash
docker exec -ti euler-5 cyberdcli query distribution rewards <delegator_address> <operator_address>
```

### Return delegator shares with current validator

```bash
docker exec -ti euler-5 cyberdcli query staking delegation <delegator_address> <operator_address>
```

### Return all delegations made from one delegator

```bash
docker exec -ti euler-5 cyberdcli query staking delegations <delegator_address>
```

### Return all unbonding delegatations from a validator

```bash
docker exec -ti euler-5 cyberdcli query staking unbonding-delegations-from <operator_address>
```

### Withdraw rewards for either a delegation

```bash
docker exec -ti euler-5 cyberdcli tx distribution withdraw-rewards <operator_address> \
  --from=<your_key_name> \
  --chain-id=euler-5
```

### Withdraw all delegation rewards

```bash
docker exec -ti euler-5 cyberdcli tx distribution withdraw-all-rewards \
  --from=<your_key_name> \
  --chain-id=euler-5
```

### Change the default withdraw address for rewards associated with an address

```bash
docker exec -ti euler-5 cyberdcli tx distribution set-withdraw-addr <your_new_address> \
  --from=<your_key_name> \
  --chain-id=euler-5
```

### Delegate liquid tokens to a validator

```bash
docker exec -ti euler-5 cyberdcli tx staking delegate <operator_address> <amount_cyb> \
  --from=<your_key_name> \
  --chain-id=euler-5
```

### Redelegate illiquid tokens from one validator to another in absolute cyb value

>instant redelegation. Amount must be less than already delegated.

```bash
docker exec -ti euler-5 cyberdcli tx staking redelegate <old_operator_address> <new_operator_address> <amount_cyb> --from=<your_key_name> --chain-id=euler-5
```

### Redelegate illiquid tokens from one validator to another in percentages

```bash
docker exec -ti euler-5 cyberdcli tx staking redelegate <old_operator_address> <new_operator_address> <shares_percentage>
  --from=<your_key_name> \
  --chain-id=euler-5
```

### Unbond shares from a validator in absolute cyb value

> 5 days for unbonding.

```bash
docker exec -ti euler-5 cyberdcli tx staking unbond <operator_address> <amount_cyb>
  --from=<your_key_name> \
  --chain-id=euler-5
```

### Unbond shares from a validator in percentages

>5 days unbonding.

```bash
docker exec -ti euler-5 cyberdcli tx staking unbond <operator_address> <shares_percentage>
  --from=<your_key_name> \
  --chain-id=euler-5
```

### Get info about unbonding delegation process to current validator

```bash
docker exec -ti euler-5 cyberdcli query staking unbonding-delegation <delegator_address> <operator_address>
```

### Get info about unbonding delegation process to all unbonded validators

```bash
docker exec -ti euler-5 cyberdcli query staking unbonding-delegation <delegator_address>
```

### Get info about redelegation process from to current validator

```bash
docker exec -ti euler-5 cyberdcli query staking redelegation <delegator_address> <old_operator_address> <new_operator_address>
```

### Get info about all redelegation processes by one delegator

```bash
docker exec -ti euler-5 cyberdcli query staking redelegations <delegator_address>
```

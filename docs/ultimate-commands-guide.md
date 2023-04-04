# Ultimate cyber CLI guide. Chain: bostrom

## Install cyber client

It is possible to interact with cyber even if you don't have your own node. All you need to do is install `cyber` client on your machine using the script below, paste it in the console (currently Linux supported):

```bash
bash < <(curl -s https://raw.githubusercontent.com/cybercongress/go-cyber/main/scripts/install_cyber.sh)
```

After installation, you will be able to use `cyber` to [import accounts](#account-management), create links or swap tokens.

In case you have your own node, which is already running inside Docker container, please add `docker exec -ti container-name` before every cyber command:

```bash
docker exec -ti bostrom cyber --help
```

First of all, I would like to encourage you to use the  `--help` feature if you want to have a better experience using cyber. This is a really easy way to find all the necessary commands with the appropriate options and flags.

For example, you can enter:

```bash
cyber --help
```

You should see this message:

```bash
Usage:
  cyber [command]

Available Commands:
  add-genesis-account Add a genesis account to genesis.json
  collect-gentxs      Collect genesis txs and output a genesis.json file
  config              Create or query an application CLI configuration file
  debug               Tool for helping with debugging your application
  export              Export state to JSON
  gentx               Generate a genesis tx carrying a self delegation
  help                Help about any command
  init                Initialize private validator, p2p, genesis, and application configuration files
  keys                Manage your applications keys
  migrate             Migrate genesis to a specified target version
  query               Querying subcommands
  start               Run the full node
  status              Query remote node for status
  tendermint          Tendermint subcommands
  testnet             Initialize files for a simapp testnet
  tx                  Transactions subcommands
  unsafe-reset-all    Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state
  validate-genesis    validates the genesis file at the default location or at the location passed as an arg
  version             Print the application binary version information
```

The help feature works like a pyramid, you can use it with any command to find available options, subcommands and flags. For example, lets explore the `query` subcommands:

```bash
cyber query --help
```

You can see the structure of the subcommand:

```bash
Usage:
  cyber query [flags]
  cyber query [command]
```

And the available subcommands and flags:

```bash
Querying subcommands

Aliases:
  query, q

Available Commands:
  account                  Query for account by address
  auth                     Querying commands for the auth module
  authz                    Querying commands for the authz module
  bandwidth                Querying commands for the bandwidth module
  bank                     Querying commands for the bank module
  block                    Get verified data for a the block at given height
  distribution             Querying commands for the distribution module
  dmn                      Querying commands for the dmn module
  evidence                 Query for evidence by hash or for all (paginated) submitted evidence
  feegrant                 Querying commands for the feegrant module
  gov                      Querying commands for the governance module
  graph                    Querying commands for the graph module
  grid                     Querying commands for the grid module
  ibc                      Querying commands for the IBC module
  ibc-transfer             IBC fungible token transfer query subcommands
  liquidity                Querying commands for the liquidity module
  mint                     Querying commands for the minting module
  params                   Querying commands for the params module
  rank                     Querying commands for the rank module
  resources                Querying commands for the resources module
  slashing                 Querying commands for the slashing module
  staking                  Querying commands for the staking module
  tendermint-validator-set Get the full tendermint validator set at given height
  tx                       Query for a transaction by hash, "<addr>/<seq>" combination or comma-separated signatures in a committed block
  txs                      Query for paginated transactions that match a set of events
  upgrade                  Querying commands for the upgrade module
  wasm                     Querying commands for the wasm module

Flags:
      --chain-id string   The network chain ID
  -h, --help              help for query

Global Flags:
      --home string         directory for config and data (default "/root/.cyber")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors

Use "cyber query [command] --help" for more information about a command.

```

Let's explore the `bank` subcommand:

```bash
cyber query bank --help
```

We can see all of the options available for these subcommands, names and account address:

```bash
Usage:
  cyber query bank [flags]
  cyber query bank [command]

Available Commands:
  balances       Query for account balances by address
  denom-metadata Query the client metadata for coin denominations
  total          Query the total supply of coins of the chain
```

In most cases you will need just two extra flags:

```bash
--from=<your_key_name> \
--chain-id=bostrom \
--node=<rpc_node_path>
```

That's it. This is a very useful tool for using cyber and troubleshooting.

## Glossary

**Bonded tokens** - Tokens that are nominated to a validator, non transferable.

**Commission** -  The tokens that you've earned via validating from delegators.

**Dyson Sphere** - Construct in cyber responsible for energy transformation and routing.

**Investminting** - Process of convertation Hydrogen to Volts or Amperes, locking a certain amount of H for a certain amount of time produces some V or A.

**Hero** - A validator.

**Hydrogen** - Token issued after boot delegation, 1:1 H to boot. Used to generate enery through the Dyson sphere (investminting process).

**Unbonding** - The process of taking back your share (delegated tokens + any rewards). 4 days for `bostrom` chain.


**link** - A reference between a CID key and a CID value. Link message cost is `100*n`, where `n` is the number of links in a message. Link finalization time is 1 block. New rank for CIDs of links will be recalculated at a period of 5 blocks.

**liquid tokens** - Transferable tokens within the cyber.

**local keystore** - A store with keys on your local machine.

**rewards** - Tokens that a hero earned via delegation. To reduce network load all the rewards are stored in a pool. You can take your part of the bounty at any time with commands from the **delegator** section.

**<comission_rate_percentage>** - The commission that a validator gets for their work. Must be a fraction >0 and <=1

**<delegator_address>** - Delegator address. Starts with `bostrom` most often coinciding with **<key_address>**

**<key_address>** - An account address. Starts with `bostrom`

**<key_name>** - The name of the account in cyber

**<operator_address>** - Validator address. Starts with `cybervaloper`

**<shares_percentage>** - The part of illiquid tokens that you want to unbond or redelegate. Must be a fraction >0 and <=1

**<chain_id>** - The current version of the chain  (bostrom).

## General commands

### Show all validators

Return the set of all active and jailed validators:

```bash
cyber query staking validators 
```

### Show chain status

Return general chain information:

```bash
 cyber status
```

### Distribution params

```bash
 cyber query distribution params 
```

### Staking params

Chain staking info:

```bash
 cyber query staking params 
```

### Staking pool

```bash
 cyber query staking pool 
```

## Account management

### Import an account with a seed phrase and store it in the local keystore

```bash
 cyber keys add <your_key_name> --recover
```

### Create a new account

```bash
 cyber keys add <your_key_name>
```

### Show account information

Name, address and the public key of the current account

```bash
 cyber keys show <your_key_name>
```

### Show account balance

Return account number and amount of tokens.

```bash
 cyber query bank balances <your_key_address>
```

### List existing keys

Return all the existing keys in cyber:

```bash
 cyber keys list
```

### Delete account from cyber

```bash
 cyber keys delete <deleting_key_name>
```

### Keyring manipulation settings

**Important note**: Starting with v.38, Cosmos-SDK uses os-native keyring to store all of your keys. We've noticed that in certain cases it does not work well by default (for example if you don't have any GUI installed on your machine). If during the execution `cyber keys add` command, you are getting this type of error:

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

You will have to use another keyring backend to keep your keys. 

In that case you'll have to use file based keyring by adding the `--keyring-backend file` option to every key manipulation command:

```bash
cyber keys add key_name --keyring-backend file
```

That means that you've set your keyring-backend to a local file. *Note*, in this case, all the keys in your keyring will be encrypted using the same password. If you would like to set up a unique password for each key, you should set a unique `--home` folder for each key. To do that, just use `--home=/<unique_path_to_key_folder>/` with setup keyring backend and at all interactions with keys when using cyber:

```bash
cyber keys add <your_second_key_name> --keyring-backend file --home=/<unique_path_to_key_folder>/
cyber keys list --home=/<unique_path_to_key_folder>/
```

### Send tokens

```bash
cyber tx bank send [from_key_or_address] [to_address] [amount] \
--from=<your_key_name> \
--chain-id=bostrom
```

### Linking content

Only IPFS hashes are available to use as CIDs

```bash
cyber tx graph cyberlink [cid-from] [cid-to] [flags] \
--from=<your_key_name> \
--chain-id=bostrom
```

Real command example:

```bash 
cyber tx graph cyberlink QmWZYRj344JSLShtBnrMS4vw5DQ2zsGqrytYKMqcQgEneB QmfZwbahFLTcB3MTMT8TA8si5khhRmzm7zbHToo4WVK3zn --from fuckgoogle --chain-id bostrom --yes
```

## Validator commands

### Get all validators

```bash
 cyber query staking validators 
```

### State of a current validator

```bash
cyber query staking validator <operator_address>
```

### Return all delegations to a validator

```bash
 cyber query staking delegations-to <operator_address>
```

### Edit the commission in an existing validator account

```bash
 cyber tx staking edit-validator \
  --from=<your_key_name> \
  --commission-rate=<new_comission_rate_percentage> \
  --chain-id=bostrom
```

### Withdraw the commission for any delegation

```bash
 cyber tx distribution withdraw-rewards <operator_address> --commission \
  --from=<your_key_name> \
  --chain-id=bostrom
```

### Edit the site and description for an existing validator account

```bash
 cyber tx staking edit-validator \
  --from=<your_key_name> \
  --details="<description>" \
  --website=<your_website> \
  --chain-id=bostrom
```

### Unjail a validator previously jailed for downtime

```bash
 cyber tx slashing unjail --from=<your_key_name> --chain-id=bostrom
```

### Get info about a redelegation process from a validator

```bash
 cyber query staking redelegations-from <operator_address>
```

## Delegator commands

### Return distribution delegator rewards for a specified validator

```bash
 cyber query distribution rewards <delegator_address> <operator_address>
```

### Return delegator shares for the specified validator

```bash
 cyber query staking delegation <delegator_address> <operator_address>
```

### Return all delegations made from a delegator

```bash
 cyber query staking delegations <delegator_address>
```

### Return all unbonding delegations from a validator

```bash
 cyber query staking unbonding-delegations-from <operator_address>
```

### Withdraw rewards for any delegation

```bash
 cyber tx distribution withdraw-rewards <operator_address> \
  --from=<your_key_name> \
  --chain-id=bostrom
```

### Withdraw all delegation rewards

```bash
 cyber tx distribution withdraw-all-rewards \
  --from=<your_key_name> \
  --chain-id=bostrom
```

### Change the default withdrawal address for rewards associated with an address

```bash
 cyber tx distribution set-withdraw-addr <your_new_address> \
  --from=<your_key_name> \
  --chain-id=bostrom
```

### Delegate liquid tokens to a validator

```bash
 cyber tx staking delegate <operator_address> <amount_boot> \
  --from=<your_key_name> \
  --chain-id=bostrom
```

### Redelegate illiquid tokens from one validator to another in absolute BOOT value

> There is a 4-day unbonding period

```bash
 cyber tx staking redelegate <old_operator_address> <new_operator_address> <amount_boot> \
 --from=<your_key_name> \
 --chain-id=bostrom
```

### Redelegate illiquid tokens from one validator to another in percentages

```bash
 cyber tx staking redelegate <old_operator_address> <new_operator_address> <shares_percentage>
  --from=<your_key_name> \
  --chain-id=bostrom
```

### Unbond shares from a validator in absolute BOOT value

> 8 days for unbonding

```bash
 cyber tx staking unbond <operator_address> <amount_boot>
  --from=<your_key_name> \
  --chain-id=bostrom
```

### Unbond shares from a validator in percentages

> 8 days for unbonding

```bash
 cyber tx staking unbond <operator_address> <shares_percentage>
  --from=<your_key_name> \
  --chain-id=bostrom
```

### Get info about the unbonding delegation process to any validator

```bash
 cyber query staking unbonding-delegation <delegator_address> <operator_address>
```

### Get info about the unbonding delegation process to all unbonded validators

```bash
 cyber query staking unbonding-delegation <delegator_address>
```

### Get info about redelegation process from to current validator

```bash
 cyber query staking redelegation <delegator_address> <old_operator_address> <new_operator_address>
```

### Get the info about all the redelegation processes by a delegator

```bash
 cyber query staking redelegations <delegator_address>
```

## Governance and voting

### Query specific proposal

```bash
cyber q gov proposal <proposal_id> 
```

### Query all proposals

```bash
cyber q gov proposals 
```

### Query votes on proposal

```bash
cyber q gov votes 
```

### Query parameters from governance module

```bash
cyber q gov params
```

### Vote for specific proposal

```bash
cyber tx gov vote <proposal_id> <vote_option:_yes_no_abstain> \
--from=<your_key_name> \
--chain-id=bostrom
```

### Submit text proposal

```bash
cyber tx gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10boot" \
--from=<your_key_name> \
--chain-id=bostrom
```

### Submit community spend proposal

```bash
cyber tx gov submit-proposal community-pool-spend <path/to/proposal.json> \
--from=<key_or_address> \
--chain-id=bostrom
```

Where `proposal.json` is a file, structured in following way:

```json
{
  "title": "Community Pool Spend",
  "description": "Pay me some boots!",
  "recipient": "bostrom1s5afhd6gxevu37mkqcvvsj8qeylhn0rz46zdlq",
  "amount": [
    {
      "denom": "boot",
      "amount": "10000"
    }
  ],
  "deposit": [
    {
      "denom": "boot",
      "amount": "10000"
    }
  ]
}
```

### Submit chain parameters change proposal

```bash
cyber tx gov submit-proposal param-change <path/to/proposal.json> \
--from=<key_or_address> \
--chain-id=bostrom
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

Few examples of real param-change proposals:

- Change max block bandwidth:

  ```json
  {
    "title": "Insrease max block bandwidth to 500000",
    "description": "Increase max block bandwidth to 500000.\n",
    "changes": [
      {
        "subspace": "bandwidth",
        "key": "MaxBlockBandwidth",
        "value": "500000"
      }
    ],
    "deposit": "10000boot"
  }
  ```
- Increase max block gas:

  ```json
  {
    "title": "Increase max gasprice",
    "description": "Increase  maximum block gas.\n",
    "changes": [
      {
        "subspace": "baseapp",
        "key": "BlockParams",
        "value": '"{\n                \"max_bytes\": \"22020096\",  \n                \"max_gas\": \"200000000\",\n                 \"time_iota_ms\": \"1000\"\n            }"'
      }
    ],
    "deposit": "10000boot"
  }
  ```

- Change rank calculation window:

  ```json
  {
    "title": "Change rank calculation window to 25",
    "description": "Increase rank calculation window from every 5 to every 25 blocks.\n",
    "changes": [
      {
        "subspace": "rank",
        "key": "CalculationPeriod",
        "value": "25"
      }
    ],
    "deposit": "10000boot"
  }
  ```

## Liquidity and pools 

Cyber has Gravity-DEX module implemented, so it is possible to create pools, swap tokens and add\remove liquidity to exisitng pools: 

```bash
Liquidity transaction subcommands

Usage:
  cyber tx liquidity [flags]
  cyber tx liquidity [command]

Available Commands:
  create-pool Create liquidity pool and deposit coins
  deposit     Deposit coins to a liquidity pool
  swap        Swap offer coin with demand coin from the liquidity pool with the given order price
  withdraw    Withdraw pool coin from the specified liquidity pool
```

### Create liquidity pool


This example creates a liquidity pool of pool-type 1 (two coins) and deposits 2000000milliampere and 200000000000boot.
New liquidity pools can be created only for coin combinations that do not already exist in the network.

[pool-type]: The id of the liquidity pool-type. The only supported pool type is 1
[deposit-coins]: The amount of coins to deposit to the liquidity pool. The number of deposit coins must be 2 in pool type 1.

Example:

```bash
cyber tx liquidity create-pool 1 2000000milliampere,200000000000boot \
--from=mykey \
--chain-id=bostrom  \
--yes
```

### Deposit tokens to liquidity pool

Deposit coins a liquidity pool.

This deposit request is not processed immediately since it is accumulated in the liquidity pool batch.
All requests in a batch are treated equally and executed at the same swap price.

Example:

```bash
cyber tx liquidity deposit 1 100000000milliampere,50000000000boot \
--from=mykey \
--chain-id=bostrom
```

This example request deposits 100000000milliampere and 50000000000boot to pool-id 1.
Deposits must be the same coin denoms as the reserve coins.

[pool-id]: The pool id of the liquidity pool
[deposit-coins]: The amount of coins to deposit to the liquidity pool

### Swap coins

Swap offer coin with demand coin from the liquidity pool with the given order price.

This swap request is not processed immediately since it is accumulated in the liquidity pool batch.
All requests in a batch are treated equally and executed at the same swap price.
The order of swap requests is ignored since the universal swap price is calculated in every batch to prevent front running.

The requested swap is executed with a swap price that is calculated from the given swap price function of the pool, the other swap requests, and the liquidity pool coin reserve status.
Swap orders are executed only when the execution swap price is equal to or greater than the submitted order price of the swap order.

Example:

```bash
cyber tx liquidity swap 1 1 50000000boot hydrogen 0.019 0.003 \
--from=mykey \
--chain-id=bostrom
```

For this example, imagine that an existing liquidity pool has with 1000000000hydrogen and 50000000000boot.
This example request swaps 50000000boot for at least 950000hydrogen with the order price of 0.019 and swap fee rate of 0.003.
A sufficient balance of half of the swap-fee-rate of the offer coin is required to reserve the offer coin fee.

The order price is the exchange ratio of X/Y, where X is the amount of the first coin and Y is the amount of the second coin when their denoms are sorted alphabetically.
Increasing order price reduces the possibility for your request to be processed and results in buying hydrogen at a lower price than the pool price.

For explicit calculations, The swap fee rate must be the value that is set as liquidity parameter in the current network.
The only supported swap-type is 1. For the detailed swap algorithm, see https://github.com/gravity-devs/liquidity

[pool-id]: The pool id of the liquidity pool
[swap-type]: The swap type of the swap message. The only supported swap type is 1 (instant swap).
[offer-coin]: The amount of offer coin to swap
[demand-coin-denom]: The denomination of the coin to exchange with offer coin
[order-price]: The limit order price for the swap order. The price is the exchange ratio of X/Y where X is the amount of the first coin and Y is the amount of the second coin when their denoms are sorted alphabetically
[swap-fee-rate]: The swap fee rate to pay for swap that is proportional to swap amount. The swap fee rate must be the value that is set as liquidity parameter in the current network.

Usage:

```bash
  cyber tx liquidity swap [pool-id] [swap-type] [offer-coin] [demand-coin-denom] [order-price] [swap-fee-rate] [flags]
```

### Withdraw tokens from liquidity pool

Withdraw pool coin from the specified liquidity pool.

This swap request is not processed immediately since it is accumulated in the liquidity pool batch.
All requests in a batch are treated equally and executed at the same swap price.

Example:

```bash
 cyber tx liquidity withdraw 1 10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295 \
 --from=mykey \
 --chain-id=bostrom
```

This example request withdraws 10000 pool coin from the specified liquidity pool.
The appropriate pool coin must be requested from the specified pool.

[pool-id]: The pool id of the liquidity pool
[pool-coin]: The amount of pool coin to withdraw from the liquidity pool

Usage:

```bash
  cyber tx liquidity withdraw [pool-id] [pool-coin] [flags]\
  --from=<key_or_address> \
  --chain-id=bostrom
```

### Query existing pools

Query details of a liquidity pool

```bash
cyber query liquidity pool 1
```

Example (with pool coin denom):

```bash
cyber query liquidity pool --pool-coin-denom=[denom]
```

Query details about all liquidity pools on a network.
Example:

```bash
cyber query liquidity pools
```

## 

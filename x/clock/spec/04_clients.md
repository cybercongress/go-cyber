<!--
order: 4
-->

# Clients

## Command Line Interface (CLI)

The CLI has been updated with new queries and transactions for the `x/clock` module. View the entire list below.

### Queries

| Command             | Subcommand  | Arguments          | Description             |
|:--------------------| :---------- | :----------------- | :---------------------- |
| `cyber query clock` | `params`    |                    | Get Clock params        |
| `cyber query clock` | `contract`  | [contract_address] | Get a Clock contract    |
| `cyber query clock` | `contracts` |                    | Get all Clock contracts |

### Transactions

| Command          | Subcommand   | Arguments          | Description                 |
|:-----------------| :----------- | :----------------- | :-------------------------- |
| `cyber tx clock` | `register`   | [contract_address] | Register a Clock contract   |
| `cyber tx clock` | `unjail`     | [contract_address] | Unjail a Clock contract     |
| `cyber tx clock` | `unregister` | [contract_address] | Unregister a Clock contract |
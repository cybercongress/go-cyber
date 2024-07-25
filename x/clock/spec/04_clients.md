<!--
order: 4
-->

# Clients

## Command Line Interface (CLI)

The CLI has been updated with new queries and transactions for the `x/clock` module. View the entire list below.

### Queries

| Command             | Subcommand  | Arguments          | Description             |
| :------------------ | :---------- | :----------------- | :---------------------- |
| `junod query clock` | `params`    |                    | Get Clock params        |
| `junod query clock` | `contract`  | [contract_address] | Get a Clock contract    |
| `junod query clock` | `contracts` |                    | Get all Clock contracts |

### Transactions

| Command          | Subcommand   | Arguments          | Description                 |
| :--------------- | :----------- | :----------------- | :-------------------------- |
| `junod tx clock` | `register`   | [contract_address] | Register a Clock contract   |
| `junod tx clock` | `unjail`     | [contract_address] | Unjail a Clock contract     |
| `junod tx clock` | `unregister` | [contract_address] | Unregister a Clock contract |
<!--
order: 1
-->

# Concepts

## Clock

The Clock module allows registered contracts to be executed at the end of every block. This allows the smart contract to perform regular and routine actions without the need for external bots. Developers can setup their contract with x/Clock by registering their contract with the module. Once registered, the contract will be executed at the end of every block. If the contract throws an error during execution or exceeds the gas limit defined in the module's parameters, the contract will be jailed and no longer executed. The contract can be unjailed by the contract admin.

## Registering a Contract

Register a contract with x/Clock by executing the following transaction:

```bash
junod tx clock register [contract_address]
```

> Note: the sender of this transaction must be the contract admin, if exists, or else the contract creator.

The `contract_address` is the bech32 address of the contract to be executed at the end of every block. Once registered, the contract will be executed at the end of every block. Please ensure that your contract follows the guidelines outlined in [Integration](03_integration.md). 

## Unjailing a Contract

A contract can be unjailed by executing the following transaction:

```bash
junod tx clock unjail [contract_address]
```

> Note: the sender of this transaction must be the contract admin, if exists, or else the contract creator.

The `contract_address` is the bech32 address of the contract to be unjailed. Unjailing a contract will allow it to be executed at the end of every block. If your contract becomes jailed, please see [Integration](03_integration.md) to ensure the contract is setup with a Sudo message. 

## Unregistering a Contract

A contract can be unregistered by executing the following transaction:

```bash
junod tx clock unregister [contract_address]
```

> Note: the sender of this transaction must be the contract admin, if exists, or else the contract creator.

The `contract_address` is the bech32 address of the contract to be unregistered. Unregistering a contract will remove it from the Clock module. This means that the contract will no longer be executed at the end of every block.